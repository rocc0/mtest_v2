package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/sirupsen/logrus"

	datapkg "mtest.com.ua/db/dataprocessor"

	"code.sajari.com/docconv"
	"github.com/gin-gonic/gin"
)

type regActUpdater interface {
	InsertRegAct(mtestID string, docText string, docName string, docType string) (string, error)
	DeleteRegAct(mtestID string, docID string) error
	GetRegAct(mtestID string, docID string) (datapkg.RegAct, error)
	ListRegActs(mtestID string) ([]datapkg.RegAct, error)
}

type regAct struct {
	MtestID string `json:"mtest_id"`
	DocID   string `json:"doc_id"`
}

func (hd *Handlers) ActUploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "No file is received"})
		return
	}
	f, err := file.Open()
	if err != nil {
		log.Fatal(err)
	}
	res, err := docconv.Convert(f, docconv.MimeTypeByExtension(file.Filename), true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid file"})
		return
	}
	hd.InsertRegAct("mid", res.Body, file.Filename, filepath.Ext(file.Filename))

	if err := hd.UpdateIndexWithFile(0, res.Body); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	/*
		todo
		reg act name, file name
		0. convert to text
		1. upload file to postgres
		2. upload text to elastic*/
	c.JSON(http.StatusOK, gin.H{
		"message": "Your file has been successfully uploaded.",
	})
}

func (hd *Handlers) ActsListHandler(c *gin.Context) {
	list, err := hd.ListRegActs("")
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"reg_acts": list})
	} else {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (hd *Handlers) ActGetHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(400)
		return
	}
	var act regAct
	if err := json.Unmarshal(x, &act); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(400)
		return
	}
	actData, err := hd.GetRegAct(act.MtestID, act.DocID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"act": actData})
}

func (hd *Handlers) ActDeleteHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(400)
		return
	}
	var act regAct
	if err := json.Unmarshal(x, &act); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(400)
		return
	}
	if err := hd.DeleteRegAct(act.MtestID, act.DocID); err == nil {
		c.JSON(200, gin.H{"title": "Документ видалено"})
	} else {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
