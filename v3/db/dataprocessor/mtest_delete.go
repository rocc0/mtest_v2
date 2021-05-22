package dataprocessor

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

func (mt *Service) DeleteMTEST(mid, email string) error {
	var (
		id        int
		dbRecords string
		records   map[string]interface{}
	)
	stmt, err := mt.db.Prepare("DELETE FROM mtests WHERE mid=?")
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

	idStmt := mt.db.QueryRow("SELECT id, records FROM users WHERE email=?", email)
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

	return mt.UpdateUser("records", string(out), id)
}
