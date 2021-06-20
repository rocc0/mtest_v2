package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type synonymRequest struct {
	Word    string `json:"word"`
	Synonym string `json:"synonym"`
}

func (hd *Handlers) GetAllSynonyms(c *gin.Context) {
	res, err := hd.LoadHandler()
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"synonyms": res})
}

func (hd *Handlers) AddSynonymHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var syn synonymRequest
	if err := json.Unmarshal(x, &syn); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.AddSynonym(syn.Word, syn.Synonym); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Synonym added"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
func (hd *Handlers) RemoveSynonymHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var syn synonymRequest
	if err := json.Unmarshal(x, &syn); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.RemoveSynonym(syn.Word, syn.Synonym); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Synonym removed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
