package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	datapkg "mtest.com.ua/db/dataprocessor"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type synonymRequest struct {
	MtestID   string `json:"mtest_id"`
	Synonym   string `json:"synonym"`
	SynonymID string `json:"synonym_id"`
}
type globalSynonymRequest struct {
	Word    string `json:"word"`
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
		go func() {
			if err := hd.UpdateIndex(syn.MtestID); err != nil {
				logrus.Error(err)
			}
		}()
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

	if err := hd.RemoveSynonym(syn.MtestID, syn.SynonymID); err == nil {
		go func() {
			if err := hd.UpdateIndex(syn.MtestID); err != nil {
				logrus.Error(err)
			}
		}()
		c.JSON(http.StatusOK, gin.H{"title": "Synonym removed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) GetGlobalSynonyms(c *gin.Context) {
	res, err := hd.LoadHandler()
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"synonyms": res})
}

func (hd *Handlers) LoadHandler() ([]datapkg.GlobalSynonym, error) {
	return hd.LoadGlobals()
}

func (hd *Handlers) AddGlobalSynonymHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var syn globalSynonymRequest
	if err := json.Unmarshal(x, &syn); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.AddGlobalSynonym(syn.Word, syn.Synonym); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Synonym added"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
func (hd *Handlers) RemoveGlobalSynonymHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var syn globalSynonymRequest
	if err := json.Unmarshal(x, &syn); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.RemoveGlobalSynonym(syn.Word, syn.Synonym); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Synonym removed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
