package main


func editGovName(id int, name string) error {
	stmt, err := db.Prepare("UPDATE governments SET gov_name=? WHERE id=?;")
	defer stmt.Close()
	check(err)
	_, err = stmt.Exec(name, id)
	if err != nil {
		return err
	}
	return nil
}


func editRegName(id int, name string) error {
	stmt, err := db.Prepare("UPDATE governments SET gov_name=? WHERE id=?;")
	defer stmt.Close()
	check(err)
	_, err = stmt.Exec(name, id)
	if err != nil {
		return err
	}
	return nil
}