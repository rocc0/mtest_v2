package dataprocessor

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const insertRegActQuery = "INSERT INTO reg_acts (mid, doc_id, doc_text, doc_name, doc_type) VALUES (?, ?, ?, ?, ?);"
const deleteRegActQuery = "DELETE FROM reg_acts WHERE mid=? AND doc_id=?;"
const getRegActQuery = "SELECT doc_text, doc_name, doc_type FROM reg_acts WHERE mid=? AND doc_id=?;"
const listRegActsQuery = `SELECT doc_id, doc_name, doc_type FROM reg_acts WHERE mid=?;`

func (mt *Service) InsertRegAct(mtestID string, docText string, docName string, docType string) (string, error) {
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

	if _, err = stmt.Exec(mtestID, docID, docText, docName, docType); err != nil {
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

	if err := res.Scan(&doc.Text, &doc.Name, &doc.Type); err != nil {
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

	res, err := mt.db.Query(listRegActsQuery)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Close(); err != nil {
			log.Error(err)
		}
	}()
	for res.Next() {
		var docID, docName, docType string
		if err := res.Scan(&docID, &docName, &docType); err != nil {
			return nil, err
		}

		users = append(users, RegAct{
			MtestID: mtestID,
			DocID:   docID,
			Name:    docName,
			Text:    docType,
		})
	}

	return users, nil
}
