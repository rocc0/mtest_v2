package hash

import (
	"math/rand"
	"time"

	"gopkg.in/mgo.v2"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type HashData struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

//todo replace with official client
type HashHandler struct {
	*mgo.DialInfo
}

func NewHashHandler(address string) (*HashHandler, error) {
	if address == "" {
		address = "mongodb://hasher:password@localhost:27017"
	}
	dialInfo, err := mgo.ParseURL(address)
	if err != nil {
		return nil, err
	}
	dialInfo.Direct = true
	dialInfo.FailFast = true
	return &HashHandler{DialInfo: dialInfo}, nil
}

func (u *HashHandler) WriteHash(email string) (string, error) {

	h := u.generateHash(20)

	session, err := mgo.DialWithInfo(u.DialInfo)
	if err != nil {
		return "", err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("hashes").C(h)

	return h, c.Insert(&HashData{Email: email, Hash: h})
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
	session, err := mgo.DialWithInfo(u.DialInfo)
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
