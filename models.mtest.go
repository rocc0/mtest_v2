package main

import (
	"encoding/json"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

type governmentRegion struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type newExecutors struct {
	Email   string `json:"email"`
	Mid     string `json:"mid"`
	Checked bool   `json:"checked"`
}

type Mtest struct {
	Id           int       `json:"id"`
	Mid          uuid.UUID `json:"mid"`
	Name         string    `json:"name"`
	Region       string    `json:"region"`
	Govern       string    `json:"govern"`
	Calculations string    `json:"calculations"`
	CalcType     int       `json:"calc_type"`
	CalcData     string    `json:"calc_data"`
	Executors    string    `json:"executors"`
	PubDate      string    `json:"pub_date"`
	Author       string    `json:"author"`
}

type (
	UserMtest struct {
		Id         string                  `json:"id"`
		Name       string                  `json:"name"`
		Region     int                     `json:"region"`
		Government int                     `json:"govern"`
		CalcType   int                     `json:"calc_type"`
		Developer  string                  `json:"developer"`
		DevMid     string                  `json:"dev_mid"`
		Executors  map[string]newExecutors `json:"executors"`
	}
	AdmAction struct {
		ActId   int    `json:"act_id"`
		ActName string `json:"act_name"`
	}
)

const calculations = `{"1":[{"type":"container","id":3,"columns":[[{"type":"itemplus","id":3,
                    "columns":[[{"type":"item","id":3,"name":"Додати дію","subsum":0},{"type":"item","id":6,"name":"Додати дію","subsum":0}]],
                    "name":"Додати складову інф. вимоги"}]],"name":"Додати інф. вимогу","contsub":0},
                {"type":"container","id":null,"columns":[[{"type":"itemplus","id":4,"columns":[[{"type":"item","id":3,"name":"Додати дію","subsum":0},
                            {"type":"item","id":4,"name":"Додати дію","subsum":0}]],"name":"Додати складову інф. вимоги"}]],"name":"Додати інф. вимогу","contsub":0}]}`

func createNewMTest(m newMtest, email string) (*map[string]interface{}, error) {
	var (
		id        int
		dbRecords string
		records   map[string]interface{}
	)

	stmt, err := db.Prepare("INSERT INTO mtests (mid, name, region, govern," +
		" calculations, calc_type, pub_date, author) VALUES (?,?,?,?,?,?,?,?)")
	if err != nil {
		return nil, err
	}

	mId := uuid.New()
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()

	result, err := stmt.Exec(mId, m.Name, m.Region, m.Government, calculations, m.CalcType, time.Now(), email)
	if err != nil {
		return nil, err
	}

	idRes, _ := result.LastInsertId()
	idStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
	if err := idStmt.Scan(&id, &dbRecords); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(dbRecords), &records); err != nil {
		return nil, err
	}

	records[mId.String()] = UserMtest{Id: mId.String(), Name: m.Name, Region: m.Region,
		Government: m.Government, CalcType: m.CalcType}

	out, err := json.Marshal(records)
	if err != nil {
		return nil, err
	}

	if err = updateUser("records", string(out), id); err != nil {
		return nil, err
	}

	//need check
	if err = updateIndex(idRes); err != nil {
		return nil, err
	}

	return &records, nil
}

func readMtest(id string) (*Mtest, error) {
	var (
		mid                  uuid.UUID
		rowId, calcType      int
		name, govern         string
		region, calculations string
		calcData, executors  string
		pubDate, author      string
	)
	res := db.QueryRow("SELECT m.id, m.mid, m.name, r.reg_name, g.gov_name, m.calculations,"+
		" m.calc_type, m.calc_data, m.executors, m.pub_date, m.author FROM mtests m JOIN govs g ON m.govern = g.id JOIN "+
		"regions r ON m.region = r.reg_id WHERE mid=?", id)

	if err := res.Scan(&rowId, &mid, &name, &region, &govern, &calculations,
		&calcType, &calcData, &executors, &pubDate, &author); err != nil {
		return nil, err
	}
	return &Mtest{Id: rowId, Mid: mid, Name: name, Region: region,
		Govern: govern, Calculations: calculations, CalcType: calcType,
		CalcData: calcData, Executors: executors, PubDate: pubDate, Author: author}, nil
}

