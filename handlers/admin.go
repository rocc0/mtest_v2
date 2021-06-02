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
	UpdateUser(field string, data interface{}, id int) error
	GetUsers(c context.Context) ([]datapkg.User, error)
}

type editRequest struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (hd *Handlers) ShowEditGovernments(c *gin.Context) {
	internal.Render(c, gin.H{"title": "Пошук відстежень"}, "index.html")
}

func (hd *Handlers) RenderAdminPage(c *gin.Context) {
	internal.Render(c, gin.H{"title": "Пошук відстежень"}, "index.html")
}

func (hd *Handlers) EditGovernmentNameHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var edit editRequest
	if err := json.Unmarshal(x, &edit); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.EditGovernmentName(edit.Id, edit.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Gov name changed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) AddGovernmentHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var edit editRequest
	if err := json.Unmarshal(x, &edit); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.AddGovernment(edit.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Adm action add"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) RemoveGovernmentHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var edit editRequest
	if err := json.Unmarshal(x, &edit); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.RemoveGovernment(edit.Id); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Adm action add"})
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
	var edit editRequest
	if err := json.Unmarshal(x, &edit); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := hd.EditRegionName(edit.Id, edit.Name); err == nil {
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

func (hd *Handlers) AddAdministrativeActionHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var edit editRequest
	if err := json.Unmarshal(x, &edit); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.AddAdministrativeAction(edit.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Adm action add"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) EditAdministrativeActionsHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var edit editRequest
	if err := json.Unmarshal(x, &edit); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.EditAdministrativeActionName(edit.Id, edit.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Adm action name changed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) DeleteAdministrativeActionsHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var edit editRequest
	if err := json.Unmarshal(x, &edit); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.DeleteAdministrativeAction(edit.Id); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Adm action name deleted"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) ValidateAdminRights(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"title": "OK"})
}
