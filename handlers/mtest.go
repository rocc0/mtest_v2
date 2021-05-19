package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"mtest.com.ua/db/dataprocessor"

	"mtest.com.ua/handlers/internal"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

type deleteRequest struct {
	Id string `json:"id"`
}

type mtestDataProcessor interface {
	UpdateMTEST(m map[string]interface{}, email string) error
	DeleteMTEST(mid, email string) error
	GetAdministrativeActions() (*[]dataprocessor.AdmAction, error)
	GetMTEST(id string) (*dataprocessor.MTEST, error)
	CreateMTEST(m dataprocessor.NewMTEST, email string) (*map[string]interface{}, error)
}

type executorDataProcessor interface {
	CreateExecutor(email string, ex dataprocessor.Executor) (string, error)
	DeleteExecutor(devEmail string, del dataprocessor.DeleteExecutor) error
}

type regionDataProcessor interface {
	GetRegions() (*[]dataprocessor.Region, error)
	EditRegionName(id int, name string) error
	GetGovernments() (*[]dataprocessor.Government, error)
	EditGovernmentName(id int, name string) error
}

type userDataProcessor interface {
	CheckUserActivation(email string) bool
	CheckUserExists(email string) bool
	CreateUser() (string, error)
	DeleteUser(id int) error
	GetUser(email string) (*dataprocessor.User, error)
	InitUsersTable() error
	PasswordCheck(email, password string) bool
	SetActiveField(email string) error
	UpdatePassword(password, hash string) error
	UpdateUser(field, data string, id int) error
}

type Handlers struct {
	mtestDataProcessor
	executorDataProcessor
	regionDataProcessor
	userDataProcessor
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
	id := c.Param("mtest_id")
	mtest, err := hd.GetMTEST(id)
	if err == nil {
		internal.Render(c, gin.H{"title": "Редагування | " + mtest.Name}, "index.html")
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}

}

func (hd *Handlers) GetMTESTHandler(c *gin.Context) {
	id := c.Param("mtest_id")
	mtest, err := hd.GetMTEST(id)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"mtest": mtest})
	} else {
		log.Print(err)
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func (hd *Handlers) CreateMTESTHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	email, ok := claims["id"].(string)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var m dataprocessor.NewMTEST
	if err := json.Unmarshal(x, &m); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if data, err := hd.CreateMTEST(m, email); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Item added", "records": data})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Print(err)
	}
}

func (hd *Handlers) UpdateMTESTHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	email, ok := claims["id"].(string)
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
	claims := jwt.ExtractClaims(c)
	email, ok := claims["id"].(string)
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
	claims := jwt.ExtractClaims(c)
	email, ok := claims["id"].(string)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var executor dataprocessor.Executor
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
	claims := jwt.ExtractClaims(c)
	email, ok := claims["id"].(string)
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var delRequest dataprocessor.DeleteExecutor
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
