package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"code.sajari.com/docconv"
	"github.com/gin-gonic/gin"
	datapkg "mtest.com.ua/db/dataprocessor"
)

type regActUpdater interface {
	InsertRegAct(mtestID string, docText string, docName string) (string, error)
	DeleteRegAct(mtestID string, docID string) error
	GetRegAct(mtestID string, docID string) (datapkg.RegAct, error)
	ListRegActs(mtestID string) ([]datapkg.RegAct, error)
}

type regAct struct {
	MtestID string `json:"mtest_id"`
	DocID   string `json:"doc_id,omitempty"`
	Name    string `json:"name,omitempty"`
}

func (hd *Handlers) ActUploadHandler(c *gin.Context) {
	var act regAct
	act.MtestID = c.PostForm("mtestID")
	file, err := c.FormFile("file")
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "No file is received"})
		return
	}
	f, err := file.Open()
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid file"})
		return
	}
	res, err := docconv.Convert(f, docconv.MimeTypeByExtension(file.Filename), true)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid file"})
		return
	}

	docID, err := hd.InsertRegAct(act.MtestID, res.Body, file.Filename)
	if err != nil {
		logrus.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := c.SaveUploadedFile(file, "/var/www/reg_acts/"+docID); err != nil {
		logrus.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := hd.UpdateIndexWithFile(act.MtestID, res.Body); err != nil {
		logrus.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	act.DocID = docID
	act.Name = file.Filename
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
		resp := gin.H{}
		if len(list) > 0 {
			resp["reg_acts"] = list[0]
		} else {
			resp["reg_acts"] = []datapkg.RegAct{}
		}
		c.JSON(http.StatusOK, resp)
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
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	actData, err := hd.GetRegAct(act.MtestID, act.DocID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", actData.Name))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Add("Content-Description", "File Transfer")
	c.Writer.Header().Add("Content-Transfer-Encoding", "binary")

	c.FileAttachment("/var/www/reg_acts/"+actData.DocID, actData.Name)
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
	defer func() {
		if err := hd.UpdateIndexWithFile(act.MtestID, ""); err != nil {
			logrus.Error(err)
		}
	}()
	if err := hd.DeleteRegAct(act.MtestID, act.DocID); err == nil {
		c.JSON(200, gin.H{"title": "Документ видалено"})
		if err := os.Remove("/var/www/reg_acts/" + act.DocID); err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	} else {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
