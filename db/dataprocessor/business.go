package dataprocessor

import log "github.com/sirupsen/logrus"

func (mt *Service) AddBusiness(name string) error {
	stmt, err := mt.db.Prepare("INSERT INTO businesses (name) VALUES (?) ;")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()

	if _, err = stmt.Exec(name); err != nil {
		return err
	}
	return nil
}

func (mt *Service) GetBusinesses() ([]Government, error) {
	var (
		govs    []Government
		govId   int
		govName string
	)
	res, err := mt.db.Query("SELECT id, name FROM businesses")
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
	return govs, nil
}

func (mt *Service) EditBusinessName(id int, name string) error {
	stmt, err := mt.db.Prepare("UPDATE businesses SET name=? WHERE id=?;")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()

	if _, err = stmt.Exec(name, id); err != nil {
		return err
	}
	return nil
}

func (mt *Service) RemoveBusiness(id int) error {
	stmt, err := mt.db.Prepare("DELETE FROM businesses WHERE id=?;")
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
