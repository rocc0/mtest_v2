package mtest

import log "github.com/sirupsen/logrus"

func (mt *Service) GetAdministrativeActions() (*[]AdmAction, error) {
	var (
		actId   int
		actName string
		actions []AdmAction
	)

	res, err := mt.db.Query("SELECT act_id, act_name FROM adm_actions")
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
