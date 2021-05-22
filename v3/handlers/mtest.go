package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"

	dataprocpkg "mtest.com.ua/v3/db/dataprocessor"
	internal2 "mtest.com.ua/v3/handlers/internal"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

type deleteRequest struct {
	Id string `json:"id"`
}

type mtestDataProcessor interface {
	UpdateMTEST(m map[string]interface{}, email string) error
	DeleteMTEST(mid, email string) error
	GetAdministrativeActions() (*[]dataprocpkg.AdmAction, error)
	GetMTEST(id string) (*dataprocpkg.MTEST, error)
	CreateMTEST(m dataprocpkg.NewMTEST, email string) (dataprocpkg.UserMtest, error)
}

type executorDataProcessor interface {
	CreateExecutor(email string, ex dataprocpkg.Executor) (string, error)
	DeleteExecutor(devEmail string, del dataprocpkg.DeleteExecutor) error
}

type regionDataProcessor interface {
	GetRegions() (*[]dataprocpkg.Region, error)
	EditRegionName(id int, name string) error
	GetGovernments() (*[]dataprocpkg.Government, error)
	EditGovernmentName(id int, name string) error
}

type indexUpdater interface {
	UpdateIndex(id int64) error
}

type Handlers struct {
	mtestDataProcessor
	executorDataProcessor
	regionDataProcessor
	userDataProcessor
	hasher
	indexUpdater
}

func (hd *Handlers) RenderIndexPage(c *gin.Context) {
	internal2.Render(c, gin.H{"title": "Калькулятор"}, "index.html")
}

func (hd *Handlers) RenderSearchPage(c *gin.Context) {
	internal2.Render(c, gin.H{"title": "Пошук АРВ"}, "index.html")
}

func (hd *Handlers) RenderUserPage(c *gin.Context) {
	internal2.Render(c, gin.H{"title": "Кабінет користувача"}, "index.html")
}

func (hd *Handlers) RenderMTESTPage(c *gin.Context) {
	mtest, err := hd.GetMTEST(c.Param("mtest_id"))
	if err == nil {
		internal2.Render(c, gin.H{"title": "Редагування | " + mtest.Name}, "index.html")
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}

}

func (hd *Handlers) GetMTESTHandler(c *gin.Context) {
	mtest, err := hd.GetMTEST(c.Param("mtest_id"))
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"mtest": mtest})
	} else {
		log.Print(err)
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func (hd *Handlers) CreateMTESTHandler(c *gin.Context) {
	em, ok := jwt.ExtractClaims(c)["id"]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	email, ok := em.(string)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var m dataprocpkg.NewMTEST
	if err := json.Unmarshal(x, &m); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if data, err := hd.CreateMTEST(m, email); err == nil {
		go func() {
			if err := hd.UpdateIndex(data.RecID); err != nil {
				logrus.Error(err)
			}
		}()
		c.JSON(http.StatusOK, gin.H{"title": "Item added", "records": data})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

}

func (hd *Handlers) UpdateMTESTHandler(c *gin.Context) {
	em, ok := jwt.ExtractClaims(c)["id"]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	email, ok := em.(string)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var form map[string]interface{}
	if err := json.Unmarshal(x, &form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.UpdateMTEST(form, email); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Mtest updated", "data": form})
	} else {
		log.Print(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) DeleteMTESTHandler(c *gin.Context) {
	em, ok := jwt.ExtractClaims(c)["id"]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	email, ok := em.(string)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var id deleteRequest
	if err := json.Unmarshal(x, &id); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.DeleteMTEST(id.Id, email); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Item removed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) GetGovernmentsHandlers(c *gin.Context) {
	res, err := hd.GetGovernments()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"govs": res})
}

func (hd *Handlers) GetRegionsHandler(c *gin.Context) {
	res, err := hd.GetRegions()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"regions": res})
}

func (hd *Handlers) CreateMTESTExecutorHandler(c *gin.Context) {
	em, ok := jwt.ExtractClaims(c)["id"]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	email, ok := em.(string)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var executor dataprocpkg.Executor
	if err := json.Unmarshal(x, &executor); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if mid, err := hd.CreateExecutor(email, executor); err == nil {
		c.JSON(http.StatusOK, gin.H{"mid": mid})
	} else {
		log.Print(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) DeleteExecutorHandler(c *gin.Context) {
	em, ok := jwt.ExtractClaims(c)["id"]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	email, ok := em.(string)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var delRequest dataprocpkg.DeleteExecutor
	if err := json.Unmarshal(x, &delRequest); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.DeleteExecutor(email, delRequest); err == nil {
		c.JSON(http.StatusOK, gin.H{"response": "ok"})
	} else {
		log.Print(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) GetAdministrativeActionsHandler(c *gin.Context) {
	res, err := hd.GetAdministrativeActions()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"actions": res})
}
