package dataprocessor

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type NewMTEST struct {
	Name       string `json:"name"`
	Region     int    `json:"region"`
	Government int    `json:"government"`
	CalcType   int    `json:"calc_type"`
}

const createMTESTQuery = `INSERT INTO mtests (mid, name, region, govern, calculations, calc_type, pub_date, author) VALUES (?,?,?,?,?,?,?,?)`
const defaultCalculations = `{"1":[{"type":"container","id":3,"columns":[[{"type":"itemplus","id":3,
                    "columns":[[{"type":"item","id":3,"name":"Додати дію","subsum":0},{"type":"item","id":6,"name":"Додати дію","subsum":0}]],
                    "name":"Додати складову інф. вимоги"}]],"name":"Додати інф. вимогу","contsub":0},
                {"type":"container","id":null,"columns":[[{"type":"itemplus","id":4,"columns":[[{"type":"item","id":3,"name":"Додати дію","subsum":0},
                            {"type":"item","id":4,"name":"Додати дію","subsum":0}]],"name":"Додати складову інф. вимоги"}]],"name":"Додати інф. вимогу","contsub":0}]}`

func (mt *Service) CreateMTEST(m NewMTEST, email string) (UserMtest, error) {
	var (
		id        int
		dbRecords string
	)

	stmt, err := mt.db.Prepare(createMTESTQuery)
	if err != nil {
		return UserMtest{}, err
	}

	mtestID := uuid.New().String()
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()

	result, err := stmt.Exec(mtestID, m.Name, m.Region, m.Government, defaultCalculations, m.CalcType, time.Now(), email)
	if err != nil {
		return UserMtest{}, err
	}

	idRes, err := result.LastInsertId()
	if err != nil {
		return UserMtest{}, err
	}
	idStmt := mt.db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
	if err := idStmt.Scan(&id, &dbRecords); err != nil {
		return UserMtest{}, err
	}

	var records map[string]interface{}
	if err := json.Unmarshal([]byte(dbRecords), &records); err != nil {
		return UserMtest{}, err
	}
	data := UserMtest{Id: mtestID, Name: m.Name, Region: m.Region, Government: m.Government, CalcType: m.CalcType, RecID: idRes}
	records[mtestID] = data

	out, err := json.Marshal(records)
	if err != nil {
		return UserMtest{}, err
	}

	if err = mt.UpdateUser("records", string(out), id); err != nil {
		return UserMtest{}, err
	}

	//need check
	return data, nil
}
