package v2

import (
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	"encoding/json"

	"github.com/sirupsen/logrus"
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

const createUsersTable = `CREATE TABLE IF NOT EXISTS
    users (id SERIAL NOT NULL PRIMARY KEY, name VARCHAR(100) NOT NULL, 
    surename VARCHAR(20) NOT NULL, email VARCHAR(100), password VARCHAR(100) NOT NULL, 
    rights VARCHAR(100) NOT NULL, records VARCHAR(100) NOT NULL, activated INTEGER DEFAULT 0);`

func userInit() error {
	stmt, err := db.Prepare(createUsersTable)
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

func passwordCheck(email, password string) bool {
	var passw string
	res := db.QueryRow("SELECT password FROM users WHERE email=?", email)
	if err := res.Scan(&passw); err != nil {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passw), []byte(password)); err != nil {
		return false
	}
	return true
}

const checkActivationQuery = "SELECT activated FROM users WHERE email=?"

func checkUserActivated(email string) bool {
	var (
		activated int
	)

	res := db.QueryRow(checkActivationQuery, email)
	if err := res.Scan(&activated); err != nil {
		log.Error(err)
		return false
	}

	return activated == 1
}

const createUserQuery = "INSERT INTO users (name, surename, email, records, password) VALUES (?,?,?,?,?)"

func (u *User) createUser() (*string, error) {
	if checkUserExists(u.Email) == false {
		return &u.Name, errors.New("користувач з цим ім'ям вже існує")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	req, err := db.Prepare(createUserQuery)
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

	hash, err := u.writeHash()
	if err != nil {
		return nil, err
	}

	if err = doSendEmail(*u, hash, "email_activate"); err != nil {
		return nil, err
	}

	return &u.Name, nil
}

const getUserQuery = "SELECT name, surename, email, id, rights, records FROM users WHERE email = ?"

func getUser(email string) (*User, error) {
	var records string
	user := &User{}

	res := db.QueryRow(getUserQuery, email)
	if err := res.Scan(&user.Name, &user.Surename, &user.Email, &user.Id, &user.Rights, &records); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(records), &user.Records); err != nil {
		return nil, err
	}

	return user, nil
}

func updateUser(field, data string, id int) error {
	stmt, err := db.Prepare("UPDATE users SET " + field + "=? WHERE id=?;")
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			logrus.Error(err)
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

func checkUserExists(email string) bool {
	var result string
	res := db.QueryRow("SELECT email FROM users WHERE email=?", email)
	if err := res.Scan(&result); err != nil {
		logrus.Error(err)
		return false
	}
	return result == ""
}
