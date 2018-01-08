package main

import (
	"io/ioutil"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func showEditGovernments(c *gin.Context) {
	render(c, gin.H{
		"title":   "Пошук відстежень",
	}, "index.html")
}



func showAdminPage(c *gin.Context) {
	render(c, gin.H{
		"title":   "Пошук відстежень",
	}, "index.html")
}


func postEditGovernments(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var editgov editGov
	err := json.Unmarshal([]byte(x), &editgov)
	check(err)
	if err := editGovName(editgov.Id, editgov.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"title": "Gov name changed",
		})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}