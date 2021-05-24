package hash

import (
	"database/sql"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type HashData struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

const createHashesTable = `CREATE TABLE IF NOT EXISTS
    hashes (id SERIAL NOT NULL PRIMARY KEY, email VARCHAR(100) NOT NULL, 
    hash VARCHAR(20) NOT NULL);`
const createHashQuery = `INSERT INTO hashes (email, hash) VALUES (?,?);`
const getHashQuery = `SELECT email FROM hashes WHERE hash=?;`
const deleteHashQuery = `DELETE FROM hashes WHERE hash=?;`

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//todo replace with official client
type HashHandler struct {
	db *sql.DB
}

func NewHashHandler(db *sql.DB) (*HashHandler, error) {
	return &HashHandler{db: db}, nil
}

func (u HashHandler) Init() error {
	stmt, err := u.db.Prepare(createHashesTable)
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

func (u *HashHandler) WriteHash(email string) (string, error) {
	h := u.generateHash(20)

	stmt, err := u.db.Prepare(createHashQuery)
	if err != nil {
		return h, err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Error(err)
		}
	}()
	_, err = stmt.Exec(email, h)
	if err != nil {
		return h, err
	}

	return h, nil
}

func (u *HashHandler) ReadHash(hash string) (HashData, error) {
	hashStmt := u.db.QueryRow(getHashQuery, hash)
	var email string
	if err := hashStmt.Scan(&email); err != nil {
		return HashData{}, err
	}
	return HashData{Email: email, Hash: hash}, nil
}

func (u *HashHandler) DeleteHash(hash string) (err error) {
	stmt, err := u.db.Prepare(deleteHashQuery)
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	if _, err = stmt.Exec(hash); err != nil {
		return err
	}
	return nil
}

func (u *HashHandler) generateHash(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
