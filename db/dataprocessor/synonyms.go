package dataprocessor

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const initSynonyms = `CREATE TABLE IF NOT EXISTS synonyms
(
	id bigint auto_increment,
	mtest_id varchar(300) not null,
	synonym varchar(300) not null,
	synonym_id varchar(300) not null,
	constraint synonyms_pk primary key (id)
);
`

const listSynonymsQuery = `SELECT synonym_id, synonym FROM synonyms WHERE mtest_id=?;`

type Synonym struct {
	SynonymID string `json:"synonym_id"`
	Synonym   string `json:"synonym"`
}

func (mt *Service) GetSynonymsByID(mtestID string) ([]Synonym, error) {
	var synonyms []Synonym
	res, err := mt.db.Query(listSynonymsQuery, mtestID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Close(); err != nil {
			log.Error(err)
		}
	}()
	for res.Next() {
		var s Synonym
		if err := res.Scan(&s.SynonymID, &s.Synonym); err != nil {
			return nil, err
		}

		synonyms = append(synonyms, s)
	}

	return synonyms, nil
}

func (mt *Service) AddSynonym(mtestID, syn string) (string, error) {
	synonymID := uuid.New().String()
	stmt, err := mt.db.Prepare("INSERT INTO synonyms (mtest_id, synonym, synonym_id) VALUES (?,?,?) ;")
	if err != nil {
		return synonymID, err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()

	if _, err = stmt.Exec(mtestID, syn, synonymID); err != nil {
		return synonymID, err
	}
	return synonymID, nil
}

func (mt *Service) RemoveSynonym(mtestID, synonymID string) error {
	stmt, err := mt.db.Prepare("DELETE FROM synonyms WHERE mtest_id=? and synonym_id=?;")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()

	if _, err = stmt.Exec(mtestID, synonymID); err != nil {
		return err
	}
	return nil
}
