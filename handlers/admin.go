package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"

	datapkg "mtest.com.ua/db/dataprocessor"
	"mtest.com.ua/handlers/internal"

	"github.com/gin-gonic/gin"
)

type userDataProcessor interface {
	CheckUserActivation(email string) bool
	CheckUserExists(email string) bool
	CreateUser() (string, error)
	DeleteUser(id int) error
	GetUser(email string) (*datapkg.User, error)
	InitUsersTable() error
	PasswordCheck(email, password string) bool
	SetActiveField(email string) error
	UpdatePassword(password, email, hash string) error
	UpdateUser(field, data string, id int) error
	GetUsers(c context.Context) ([]datapkg.User, error)
}

type editGovRequest struct {
	Id   int
	Name string
}

func (hd *Handlers) ShowEditGovernments(c *gin.Context) {
	internal.Render(c, gin.H{"title": "Пошук відстежень"}, "index.html")
}

func (hd *Handlers) RenderAdminPage(c *gin.Context) {
	internal.Render(c, gin.H{"title": "Пошук відстежень"}, "index.html")
}

func (hd *Handlers) PostEditGovernments(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var editGov editGovRequest
	if err := json.Unmarshal(x, &editGov); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.EditGovernmentName(editGov.Id, editGov.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Gov name changed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) PostEditRegions(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var editGov editGovRequest
	if err := json.Unmarshal(x, &editGov); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := hd.EditRegionName(editGov.Id, editGov.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Gov name changed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) GetUsersHandler(c *gin.Context) {
	users, err := hd.GetUsers(c)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"users": users})
	} else {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
