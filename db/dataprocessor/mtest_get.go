package dataprocessor

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type MTEST struct {
	Id           int       `json:"id"`
	Mid          uuid.UUID `json:"mid"`
	Name         string    `json:"name"`
	Region       string    `json:"region"`
	Govern       string    `json:"govern"`
	Business     string    `json:"business"`
	Calculations string    `json:"calculations"`
	CalcType     int       `json:"calc_type"`
	CalcData     string    `json:"calc_data"`
	Executors    string    `json:"executors"`
	PubDate      string    `json:"pub_date"`
	Author       string    `json:"author"`
}

func (mt *Service) GetMTEST(id string) (*MTEST, error) {
	var mtest MTEST
	res := mt.db.QueryRow("SELECT m.id, m.mid, m.name, r.reg_name, g.gov_name, m.calculations,"+
		" m.calc_type, m.calc_data, m.executors, m.pub_date, m.author, b.name FROM mtests m JOIN govs g ON m.govern = g.id JOIN "+
		"regions r ON m.region = r.reg_id JOIN businesses b ON m.business = b.id WHERE mid=?", id)

	if err := res.Scan(&mtest.Id, &mtest.Mid, &mtest.Name, &mtest.Region, &mtest.Govern, &mtest.Calculations,
		&mtest.CalcType, &mtest.CalcData, &mtest.Executors, &mtest.PubDate, &mtest.Author, &mtest.Business); err != nil {
		log.Error(err)
		return nil, err
	}
	return &mtest, nil
}
