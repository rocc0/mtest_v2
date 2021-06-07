package dataprocessor

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type DeleteExecutor struct {
	ExEmail    string `json:"ex_email"`
	ExMtestId  string `json:"ex_mtest_id"`
	DevMtestId string `json:"dev_mtest_id"`
}

func (mt *Service) DeleteExecutor(devEmail string, del DeleteExecutor) error {
	var (
		id         int
		dbRecords  string
		records    map[string]interface{}
		devRecords map[string]MTestData
	)
	//delete mtest

	stmt, err := mt.db.Prepare("DELETE FROM mtests WHERE mid=?")
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
	if err := mt.db.QueryRow("SELECT id, records FROM users WHERE email=?", del.ExEmail).Scan(&id, &dbRecords); err != nil {
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

	if err := mt.UpdateUser("records", string(out), id); err != nil {
		return err
	}

	//delete from developers mtests
	if err := mt.db.QueryRow("SELECT id, records FROM users WHERE email=?", devEmail).Scan(&id, &dbRecords); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(dbRecords), &devRecords); err != nil {
		return err
	}

	if _, ok := devRecords[del.DevMtestId]; ok {
		delete(devRecords[del.DevMtestId].Executors, del.ExMtestId)
	}

	devOut, err := json.Marshal(devRecords)
	if err != nil {
		return err
	}

	if err := mt.UpdateUser("records", string(devOut), id); err != nil {
		return err
	}

	//delete from developer mtest
	if err := mt.db.QueryRow("SELECT executors FROM mtests WHERE mid=?", del.DevMtestId).Scan(&dbRecords); err != nil {
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

	mtSaveStmt, err := mt.db.Prepare("UPDATE mtests SET executors=? WHERE mid=?")
	if err != nil {
		return err
	}

	if _, err := mtSaveStmt.Exec(mtOut, del.DevMtestId); err != nil {
		return err
	}
	return nil

}
