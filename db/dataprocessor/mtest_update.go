package dataprocessor

import (
	"database/sql"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	db *sql.DB
}

type (
	MTestData struct {
		Id         string                  `json:"id"`
		Name       string                  `json:"name"`
		Region     int                     `json:"region"`
		Government int                     `json:"govern"`
		CalcType   int                     `json:"calc_type"`
		Developer  string                  `json:"developer"`
		DevMid     string                  `json:"dev_mid"`
		Business   int                     `json:"business,omitempty"`
		Executors  map[string]ExecutorInfo `json:"executors"`
		Files      []RegAct                `json:"files"`
		Synonyms   []Synonym               `json:"synonyms"`
	}
	AdmAction struct {
		ActId   int    `json:"id"`
		ActName string `json:"act_name"`
	}
	ExecutorInfo struct {
		Email   string `json:"email"`
		Mid     string `json:"mid"`
		Checked bool   `json:"checked"`
	}
)

func (mt *Service) UpdateMTEST(m map[string]interface{}, email string) error {
	if m["calculations"] != nil && m["name"] == nil {
		return mt.updateCalculations(m)
	} else if m["executors"] != nil && m["name"] == nil {
		return mt.updateExecutors(m)
	}

	return mt.updateMTESTAndUser(m, email)
}

func (mt *Service) updateCalculations(m map[string]interface{}) error {
	stmt, err := mt.db.Prepare("UPDATE mtests SET calculations= ?, corr_result=?, math_result=? WHERE mid=?;")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()
	if _, err := stmt.Exec(m["calculations"], m["corr_total"], m["calc_total"], m["id"]); err != nil {
		return err
	}
	return nil
}

func (mt *Service) updateExecutors(m map[string]interface{}) error {
	stmt, err := mt.db.Prepare("UPDATE mtests SET executors= ? WHERE mid=?;")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()
	if _, err := stmt.Exec(m["executors"], m["id"]); err != nil {
		return err
	}
	return nil
}

func (mt *Service) updateMTESTAndUser(m map[string]interface{}, email string) error {
	stmt, err := mt.db.Prepare("UPDATE mtests SET name=?, region=?, govern=? WHERE mid=?;")
	if err != nil {
		return err
	}
	var (
		id        int
		dbRecords string
		records   map[string]MTestData
	)

	govern := int(m["govern"].(float64))
	region := int(m["region"].(float64))
	calcType := int(m["calc_type"].(float64))
	if _, err := stmt.Exec(m["name"], region, govern, m["mid"]); err != nil {
		return err
	}

	idStmt := mt.db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
	if err := idStmt.Scan(&id, &dbRecords); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(dbRecords), &records); err != nil {
		return err
	}
	record := records[m["mid"].(string)]
	records[m["mid"].(string)] = MTestData{Id: m["mid"].(string), Name: m["name"].(string),
		Region: region, Government: govern, CalcType: calcType, Executors: record.Executors}

	out, err := json.Marshal(records)
	if err != nil {
		return err
	}
	return mt.UpdateUser("records", string(out), id)
}
