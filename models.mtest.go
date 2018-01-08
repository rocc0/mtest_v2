package main

import (
	"log"
	"time"
	"github.com/google/uuid"
	"encoding/json"
)

type Gov struct {
	Id int `json:"id"`
	Name string	`json:"name"`
}

type Mtest struct {
	Id int `json:"id"`
	Mid uuid.UUID `json:"mid"`
	Name string `json:"name"`
	Region int `json:"region"`
	Govern int `json:"govern"`
	Calculations string `json:"calculations"`
	Corruptions string `json:"corruptions"`
	PubDate string `json:"pub_date"`
	Author string `json:"author"`
}

type UserMtest struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Region int `json:"region"`
	Govern int `json:"govern"`
}

const corr = `{}`
const calcs = `{"1":[{"type":"container","id":3,"columns":[[{"type":"itemplus","id":3,
                    "columns":[[{"type":"item","id":3,"name":"Додати дію","subsum":0},{"type":"item","id":6,"name":"Додати дію","subsum":0}]],
                    "name":"Додати складову інф. вимоги"}]],"name":"Додати інф. вимогу","contsub":0},
                {"type":"container","id":null,"columns":[[{"type":"itemplus","id":4,"columns":[[{"type":"item","id":3,"name":"Додати дію","subsum":0},
                            {"type":"item","id":4,"name":"Додати дію","subsum":0}]],"name":"Додати складову інф. вимоги"}]],"name":"Додати інф. вимогу","contsub":0}]}`



func createNewMtest(m newMtest, email string) (*map[string]interface{}, error) {
	var (
		id int
		db_records string
		records map[string]interface{}
		)

	stmt, err := db.Prepare("INSERT INTO mtests (mid, name, region, govern," +
		" calculations, corruptions, pub_date, author) VALUES (?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Print(err)
		return nil, err
	}
	m_id := uuid.New()
	result, err := stmt.Exec(m_id, m.Name, m.Region, m.Government, calcs, corr, time.Now(), email)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	id_res, _ := result.LastInsertId()

	id_stmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
	id_stmt.Scan(&id, &db_records)

	json.Unmarshal([]byte(db_records), &records)
	records[m_id.String()] = UserMtest{m_id.String(), m.Name, m.Region, m.Government}

	out, err := json.Marshal(records)
	check(err)

	id_err := updateUser("records", string(out), id)

	if id_err != nil {
		return nil, id_err
	}
//need check
	idx_err := updateIndex(id_res)
	if idx_err != nil {
		return nil, idx_err
	}


	return &records, nil
}

func readMtest(id string) (*Mtest, error) {
	var (
		mtest Mtest
		mid uuid.UUID
		row_id, govern, region int
		name,  calculations, corruptions, pub_date, author string
	)
	res := db.QueryRow("SELECT id, mid, name, region, govern, calculations," +
		" corruptions, pub_date, author FROM mtests WHERE mid=?", id)

	err := res.Scan(&row_id, &mid, &name, &region, &govern, &calculations, &corruptions, &pub_date, &author)
	if err != nil {
		log.Print(err, "kek")
		return nil, err
	}

	mtest = Mtest{row_id, mid, name, region, govern,calculations,
	corruptions,pub_date, author}

	return &mtest, nil
}

func updateMtest(m map[string]interface{}, email string) error{
	var (
		id int
		db_records string
		records map[string]interface{}
	)

	if m["corruptions"] == nil && m["name"] != nil {

		stmt, err := db.Prepare("UPDATE mtests SET name=?, region=?, govern=? WHERE mid=?;")
		if err != nil {
			return err
		} else {
			region := int(m["region"].(float64))
			govern := int(m["govern"].(float64))

			_, err := stmt.Exec(m["name"],region,govern,m["mid"])
			check(err)

			id_stmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
			id_stmt.Scan(&id, &db_records)

			json.Unmarshal([]byte(db_records), &records)

			records[m["mid"].(string)] = UserMtest{m["mid"].(string),m["name"].(string),
				region,govern}

			out, err := json.Marshal(records)
			check(err)

			id_err := updateUser("records", string(out), id)

			if id_err != nil {
				return id_err
			}

			return nil
		}
	} else if m["corruptions"] != nil && m["name"] == nil {
		stmt, err := db.Prepare("UPDATE mtests SET calculations= ?, corruptions=?, WHERE id=?;")
		if err != nil {
			return err
		} else {
			_, err := stmt.Exec(m["calculations"],m["corruptions"],m["id"])
			check(err)
			return nil
		}
	}

	//update index!!!!!!!!
	return nil
}

func deleteMtest(mid, email string) error {
	var (
		id int
		db_records string
		records map[string]interface{}
	)
	if stmt, err := db.Prepare("DELETE FROM mtests WHERE mid=?"); err != nil {
		log.Print("\n",err,mid,"\n")
		return err
	} else {
		if res, err := stmt.Exec(mid); err != nil {
			log.Print("\n",err,res,"\n")
			return err
		}
	}
	id_stmt := db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
	id_stmt.Scan(&id, &db_records)

	json.Unmarshal([]byte(db_records), &records)

    delete(records,mid)

	out, err := json.Marshal(records)
	check(err)

	if err := updateUser("records", string(out), id); err == nil {
		return nil
	} else {
		return err
	}

}

func editGovName(id int, name string) error {
	stmt, err := db.Prepare("UPDATE governments SET gov_name=? WHERE id=?;")
	check(err)
	_, err = stmt.Exec(name, id)
	if err != nil {
		return err
	}
	return nil
}


func getGovs() (*[]Gov, error){
	var (
		govs []Gov
		gov_id int
		gov_name string
	)
	res, err := db.Query("SELECT id, gov_name FROM governments")
	check(err)

	for res.Next() {
		err = res.Scan(&gov_id, &gov_name)
		check(err)

		govs = append(govs, Gov{gov_id, gov_name })
	}
	return &govs, nil
}