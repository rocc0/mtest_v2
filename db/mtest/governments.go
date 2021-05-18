package mtest

import log "github.com/sirupsen/logrus"

type Government struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (mt *Service) GetGovernments() (*[]Government, error) {
	var (
		govs    []Government
		govId   int
		govName string
	)
	res, err := mt.db.Query("SELECT gov_id, gov_name FROM govs")
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Close(); err != nil {
			log.Error(err)
		}
	}()
	for res.Next() {
		if err = res.Scan(&govId, &govName); err != nil {
			return nil, err
		}
		govs = append(govs, Government{Id: govId, Name: govName})
	}
	return &govs, nil
}
