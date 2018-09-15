package main

import (
	"math/rand"
	"time"
	"gopkg.in/mgo.v2"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Hash struct {
	Email string `json:"email"`
	Hash string `json:"hash"`
}

func (u *User) writeHash() (*Hash, error) {
	hash := u.generateHash(20)

	dialInfo, err := mgo.ParseURL("mongodb://hasher:password@localhost:27017")
	dialInfo.Direct = true
	dialInfo.FailFast = true
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("hashes").C(hash)
	err = c.Insert(&Hash{u.Email, hash})
	if err != nil {
		return nil, err
	}
	h := Hash{u.Email, hash}

	return &h, nil
}

func (u *User) readHash(hash string) (*Hash, error) {
	var h Hash
	dialInfo, err := mgo.ParseURL("mongodb://hasher:password@localhost:27017")
	dialInfo.Direct = true
	dialInfo.FailFast = true
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("hashes").C(hash)

	err = c.Find(nil).One(&h)

	if err != nil {
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
	err = c.DropCollection()
	
	if err != nil {
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
	_, err = stmt.Exec(email)
	if err != nil {
		return err
	}
	return nil
}

func passwordResetter(password, hash string) error {
	var u User
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	h, err := u.readHash(hash)
	if err != nil {
		return err
	}
	if h == nil {
		return errors.New("посилання не існує")
	}
	res, _ := db.Prepare("UPDATE users SET password=? WHERE email=?")
	_, err = res.Exec(hashedPassword, h.Email)
	if err != nil {
		return err
	}
	err = u.deleteHash(h.Hash)
	if err != nil {
		return err
	}
	return nil
}