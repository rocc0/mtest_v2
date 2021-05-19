package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"mtest.com.ua/handlers/internal"

	"github.com/gin-gonic/gin"
)

type editGov struct {
	Id   int
	Name string
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
