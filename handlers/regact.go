package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"code.sajari.com/docconv"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	datapkg "mtest.com.ua/db/dataprocessor"
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
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var act regAct
	if err := json.Unmarshal(x, &act); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
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

	docID, err := hd.InsertRegAct(act.MtestID, res.Body, file.Filename, filepath.Ext(file.Filename))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := hd.UpdateIndexWithFile(act.MtestID, res.Body); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	act.DocID = docID
	c.JSON(http.StatusOK, gin.H{"act": act})
}

func (hd *Handlers) ActsListHandler(c *gin.Context) {
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

	list, err := hd.ListRegActs(act.MtestID)
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

	pdf := gofpdf.New("P", "mm", "A4", "")
	fontSz := float64(16)
	lineSz := pdf.PointToUnitConvert(fontSz)
	pdf.SetFont("Arial", "B", fontSz)
	pdf.Write(lineSz, actData.Text)
	defer pdf.Close()
	filename := actData.Name + "." + actData.Type
	if err := pdf.OutputFileAndClose("/tmp/" + filename); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", actData.Name+"."+actData.Type))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Status(http.StatusOK)
	c.File("/tmp/" + filename)
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
