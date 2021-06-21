package dataprocessor

import log "github.com/sirupsen/logrus"

const initSynonyms = `CREATE TABLE IF NOT EXISTS synonyms
(
	id bigint auto_increment,
	word varchar(300) not null,
	synonym varchar(300) not null,
	constraint synonyms_pk primary key (id)
);
`

type Synonym struct {
	Word     string   `json:"word"`
	Synonyms []string `json:"synonyms"`
}

func (mt *Service) Load() ([]Synonym, error) {
	var (
		word, synonym string
	)

	synonyms := make(map[string][]string)
	result := []Synonym{}

	res, err := mt.db.Query("SELECT word, synonym FROM synonyms")
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
		result = append(result, Synonym{
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

func (mt *Service) AddSynonym(word, syn string) error {
	stmt, err := mt.db.Prepare("INSERT INTO synonyms (word, synonym) VALUES (?,?) ;")
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

func (mt *Service) RemoveSynonym(word, syn string) error {
	stmt, err := mt.db.Prepare("DELETE FROM synonyms WHERE word=? and synonym=?;")
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
