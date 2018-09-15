package main

import (
	"time"
	"errors"
	"github.com/google/uuid"
	"encoding/json"
)

type governmentRegion struct {
	Id int `json:"id"`
	Name string	`json:"name"`
}

type newExecutors struct {
	Email string `json:"email"`
	Mid string `json:"mid"`
	Checked bool `json:"checked"`
}

type Mtest struct {
	Id int `json:"id"`
	Mid uuid.UUID `json:"mid"`
	Name string `json:"name"`
	Region string `json:"region"`
	Govern string `json:"govern"`
	Calculations string `json:"calculations"`
	CalcType int `json:"calc_type"`
	CalcData string `json:"calc_data"`
	Executors string `json:"executors"`
	PubDate string `json:"pub_date"`
	Author string `json:"author"`
}

type UserMtest struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Region int `json:"region"`
	Government int `json:"govern"`
	CalcType int `json:"calc_type"`
	Developer string `json:"developer"`
	DevMid string `json:"dev_mid"`
	Executors map[string]newExecutors `json:"executors"`
}

type AdmAction struct {
	ActId int `json:"act_id"`
	ActName string `json:"act_name"`
}

const calculations = `{"1":[{"type":"container","id":3,"columns":[[{"type":"itemplus","id":3,
                    "columns":[[{"type":"item","id":3,"name":"Додати дію","subsum":0},{"type":"item","id":6,"name":"Додати дію","subsum":0}]],
                    "name":"Додати складову інф. вимоги"}]],"name":"Додати інф. вимогу","contsub":0},
                {"type":"container","id":null,"columns":[[{"type":"itemplus","id":4,"columns":[[{"type":"item","id":3,"name":"Додати дію","subsum":0},
                            {"type":"item","id":4,"name":"Додати дію","subsum":0}]],"name":"Додати складову інф. вимоги"}]],"name":"Додати інф. вимогу","contsub":0}]}`

func createNewMtest(m newMtest, email string) (*map[string]interface{}, error) {
	var (
		id int
		dbRecords string
		records map[string]interface{}
		)

	stmt, err := db.Prepare("INSERT INTO mtests (mid, name, region, govern," +
		" calculations, calc_type, pub_date, author) VALUES (?,?,?,?,?,?,?,?)")
	if err != nil {

		return nil, err
	}
	mId := uuid.New()
	defer stmt.Close()

	result, err := stmt.Exec(mId, m.Name, m.Region, m.Government, calculations,m.CalcType, time.Now(), email )

	if err != nil {

		return nil, err
	}
	idRes, _ := result.LastInsertId()

	idStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
	idStmt.Scan(&id, &dbRecords)

	json.Unmarshal([]byte(dbRecords), &records)
	records[mId.String()] = UserMtest{mId.String(), m.Name, m.Region,
		m.Government, m.CalcType, "", "", nil}

	out, err := json.Marshal(records)
	check(err)

	err = updateUser("records", string(out), id)

	if err != nil {
		return nil, err
	}
//need check
	err = updateIndex(idRes)
	if err != nil {
		return nil, err
	}



	return &records, nil
}

func readMtest(id string) (*Mtest, error) {
	var (
		mtest Mtest
		mid uuid.UUID
		rowId, calcType int
		name, govern, region, calculations, calcData, executors, pubDate, author string
	)
	res := db.QueryRow("SELECT m.id, m.mid, m.name, r.reg_name, g.gov_name, m.calculations," +
		" m.calc_type, m.calc_data, m.executors, m.pub_date, m.author FROM mtests m JOIN govs g ON m.govern = g.id JOIN " +
			"regions r ON m.region = r.reg_id WHERE mid=?", id)

	err := res.Scan(&rowId, &mid, &name, &region, &govern, &calculations, &calcType, &calcData, &executors, &pubDate, &author)
	if err != nil {

		return nil, err
	}

	mtest = Mtest{rowId, mid, name, region,
	govern,calculations,calcType,calcData, executors,pubDate, author}

	return &mtest, nil
}

