package main

import (
	"errors"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Hash struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

func (u *User) writeHash() (*Hash, error) {
	hash := u.generateHash(20)

	dialInfo, err := mgo.ParseURL("mongodb://hasher:password@localhost:27017")
	if err != nil {
		return nil, err
	}
	dialInfo.Direct = true
	dialInfo.FailFast = true
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("hashes").C(hash)
	if err = c.Insert(&Hash{u.Email, hash}); err != nil {
		return nil, err
	}
	return &Hash{u.Email, hash}, nil
}

func (u *User) readHash(hash string) (*Hash, error) {
	dialInfo, err := mgo.ParseURL("mongodb://hasher:password@localhost:27017")
	if err != nil {
		return nil, err
	}
	dialInfo.Direct = true
	dialInfo.FailFast = true
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("hashes").C(hash)
	var h Hash
	if err = c.Find(nil).One(&h); err != nil {
		return nil, err
	}

	return &h, nil
}

func (u *User) deleteHash(h string) (err error) {
	session, err := mgo.Dial("mongodb://hasher:password@localhost:27017")
	if err != nil {
		return err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("hashes").C(h)
	if err = c.DropCollection(); err != nil {
		return err
	}

	return nil
}

func (u *User) generateHash(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func setActiveField(email string) error {
	stmt, err := db.Prepare("UPDATE users SET activated=1 WHERE email=?")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(email); err != nil {
		return err
	}
	return nil
}

func passwordResetter(password, hash string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	var u User
	h, err := u.readHash(hash)
	if err != nil {
		return err
	}
	if h == nil {
		return errors.New("посилання не існує")
	}
	res, _ := db.Prepare("UPDATE users SET password=? WHERE email=?")
	if _, err = res.Exec(hashedPassword, h.Email); err != nil {
		return err
	}

	if err = u.deleteHash(h.Hash); err != nil {
		return err
	}
	return nil
}