func updateMtest(m map[string]interface{}, email string) error {
	var (
		id        int
		dbRecords string
		records   map[string]UserMtest
	)

	if m["calculations"] == nil && m["name"] != nil {
		stmt, err := db.Prepare("UPDATE mtests SET name=?, region=?, govern=? WHERE mid=?;")
		if err != nil {
			return err
		}
		defer func() {
			if err := stmt.Close(); err != nil {
				log.Error(err)
			}
		}()
		if m["calculations"] != nil && m["name"] == nil {
			stmt, err := db.Prepare("UPDATE mtests SET calculations= ? WHERE mid=?;")
			if err != nil {
				return err
			}
			if _, err := stmt.Exec(m["calculations"], m["id"]); err != nil {
				return err
			}
			return nil
		} else if m["executors"] != nil && m["name"] == nil {
			stmt, err := db.Prepare("UPDATE mtests SET executors= ? WHERE mid=?;")
			if err != nil {
				return err
			}
			if _, err := stmt.Exec(m["executors"], m["id"]); err != nil {
				return err
			}
			return nil
		} else {
			govern := int(m["govern"].(float64))
			region := int(m["region"].(float64))
			calcType := int(m["calc_type"].(float64))
			if _, err := stmt.Exec(m["name"], region, govern, m["mid"]); err != nil {
				return err
			}

			idStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
			if err := idStmt.Scan(&id, &dbRecords); err != nil {
				return err
			}

			if err := json.Unmarshal([]byte(dbRecords), &records); err != nil {
				return err
			}
			record := records[m["mid"].(string)]
			records[m["mid"].(string)] = UserMtest{Id: m["mid"].(string), Name: m["name"].(string),
				Region: region, Government: govern, CalcType: calcType, Executors: record.Executors}

			out, err := json.Marshal(records)
			if err != nil {
				return err
			}
			if err := updateUser("records", string(out), id); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func deleteMtest(mid, email string) error {
	var (
		id        int
		dbRecords string
		records   map[string]interface{}
	)
	stmt, err := db.Prepare("DELETE FROM mtests WHERE mid=?")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()
	if _, err := stmt.Exec(mid); err != nil {
		return err
	}

	idStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
	if err := idStmt.Scan(&id, &dbRecords); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(dbRecords), &records); err != nil {
		return err
	}

	delete(records, mid)
	out, err := json.Marshal(records)
	if err != nil {
		return err
	}

	if err := updateUser("records", string(out), id); err == nil {
		return err
	}
	return nil
}

func getGovs() (*[]governmentRegion, error) {
	var (
		govs    []governmentRegion
		govId   int
		govName string
	)
	res, err := db.Query("SELECT gov_id, gov_name FROM govs")
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Close(); err != nil {
			log.Error(err)
		}
	}()
	for res.Next() {
		if err = res.Scan(&govId, &govName); err != nil {
			return nil, err
		}
		govs = append(govs, governmentRegion{Id: govId, Name: govName})
	}
	return &govs, nil
}

func getRegs() (*[]governmentRegion, error) {
	var (
		regions []governmentRegion
		regId   int
		regName string
	)

	res, err := db.Query("SELECT reg_id, reg_name FROM regions")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Close(); err != nil {
			log.Error(err)
		}
	}()

	for res.Next() {
		if err = res.Scan(&regId, &regName); err != nil {
			return nil, err
		}
		regions = append(regions, governmentRegion{Id: regId, Name: regName})
	}
	return &regions, nil
}

func getAdmactions() (*[]AdmAction, error) {
	var (
		actId   int
		actName string
		actions []AdmAction
	)

	res, err := db.Query("SELECT act_id, act_name FROM adm_actions")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Close(); err != nil {
			log.Error(err)
		}
	}()
	for res.Next() {
		if err := res.Scan(&actId, &actName); err != nil {
			return nil, err
		}
		actions = append(actions, AdmAction{ActId: actId, ActName: actName})
	}

	return &actions, nil

}