func updateMtest(m map[string]interface{}, email string) error {
	var (
		id        int
		dbRecords string
		records   map[string]UserMtest
	)

	if m["calculations"] == nil && m["name"] != nil {

		stmt, err := db.Prepare("UPDATE mtests SET name=?, region=?, govern=? WHERE mid=?;")
		defer stmt.Close()
		if err != nil {
			return err
		} else if m["calculations"] != nil && m["name"] == nil {
		stmt, err := db.Prepare("UPDATE mtests SET calculations= ? WHERE mid=?;")
		if err != nil {
			return err
		} else {
			_, err := stmt.Exec(m["calculations"], m["id"])
			check(err)
			return nil
		}
	} else if m["executors"] != nil && m["name"] == nil {

		stmt, err := db.Prepare("UPDATE mtests SET executors= ? WHERE mid=?;")
		if err != nil {
			return err
		} else {
			_, err := stmt.Exec(m["executors"], m["id"])
			check(err)
			return nil
		}
	} else {
			govern := int(m["govern"].(float64))
			region := int(m["region"].(float64))
			calcType := int(m["calc_type"].(float64))

			_, err := stmt.Exec(m["name"], region, govern, m["mid"])
			check(err)

			idStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
			idStmt.Scan(&id, &dbRecords)

			json.Unmarshal([]byte(dbRecords), &records)
			record := records[m["mid"].(string)]
			records[m["mid"].(string)] = UserMtest{m["mid"].(string), m["name"].(string),
				region, govern, calcType,
				"", "", record.Executors}

			out, err := json.Marshal(records)
			check(err)

			idErr := updateUser("records", string(out), id)

			if idErr != nil {
				return idErr
			}

			return nil
		}
	}
	//update index!!!!!!!!
	return nil
}

func deleteMtest(mid, email string) error {
	var (
		id int
		dbRecords string
		records map[string]interface{}
	)
	if stmt, err := db.Prepare("DELETE FROM mtests WHERE mid=?"); err != nil {

		defer stmt.Close()
		return err
	} else {
		defer stmt.Close()
		if _, err := stmt.Exec(mid); err != nil {
			return err
		}
	}
	idStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
	idStmt.Scan(&id, &dbRecords)

	json.Unmarshal([]byte(dbRecords), &records)

    delete(records,mid)

	out, err := json.Marshal(records)
	check(err)

	if err := updateUser("records", string(out), id); err == nil {
		return nil
	} else {
		return err
	}

}

func getGovs() (*[]governmentRegion, error){
	var (
		govs []governmentRegion
		govId int
		govName string
	)
	res, err := db.Query("SELECT gov_id, gov_name FROM govs")
	check(err)
	defer res.Close()

	for res.Next() {
		err = res.Scan(&govId, &govName)
		check(err)

		govs = append(govs, governmentRegion{govId, govName })
	}
	return &govs, nil
}

func getRegs() (*[]governmentRegion, error) {
	var (
		regions []governmentRegion
		regId int
		regName string
	)

	res, err := db.Query("SELECT reg_id, reg_name FROM regions")
	check(err)
	defer res.Close()

	for res.Next() {
		err = res.Scan(&regId, &regName)
		check(err)

		regions = append(regions, governmentRegion{regId, regName})
	}
	return &regions, nil
}

func getAdmactions() (*[]AdmAction, error) {
	var (
		actId int
		actName string
		actions []AdmAction
	)

	res, err := db.Query("SELECT act_id, act_name FROM adm_actions")
	defer res.Close()
	check(err)
	for res.Next() {
		err := res.Scan(&actId, &actName)
		check(err)
		action := AdmAction{actId, actName}
		actions = append(actions, action)
	}

	return &actions, nil

}

