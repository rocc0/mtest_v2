package mtest

import log "github.com/sirupsen/logrus"

type Region struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (mt *Service) getRegions() (*[]Region, error) {
	var (
		regions []Region
		regId   int
		regName string
	)

	res, err := mt.db.Query("SELECT reg_id, reg_name FROM regions")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Close(); err != nil {
			log.Error(err)
		}
	}()

	for res.Next() {
		if err = res.Scan(&regId, &regName); err != nil {
			return nil, err
		}
		regions = append(regions, Region{Id: regId, Name: regName})
	}
	return &regions, nil
}
