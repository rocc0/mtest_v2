package dataprocessor

import (
	"encoding/json"
	"errors"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	createUsersTable = `CREATE TABLE IF NOT EXISTS
    users (id SERIAL NOT NULL PRIMARY KEY, name VARCHAR(100) NOT NULL, 
    surename VARCHAR(20) NOT NULL, email VARCHAR(100), password VARCHAR(100) NOT NULL, 
    rights VARCHAR(100) NOT NULL, records VARCHAR(100) NOT NULL, activated INTEGER DEFAULT 0);`
	checkActivationQuery = `SELECT activated FROM users WHERE email=?`
	createUserQuery      = `INSERT INTO users (name, surename, email, records, password) VALUES (?,?,?,?,?)`
	getUserQuery         = `SELECT name, surename, email, id, rights, records FROM users WHERE email = ?`
	deleteUserQuery      = `DELETE FROM users WHERE id=?;`
	checkUserExistsQuery = `SELECT email FROM users WHERE email=?`
	passwordCheckQuery   = `SELECT password FROM users WHERE email=?`
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

func (mt *Service) InitUsersTable() error {
	stmt, err := mt.db.Prepare(createUsersTable)
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

func (mt *Service) PasswordCheck(email, password string) bool {
	var dbPassword string
	res := mt.db.QueryRow(passwordCheckQuery, email)
	if err := res.Scan(&dbPassword); err != nil {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err != nil {
		return false
	}
	return true
}

func (mt *Service) CheckUserActivation(email string) bool {
	var activated int

	res := mt.db.QueryRow(checkActivationQuery, email)
	if err := res.Scan(&activated); err != nil {
		log.Error(err)
		return false
	}

	return activated == 1
}

func (mt *Service) CreateUser() (string, error) {
	var u User
	if mt.CheckUserExists(u.Email) == false {
		return u.Name, errors.New("користувач з цим ім'ям вже існує")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	req, err := mt.db.Prepare(createUserQuery)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := req.Close(); err != nil {
			log.Error(err)
		}
	}()
	if _, err = req.Exec(u.Name, u.Surename, u.Email, "{}", hashedPassword); err != nil {
		return "", err
	}

	//todo move this part outside this library
	//hash, err := u.writeHash()
	//if err != nil {
	//	return "", err
	//}

	//if err = doSendEmail(u, hash, "email_activate"); err != nil {
	//	return "", err
	//}

	return u.Name, nil
}

func (mt *Service) GetUser(email string) (*User, error) {
	var records string
	user := &User{}

	res := mt.db.QueryRow(getUserQuery, email)
	if err := res.Scan(&user.Name, &user.Surename, &user.Email, &user.Id, &user.Rights, &records); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(records), &user.Records); err != nil {
		return nil, err
	}

	return user, nil
}

func (mt *Service) UpdateUser(field, data string, id int) error {
	stmt, err := mt.db.Prepare("UPDATE users SET " + field + "=? WHERE id=?;")
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

func (mt *Service) DeleteUser(id int) error {
	stmt, err := mt.db.Prepare(deleteUserQuery)
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	if _, err = stmt.Exec(id); err != nil {
		return err
	}
	return nil
}

func (mt *Service) CheckUserExists(email string) bool {
	var result string
	res := mt.db.QueryRow(checkUserExistsQuery, email)
	if err := res.Scan(&result); err != nil {
		logrus.Error(err)
		return false
	}
	return result == ""
}

func (mt *Service) SetActiveField(email string) error {
	stmt, err := mt.db.Prepare("UPDATE users SET activated=1 WHERE email=?")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(email); err != nil {
		return err
	}
	return nil
}

const updatePasswordQuery = `UPDATE users SET password=? WHERE email=?`

func (mt *Service) UpdatePassword(password, hash string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	var u User
	//todo move outside
	h, err := u.ReadHash(hash)
	if err != nil {
		return err
	}
	if h == nil {
		return errors.New("посилання не існує")
	}
	res, err := mt.db.Prepare(updatePasswordQuery)
	if err != nil {
		return err
	}
	if _, err = res.Exec(hashedPassword, h.Email); err != nil {
		return err
	}

	//todo move outside
	if err = u.DeleteHash(h.Hash); err != nil {
		return err
	}
	return nil
}