func createMtestExecutor(email string, ex newExecutor) (*uuid.UUID, error) {
	var (
		id, devId                             int
		dbRecords, devDbRecords, getExecutors string
		records                               map[string]interface{}
		devRecords                            map[string]UserMtest
	)

	if isUsernameAvailable(ex.Email) == true {
		return nil, errors.New("користувач не зареєстрований")
	}

	//add mtest type 3
	stmt, _ := db.Prepare("INSERT INTO mtests (mid, name, region, govern," +
		" calculations, calc_type, developer, dev_mid, pub_date, author) VALUES (?,?,?,?,?,?,?,?,?,?)")

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()
	mId := uuid.New()
	if _, err := stmt.Exec(mId, ex.Title, ex.Region, ex.Government,
		calculations, 3, email, ex.DevMid, time.Now(), ex.Email); err != nil {
		return nil, err
	}

	//UPDATE MAIN MTEST executors!!!!!!!!!!!
	saveExecutors := map[string]newExecutors{}
	exStmt := db.QueryRow("SELECT executors FROM mtests WHERE mid=?", ex.DevMid)
	if err := exStmt.Scan(&getExecutors); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(getExecutors), &saveExecutors); err != nil {
		return nil, err
	}

	saveExecutors[mId.String()] = newExecutors{ex.Email, mId.String(), true}
	updOut, updOutErr := json.Marshal(saveExecutors)
	if updOutErr != nil {
		return nil, updOutErr
	}

	updMtest, updMtestErr := db.Prepare("UPDATE mtests SET executors=? WHERE mid=?")
	if updMtestErr != nil {
		return nil, updMtestErr
	}
	_, updError := updMtest.Exec(updOut, ex.DevMid)
	if updError != nil {
		return nil, updError
	}

	//add mtest to executor mtests
	idStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", ex.Email)

	if err := idStmt.Scan(&id, &dbRecords); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(dbRecords), &records); err != nil {
		return nil, err
	}

	records[mId.String()] = UserMtest{Id: mId.String(), Name: ex.Title, Region: ex.Region,
		Government: ex.Government, CalcType: 3, Developer: email, DevMid: ex.DevMid}

	out, marshErr := json.Marshal(records)
	if marshErr != nil {
		return nil, marshErr
	}

	if idErr := updateUser("records", string(out), id); idErr != nil {
		return nil, idErr
	}

	// add mtest to developer mtest to executors section
	devStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
	if err := devStmt.Scan(&devId, &devDbRecords); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(devDbRecords), &devRecords); err != nil {
		return nil, err
	}
	record := devRecords[ex.DevMid]
	devRecords[ex.DevMid] = UserMtest{Id: ex.DevMid, Name: record.Name, Region: record.Region,
		Government: record.Government, CalcType: record.CalcType, Executors: saveExecutors}

	devOut, err := json.Marshal(devRecords)
	if err != nil {
		return nil, err
	}

	if devErr := updateUser("records", string(devOut), devId); devErr != nil {
		return nil, devErr
	}

	return &mId, nil
}

func deleteMtestExecutor(devEmail string, del delExecutorReq) error {
	var (
		id         int
		dbRecords  string
		records    map[string]interface{}
		devRecords map[string]UserMtest
	)
	//delete mtest

	stmt, err := db.Prepare("DELETE FROM mtests WHERE mid=?")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()
	if _, err := stmt.Exec(del.ExMtestId); err != nil {
		return err
	}

	//delete from executors mtests
	idStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", del.ExEmail)
	if err := idStmt.Scan(&id, &dbRecords); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(dbRecords), &records); err != nil {
		return err
	}

	delete(records, del.ExMtestId)
	out, err := json.Marshal(records)
	if err != nil {
		return err
	}

	if err := updateUser("records", string(out), id); err != nil {
		return err
	}

	//delete from developers mtests
	devStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", devEmail)
	if err := devStmt.Scan(&id, &dbRecords); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(dbRecords), &devRecords); err != nil {
		return err
	}
	delete(devRecords[del.DevMtestId].Executors, del.ExMtestId)
	devOut, err := json.Marshal(devRecords)
	if err != nil {
		return err
	}

	if err := updateUser("records", string(devOut), id); err != nil {
		return err
	}

	//delete from developer mtest
	mtStmt := db.QueryRow("SELECT executors FROM mtests WHERE mid=?", del.DevMtestId)
	if err := mtStmt.Scan(&dbRecords); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(dbRecords), &records); err != nil {
		return err
	}
	delete(records, del.ExMtestId)
	mtOut, err := json.Marshal(records)
	if err != nil {
		return err
	}

	mtSaveStmt, _ := db.Prepare("UPDATE mtests SET executors=? WHERE mid=?")

	if _, err := mtSaveStmt.Exec(mtOut, del.DevMtestId); err != nil {
		return err
	}
	return nil

}
