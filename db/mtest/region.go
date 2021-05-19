package mtest

import "github.com/sirupsen/logrus"

type Region struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

const (
	editRegNameQuery = `UPDATE regions SET reg_name=? WHERE reg_id=?;`
	getRegionsQuery  = `SELECT reg_id, reg_name FROM regions`
)

func (mt *Service) GetRegions() (*[]Region, error) {
	var (
		regions []Region
		regId   int
		regName string
	)

	res, err := mt.db.Query(getRegionsQuery)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Close(); err != nil {
			logrus.Error(err)
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

func (mt *Service) EditRegionName(id int, name string) error {
	stmt, err := mt.db.Prepare(editRegNameQuery)
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	if _, err = stmt.Exec(name, id); err != nil {
		return err
	}
	return nil
}
