package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func (hd *Handlers) GetBusinessesHandler(c *gin.Context) {
	res, err := hd.GetBusinesses()
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"businesses": res})
}

func (hd *Handlers) EditBusinessHandler(c *gin.Context) {
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

	if err := hd.EditBusinessName(edit.Id, edit.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Business name changed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) AddBusinessHandler(c *gin.Context) {
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

	if err := hd.AddBusiness(edit.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Business added"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) DeleteBusinessHandler(c *gin.Context) {
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

	if err := hd.RemoveBusiness(edit.Id); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Business removed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
