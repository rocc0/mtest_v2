package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"

	datapkg "mtest.com.ua/db/dataprocessor"
	"mtest.com.ua/handlers/internal"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type deleteRequest struct {
	Id string `json:"id"`
}

type mtestDataProcessor interface {
	UpdateMTEST(m map[string]interface{}, email string) error
	DeleteMTEST(mid, email string) error
	GetMTEST(id string) (*datapkg.MTEST, error)
	CreateMTEST(m datapkg.NewMTEST, email string) (datapkg.MTestData, error)
}

type admActionsProcessor interface {
	GetAdministrativeActions() (*[]datapkg.AdmAction, error)
	EditAdministrativeActionName(id int, name string) error
	AddAdministrativeAction(name string) error
	DeleteAdministrativeAction(id int) error
}

type executorDataProcessor interface {
	CreateExecutor(email string, ex datapkg.Executor) (string, error)
	DeleteExecutor(devEmail string, del datapkg.DeleteExecutor) error
}

type regionDataProcessor interface {
	GetRegions() (*[]datapkg.Region, error)
	EditRegionName(id int, name string) error
}

type governmentDataProcessor interface {
	AddGovernment(name string) error
	GetGovernments() ([]datapkg.Government, error)
	EditGovernmentName(id int, name string) error
	RemoveGovernment(id int) error
}

type BusinessDataProcessor interface {
	AddBusiness(name string) error
	GetBusinesses() ([]datapkg.Government, error)
	EditBusinessName(id int, name string) error
	RemoveBusiness(id int) error
}

type indexUpdater interface {
	UpdateIndex(id string) error
	UpdateIndexWithFile(id string, text string) error
}

type Handlers struct {
	mtestDataProcessor
	executorDataProcessor
	regionDataProcessor
	userDataProcessor
	hasher
	indexUpdater
	admActionsProcessor
	governmentDataProcessor
	BusinessDataProcessor
	regActUpdater
}

func (hd *Handlers) RenderIndexPage(c *gin.Context) {
	internal.Render(c, gin.H{"title": "Калькулятор"}, "index.html")
}

func (hd *Handlers) RenderSearchPage(c *gin.Context) {
	internal.Render(c, gin.H{"title": "Пошук АРВ"}, "index.html")
}

func (hd *Handlers) RenderUserPage(c *gin.Context) {
	internal.Render(c, gin.H{"title": "Кабінет користувача"}, "index.html")
}

func (hd *Handlers) RenderMTESTPage(c *gin.Context) {
	mtest, err := hd.GetMTEST(c.Param("mtest_id"))
	if err == nil {
		internal.Render(c, gin.H{"title": "Редагування | " + mtest.Name}, "index.html")
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}

}

func (hd *Handlers) GetMTESTHandler(c *gin.Context) {
	mtest, err := hd.GetMTEST(c.Param("mtest_id"))
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"mtest": mtest})
	} else {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func (hd *Handlers) CreateMTESTHandler(c *gin.Context) {
	em, ok := jwt.ExtractClaims(c)["id"]
	if !ok {
		logrus.Error("no email in request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	email, ok := em.(string)
	if !ok {
		logrus.Error("no enail in request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var m datapkg.NewMTEST
	if err := json.Unmarshal(x, &m); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if data, err := hd.CreateMTEST(m, email); err == nil {
		logrus.Error(err)
		go func() {
			if err := hd.UpdateIndex(data.Id); err != nil {
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
		logrus.Error("no email in request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	email, ok := em.(string)
	if !ok {
		logrus.Error("no email in request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var form map[string]interface{}
	if err := json.Unmarshal(x, &form); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.UpdateMTEST(form, email); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Mtest updated", "data": form})
	} else {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) DeleteMTESTHandler(c *gin.Context) {
	em, ok := jwt.ExtractClaims(c)["id"]
	if !ok {
		logrus.Error("no ID in request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	email, ok := em.(string)
	if !ok {
		logrus.Error("no ID in request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var id deleteRequest
	if err := json.Unmarshal(x, &id); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.DeleteMTEST(id.Id, email); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Item removed"})
	} else {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) GetGovernmentsHandlers(c *gin.Context) {
	res, err := hd.GetGovernments()
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"govs": res})
}

func (hd *Handlers) GetRegionsHandler(c *gin.Context) {
	res, err := hd.GetRegions()
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"regions": res})
}

func (hd *Handlers) CreateMTESTExecutorHandler(c *gin.Context) {
	em, ok := jwt.ExtractClaims(c)["id"]
	if !ok {
		logrus.Error("no email in request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	email, ok := em.(string)
	if !ok {
		logrus.Error("no email in request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var executor datapkg.Executor
	if err := json.Unmarshal(x, &executor); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if mid, err := hd.CreateExecutor(email, executor); err == nil {
		c.JSON(http.StatusOK, gin.H{"mid": mid})
	} else {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) DeleteExecutorHandler(c *gin.Context) {
	em, ok := jwt.ExtractClaims(c)["id"]
	if !ok {
		logrus.Error("no email in request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	email, ok := em.(string)
	if !ok {
		logrus.Error("no email in request")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var delRequest datapkg.DeleteExecutor
	if err := json.Unmarshal(x, &delRequest); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := hd.DeleteExecutor(email, delRequest); err == nil {
		c.JSON(http.StatusOK, gin.H{"response": "ok"})
	} else {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (hd *Handlers) GetAdministrativeActionsHandler(c *gin.Context) {
	res, err := hd.GetAdministrativeActions()
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"actions": res})
}
