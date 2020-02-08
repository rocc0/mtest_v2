package main

import (
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type User struct {
	Id       int                    `json:"id"`
	Name     string                 `json:"name"`
	Surename string                 `json:"surename"`
	Email    string                 `json:"email"`
	Rights   int                    `json:"rights"`
	Password string                 `json:"password"`
	Records  map[string]interface{} `json:"records"`
}

func userInit() error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id SERIAL NOT NULL PRIMARY KEY, " +
		"name VARCHAR(100) NOT NULL, surename VARCHAR(20) NOT NULL, email VARCHAR(100)," +
		" password VARCHAR(100) NOT NULL, rights VARCHAR(100) NOT NULL, records VARCHAR(100) NOT NULL);")
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()
	if _, err = stmt.Exec(); err != nil {
		return err
	}
	return nil
}

func loginCheck(email, password string) bool {
	var (
		eMail, passw string
	)
	res := db.QueryRow("SELECT email, password FROM users WHERE email=?", email)
	if err := res.Scan(&eMail, &passw); err != nil {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passw), []byte(password)); err != nil {
		return false
	}
	return true
}

func authCheck(email string) bool {
	var (
		privileged int
	)
	res := db.QueryRow("SELECT activated FROM users WHERE email=?", email)
	if err := res.Scan(&privileged); err != nil {
		log.Error(err)
		return false
	}

	return privileged == 1
}

func (u *User) createUser() (*string, error) {
	if isUsernameAvailable(u.Email) == false {
		return &u.Name, errors.New("користувач з цим ім'ям вже існує")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	req, err := db.Prepare("INSERT INTO users (name, surename, email, records, password) VALUES (?,?,?,?,?)")
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := req.Close(); err != nil {
			log.Error(err)
		}
	}()
	if _, err = req.Exec(u.Name, u.Surename, u.Email, "{}", hashedPassword); err != nil {
		return nil, err
	}

	usrHash, err := u.writeHash()
	if err != nil {
		return nil, err
	}

	if err = doSendEmail(*u, *usrHash, "email_activate"); err != nil {
		return nil, err
	}

	return &u.Name, nil
}

func readUser(email string) (*User, error) {
	log.Print(email)
	var (
		name, sureName, eMail string
		jsonRecords           map[string]interface{}
		records               string
		id, rights            int
	)

	res := db.QueryRow("SELECT name, surename, email, id, rights, records FROM users WHERE email = ?", email)
	if err := res.Scan(&name, &sureName, &eMail, &id, &rights, &records); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(records), &jsonRecords); err != nil {
		return nil, err
	}

	return &User{Id: id, Name: name, Surename: sureName,
		Email: eMail, Rights: rights, Records: jsonRecords}, nil
}

func updateUser(field, data string, id int) error {

	stmt, err := db.Prepare("UPDATE users SET " + field + "=? WHERE id=?;")
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()
	if field == "password" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
		if _, err = stmt.Exec(field, hashedPassword, id); err != nil {
			return err
		}
	}

	if _, err = stmt.Exec(data, id); err != nil {
		return err
	}
	return nil
}

func isUsernameAvailable(email string) bool {
	var result string
	res := db.QueryRow("SELECT email FROM users WHERE email=?", email)
	if err := res.Scan(&result); err != nil {
		return false
	}
	if result == "" {
		return true
	}
	return false
}
