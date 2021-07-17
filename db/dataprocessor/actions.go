package dataprocessor

import (
	"math/rand"

	log "github.com/sirupsen/logrus"
)

func (mt *Service) GetAdministrativeActions() (*[]AdmAction, error) {
	var (
		actId   int
		actName string
		actions []AdmAction
	)

	res, err := mt.db.Query("SELECT id, act_name FROM adm_actions")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Close(); err != nil {
			log.Error(err)
		}
	}()
	for res.Next() {
		if err := res.Scan(&actId, &actName); err != nil {
			return nil, err
		}
		actions = append(actions, AdmAction{ActId: actId, ActName: actName})
	}

	return &actions, nil

}

func (mt *Service) DeleteAdministrativeAction(id int) error {
	stmt, err := mt.db.Prepare("DELETE FROM adm_actions WHERE id=?;")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()

	if _, err = stmt.Exec(id); err != nil {
		return err
	}
	return nil
}

func (mt *Service) EditAdministrativeActionName(id int, name string) error {
	stmt, err := mt.db.Prepare("UPDATE adm_actions SET act_name=? WHERE id=?;")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()
	log.Error(name, id)
	if _, err = stmt.Exec(name, id); err != nil {
		return err
	}
	return nil
}

func (mt *Service) AddAdministrativeAction(name string) error {
	stmt, err := mt.db.Prepare("INSERT INTO adm_actions (act_name,act_id) VALUES (?,?) ;")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()

	if _, err = stmt.Exec(name, rand.Int31()); err != nil {
		return err
	}
	return nil
}
