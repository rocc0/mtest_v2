package mtest

import (
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type HashData struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

type hasher interface {
	WriteHash(hash, email string) (HashData, error)
	ReadHash(hash string) (HashData, error)
	DeleteHash(hash string) (err error)
}

type HashHandler struct {
	*mongo.Client
}

func (u *HashHandler) WriteHash(hash, email string) (HashData, error) {
	dialInfo, err := mgo.ParseURL("mongodb://hasher:password@localhost:27017")
	if err != nil {
		return HashData{}, err
	}
	dialInfo.Direct = true
	dialInfo.FailFast = true
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return HashData{}, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("hashes").C(hash)
	h := HashData{email, hash}
	return h, c.Insert(&h)
}

func (u *HashHandler) ReadHash(hash string) (HashData, error) {
	dialInfo, err := mgo.ParseURL("mongodb://hasher:password@localhost:27017")
	if err != nil {
		return HashData{}, err
	}
	dialInfo.Direct = true
	dialInfo.FailFast = true
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return HashData{}, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("hashes").C(hash)
	var h HashData
	if err = c.Find(nil).One(&h); err != nil {
		return h, err
	}

	return h, nil
}

func (u *HashHandler) DeleteHash(hash string) (err error) {
	session, err := mgo.Dial("mongodb://hasher:password@localhost:27017")
	if err != nil {
		return err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	return session.DB("hashes").C(hash).DropCollection()
}

func (u *HashHandler) generateHash(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
