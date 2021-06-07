package dataprocessor

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Executor struct {
	Title      string `json:"title"`
	Email      string `json:"email"`
	Region     int    `json:"region"`
	Government int    `json:"government"`
	DevMid     string `json:"dev_mid"`
}

func (mt *Service) CreateExecutor(email string, ex Executor) (string, error) {
	var (
		id, devId                             int
		dbRecords, devDbRecords, getExecutors string
		records                               map[string]interface{}
		devRecords                            map[string]MTestData
	)

	if ok := mt.CheckUserExists(ex.Email); !ok {
		return "", errors.New("користувач не зареєстрований")
	}

	//add mtest type 3
	stmt, err := mt.db.Prepare("INSERT INTO mtests (mid, name, region, govern," +
		" calculations, calc_type, developer, dev_mid, pub_date, author) VALUES (?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return "", err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()
	mtestID := uuid.New().String()
	if _, err := stmt.Exec(mtestID, ex.Title, ex.Region, ex.Government,
		defaultCalculations, 3, email, ex.DevMid, time.Now(), ex.Email); err != nil {
		return "", err
	}

	//UPDATE MAIN MTEST executors!!!!!!!!!!!
	saveExecutors := map[string]ExecutorInfo{}
	if err := mt.db.QueryRow("SELECT executors FROM mtests WHERE mid=?", ex.DevMid).Scan(&getExecutors); err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(getExecutors), &saveExecutors); err != nil {
		return "", err
	}

	saveExecutors[mtestID] = ExecutorInfo{ex.Email, mtestID, true}
	updOut, updOutErr := json.Marshal(saveExecutors)
	if updOutErr != nil {
		return "", updOutErr
	}

	updStmt, err := mt.db.Prepare("UPDATE mtests SET executors=? WHERE mid=?")
	if err != nil {
		return "", err
	}

	if _, err = updStmt.Exec(updOut, ex.DevMid); err != nil {
		return "", err
	}

	//add mtest to executor mtests
	if err := mt.db.QueryRow("SELECT id, records FROM users WHERE email=?", ex.Email).Scan(&id, &dbRecords); err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(dbRecords), &records); err != nil {
		return "", err
	}

	records[mtestID] = MTestData{Id: mtestID, Name: ex.Title, Region: ex.Region,
		Government: ex.Government, CalcType: 3, Developer: email, DevMid: ex.DevMid}

	out, err := json.Marshal(records)
	if err != nil {
		return "", err
	}

	if idErr := mt.UpdateUser("records", string(out), id); idErr != nil {
		return "", idErr
	}

	// add mtest to developer mtest to executors section
	if err := mt.db.QueryRow("SELECT id, records FROM users WHERE email=?", email).Scan(&devId, &devDbRecords); err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(devDbRecords), &devRecords); err != nil {
		return "", err
	}

	record := devRecords[ex.DevMid]
	devRecords[ex.DevMid] = MTestData{Id: ex.DevMid, Name: record.Name, Region: record.Region,
		Government: record.Government, CalcType: record.CalcType, Executors: saveExecutors}

	devOut, err := json.Marshal(devRecords)
	if err != nil {
		return "", err
	}

	return mtestID, mt.UpdateUser("records", string(devOut), devId)
}
