package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"mtest.com.ua/db/dataprocessor"
	"mtest.com.ua/handlers/internal"

	"github.com/gin-gonic/gin"
)

type userDataProcessor interface {
	CheckUserActivation(email string) bool
	CheckUserExists(email string) bool
	CreateUser() (string, error)
	DeleteUser(id int) error
	GetUser(email string) (*dataprocessor.User, error)
	InitUsersTable() error
	PasswordCheck(email, password string) bool
	SetActiveField(email string) error
	UpdatePassword(password, hash string) error
	UpdateUser(field, data string, id int) error
}

type editGov struct {
	Id   int
	Name string
}

type hasher interface {
	WriteHash(hash, email string) (HashData, error)
	ReadHash(hash string) (HashData, error)
	DeleteHash(hash string) (err error)
}

func (hd *Handlers) showEditGovernments(c *gin.Context) {
	internal.Render(c, gin.H{"title": "Пошук відстежень"}, "index.html")
}

func (hd *Handlers) showAdminPage(c *gin.Context) {
	internal.Render(c, gin.H{"title": "Пошук відстежень"}, "index.html")
}

func (hd *Handlers) postEditGovernments(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var editgov editGov
	if err := json.Unmarshal(x, &editgov); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.EditGovernmentName(editgov.Id, editgov.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Gov name changed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) postEditRegions(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var editgov editGov
	if err := json.Unmarshal(x, &editgov); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := hd.EditRegionName(editgov.Id, editgov.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Gov name changed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
