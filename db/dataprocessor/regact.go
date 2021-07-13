package dataprocessor

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const insertRegActQuery = "INSERT INTO reg_acts (mid, doc_id, doc_text, doc_name) VALUES (?, ?, ?, ?);"
const deleteRegActQuery = "DELETE FROM reg_acts WHERE mid=? AND doc_id=?;"
const getRegActQuery = "SELECT doc_text, doc_name FROM reg_acts WHERE mid=? AND doc_id=?;"
const listRegActsQuery = `SELECT doc_id, doc_name FROM reg_acts WHERE mid=?;`

func (mt *Service) InsertRegAct(mtestID string, docText string, docName string) (string, error) {
	stmt, err := mt.db.Prepare(insertRegActQuery)
	if err != nil {
		return "", err
	}
	docID := uuid.New().String()
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()

	if _, err = stmt.Exec(mtestID, docID, docText, docName); err != nil {
		return "", err
	}
	return docID, nil
}

func (mt *Service) DeleteRegAct(mtestID string, docID string) error {
	stmt, err := mt.db.Prepare(deleteRegActQuery)
	if err != nil {
		return err

	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()
	if _, err := stmt.Exec(mtestID, docID); err != nil {
		return err
	}

	return nil
}

type RegAct struct {
	MtestID string
	DocID   string
	Name    string
	Text    string
	Type    string
}

func (mt *Service) GetRegAct(mtestID string, docID string) (RegAct, error) {
	var doc RegAct
	res := mt.db.QueryRow(getRegActQuery, mtestID, docID)

	if err := res.Scan(&doc.Text, &doc.Name); err != nil {
		return doc, err
	}

	doc.DocID = docID
	doc.MtestID = mtestID

	return doc, nil
}

func (mt *Service) ListRegActs(mtestID string) ([]RegAct, error) {
	var (
		users []RegAct
	)

	res, err := mt.db.Query(listRegActsQuery, mtestID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Close(); err != nil {
			log.Error(err)
		}
	}()
	for res.Next() {
		var docID, docName string
		if err := res.Scan(&docID, &docName); err != nil {
			return nil, err
		}

		users = append(users, RegAct{
			MtestID: mtestID,
			DocID:   docID,
			Name:    docName,
		})
	}

	return users, nil
}
