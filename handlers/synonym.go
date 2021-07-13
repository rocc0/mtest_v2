package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type synonymRequest struct {
	MtestID string `json:"mtest_id"`
	Synonym string `json:"synonym"`
}

func (hd *Handlers) GetAllSynonyms(c *gin.Context) {
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

	res, err := hd.GetSynonymsByID(syn.MtestID)
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

	if synID, err := hd.AddSynonym(syn.MtestID, syn.Synonym); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Synonym added", "synonym_id": synID})
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

	if err := hd.RemoveSynonym(syn.MtestID, syn.Synonym); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Synonym removed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
