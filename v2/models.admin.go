package v2

import log "github.com/sirupsen/logrus"

func editGovName(id int, name string) error {
	stmt, err := db.Prepare("UPDATE governments SET gov_name=? WHERE id=?;")
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

func editRegName(id int, name string) error {
	stmt, err := db.Prepare("UPDATE governments SET gov_name=? WHERE id=?;")
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
