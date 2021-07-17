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

const initGlobalSynonyms = `CREATE TABLE IF NOT EXISTS synonyms
(
	id bigint auto_increment,
	word varchar(300) not null,
	synonym varchar(300) not null,
	constraint synonyms_pk primary key (id)
);
`

const listSynonymsQuery = `SELECT synonym_id, synonym FROM synonyms WHERE mtest_id=?;`

type Synonym struct {
	SynonymID string `json:"synonym_id"`
	Synonym   string `json:"synonym"`
}

type GlobalSynonym struct {
	Word     string
	Synonyms []string
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

func (mt *Service) LoadGlobals() ([]GlobalSynonym, error) {
	var (
		word, synonym string
	)

	synonyms := make(map[string][]string)
	var result []GlobalSynonym

	res, err := mt.db.Query("SELECT word, synonym FROM global_synonyms")
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Close(); err != nil {
			log.Error(err)
		}
	}()

	for res.Next() {
		if err = res.Scan(&word, &synonym); err != nil {
			return nil, err
		}
		if s, ok := synonyms[word]; ok {
			s = append(s, synonym)
			if !contains(word, synonyms[word]) {
				s = append(s, synonym)
			}
		} else {
			synonyms[word] = []string{synonym}
		}

		if s, ok := synonyms[synonym]; ok {
			if !contains(synonym, synonyms[synonym]) {
				s = append(s, word)
			}
		} else {
			synonyms[synonym] = []string{word}
		}
	}

	for k, s := range synonyms {
		result = append(result, GlobalSynonym{
			Word:     k,
			Synonyms: s,
		})
	}

	return result, nil
}

func contains(word string, arr []string) bool {
	for _, s := range arr {
		if s == word {
			return true
		}
	}
	return false

}

func (mt *Service) AddGlobalSynonym(word, syn string) error {
	stmt, err := mt.db.Prepare("INSERT INTO global_synonyms (word, synonym) VALUES (?,?) ;")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()

	if _, err = stmt.Exec(word, syn); err != nil {
		return err
	}
	return nil
}

func (mt *Service) RemoveGlobalSynonym(word, syn string) error {
	stmt, err := mt.db.Prepare("DELETE FROM global_synonyms WHERE word=? and synonym=?;")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()

	if _, err = stmt.Exec(word, syn); err != nil {
		return err
	}
	return nil
}
