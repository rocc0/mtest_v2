package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type editGov struct {
	Id   int
	Name string
}

func showEditGovernments(c *gin.Context) {
	render(c, gin.H{"title": "Пошук відстежень"}, "index.html")
}

func showAdminPage(c *gin.Context) {
	render(c, gin.H{"title": "Пошук відстежень"}, "index.html")
}

func postEditGovernments(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var editgov editGov
	if err := json.Unmarshal(x, &editgov); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := editGovName(editgov.Id, editgov.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Gov name changed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func postEditRegions(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var editgov editGov
	if err := json.Unmarshal(x, &editgov); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := editRegName(editgov.Id, editgov.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Gov name changed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