func createMtestExecutor(email string, ex newExecutor) (*uuid.UUID, error) {
	var mId = uuid.New()

	var (
		id, devId int
		dbRecords, devDbRecords,getExecutors string
		records map[string]interface{}
		devRecords map[string]UserMtest
	)


	//check if user exists !!!!

	if isUsernameAvailable(ex.Email) == true {
		return nil, errors.New("користувач не зареєстрований")
	}

	//add mtest type 3
	stmt, _ := db.Prepare("INSERT INTO mtests (mid, name, region, govern," +
		" calculations, calc_type, developer, dev_mid, pub_date, author) VALUES (?,?,?,?,?,?,?,?,?,?)")

	defer stmt.Close()

	_, err := stmt.Exec(mId, ex.Title, ex.Region, ex.Government,
		calculations, 3, email, ex.DevMid, time.Now(), ex.Email)
	if err != nil {

		return nil, err
	}


	//UPDATE MAIN MTEST executors!!!!!!!!!!!

	var saveExecutors = map[string]newExecutors{}

	exStmt := db.QueryRow("SELECT executors FROM mtests WHERE mid=?", ex.DevMid)
	exStmt.Scan(&getExecutors)
	json.Unmarshal([]byte(getExecutors), &saveExecutors)

	saveExecutors[mId.String()] = newExecutors{ex.Email, mId.String(), true}

	updOut, updOutErr := json.Marshal(saveExecutors)
	if updOutErr != nil {
		return nil, updOutErr
	}

	updMtest, updMtestErr := db.Prepare("UPDATE mtests SET executors=? WHERE mid=?")
	if updMtestErr != nil {
		return nil, updMtestErr
	}
	_, updError := updMtest.Exec(updOut,ex.DevMid)
	if updError != nil {
		return nil, updError
	}



	//add mtest to executor mtests
	idStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", ex.Email)
	idStmt.Scan(&id, &dbRecords)

	json.Unmarshal([]byte(dbRecords), &records)
	records[mId.String()] = UserMtest{mId.String(), ex.Title, ex.Region,
	ex.Government, 3, email, ex.DevMid, nil}

	out, marshErr := json.Marshal(records)
	if marshErr != nil {
		return nil, marshErr
	}

	idErr := updateUser("records", string(out), id)
	if idErr != nil {
		return nil, idErr
	}

	// add mtest to developer mtest to executors section
	devStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
	devStmt.Scan(&devId, &devDbRecords)

	json.Unmarshal([]byte(devDbRecords), &devRecords)
	record := devRecords[ex.DevMid]
	devRecords[ex.DevMid] = UserMtest{ex.DevMid, record.Name, record.Region,
		record.Government, record.CalcType, "", "", saveExecutors}

	devOut, devOutErr := json.Marshal(devRecords)
	if devOutErr != nil {
		return nil, devOutErr
	}

	devErr := updateUser("records", string(devOut), devId)
	if devErr != nil {
		return nil, devErr
	}

	return &mId, nil
}

func deleteMtestExecutor(devEmail string, del delExecutorReq) error {
	var (
		id int
		dbRecords string
		records map[string]interface{}
		devRecords map[string]UserMtest
	)
	//delete mtest

	stmt, _ := db.Prepare("DELETE FROM mtests WHERE mid=?")
	_, err := stmt.Exec(del.ExMtestId)
	if err != nil {
		return err
	}

	//delete from executors mtests
	idStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", del.ExEmail)
	idStmt.Scan(&id, &dbRecords)

	json.Unmarshal([]byte(dbRecords), &records)

	delete(records, del.ExMtestId)

	out, err := json.Marshal(records)
	check(err)

	if err := updateUser("records", string(out), id); err == nil {

	} else {
		return err
	}

	//delete from developers mtests

	devStmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", devEmail)
	devStmt.Scan(&id, &dbRecords)

	unmarshErr := json.Unmarshal([]byte(dbRecords), &devRecords)
	if unmarshErr != nil {
		return unmarshErr
	}
	delete(devRecords[del.DevMtestId].Executors, del.ExMtestId)
	devOut, devErr := json.Marshal(devRecords)
	if devErr != nil {
		return devErr
	}

	if err := updateUser("records", string(devOut), id); err == nil {

	} else {
		return err
	}

	//delete from developer mtest
	mtStmt := db.QueryRow("SELECT executors FROM mtests WHERE mid=?", del.DevMtestId)
	mtStmt.Scan(&dbRecords)

	json.Unmarshal([]byte(dbRecords), &records)
	delete(records, del.ExMtestId)
	mtOut, err := json.Marshal(records)
	check(err)

	mtSaveStmt, _ := db.Prepare("UPDATE mtests SET executors=? WHERE mid=?", )
	_, mtErr := mtSaveStmt.Exec(mtOut, del.DevMtestId)
	if mtErr != nil {
		return mtErr
	}

	return nil

}