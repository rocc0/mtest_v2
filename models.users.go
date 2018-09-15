package main

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/go-sql-driver/mysql"

	"log"
	"encoding/json"
)

type User struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Surename string `json:"surename"`
	Email string `json:"email"`
	Rights int `json:"rights"`
	Password string `json:"password"`
	Records map[string]interface{} `json:"records"`
}

func userInit() {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id SERIAL NOT NULL PRIMARY KEY, " +
		"name VARCHAR(100) NOT NULL, surename VARCHAR(20) NOT NULL, email VARCHAR(100)," +
			" password VARCHAR(100) NOT NULL, rights VARCHAR(100) NOT NULL, records VARCHAR(100) NOT NULL);")

	defer stmt.Close()

	check(err)

	_, err = stmt.Exec()
	check(err)

}


func loginCheck(email, password string) bool{
	var (
		eMail, passw string
	)
	res := db.QueryRow("SELECT email, password FROM users WHERE email=?", email)
	res.Scan(&eMail, &passw)
	err := bcrypt.CompareHashAndPassword([]byte(passw), []byte(password))

	if err != nil {
		return false
	}
	return true
}

func authCheck(email string) bool {
	var (
		privileged int
	)
	res := db.QueryRow("SELECT activated FROM users WHERE email=?", email)
	res.Scan(&privileged)

	return privileged == 1
}

func (u *User) createUser() (*string, error) {
	if isUsernameAvailable(u.Email) == false {
		return &u.Name, errors.New("користувач з цим ім'ям вже існує")
	} else {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		req, err := db.Prepare("INSERT INTO users (name, surename, email, records, password) VALUES (?,?,?,?,?)")

		if err != nil {
			return nil, err
		}
		defer req.Close()
		_, err = req.Exec(u.Name, u.Surename, u.Email,"{}", hashedPassword)

		if err != nil {
			return nil, err
		}

		usrHash, err := u.writeHash()
		if err != nil {
			return nil, err
		}

		err = doSendEmail(*u, *usrHash, "email_activate")
		if err != nil {
			return nil, err
		}

		return &u.Name, nil
	}
}

func readUser(email string) (*User, error) {
	log.Print(email)
	var (
		name, surename, eMail string
		jsonRecords map[string]interface{}
		records string
		id,rights int
	)

	res := db.QueryRow("SELECT name, surename, email, id, rights, records FROM users WHERE email = ?", email)
	err := res.Scan(&name, &surename, &eMail, &id, &rights, &records)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(records), &jsonRecords)
	userData := User{id, name, surename, eMail, rights, "", jsonRecords}

	return &userData, nil
}


func updateUser(field, data string, id int) error {

	stmt, err := db.Prepare("UPDATE users SET "+ field + "=? WHERE id=?;")
	check(err)

	defer stmt.Close()

	if field == "password" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
		_, err = stmt.Exec(field, hashedPassword, id)
		check(err)
		return nil
	} else {
		_, err = stmt.Exec(data, id)
		check(err)
		return nil
	}
}

func isUsernameAvailable(email string) bool {
	var result string
	res := db.QueryRow("SELECT email FROM users WHERE email=?", email)
	res.Scan(&result)
	if result == "" {
		return true
	}
	return false
}