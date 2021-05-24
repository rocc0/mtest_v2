package v2

import (
	"github.com/gin-gonic/gin"

	"net/http"

	"encoding/json"
	"io/ioutil"
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
)

type (
	newMtest struct {
		Name       string `json:"name"`
		Region     int    `json:"region"`
		Government int    `json:"government"`
		CalcType   int    `json:"calc_type"`
	}
	deleteRequest struct {
		Id string `json:"id"`
	}
	delExecutorReq struct {
		ExEmail    string `json:"ex_email"`
		ExMtestId  string `json:"ex_mtest_id"`
		DevMtestId string `json:"dev_mtest_id"`
	}
	newExecutor struct {
		Title      string `json:"title"`
		Email      string `json:"email"`
		Region     int    `json:"region"`
		Government int    `json:"government"`
		DevMid     string `json:"dev_mid"`
	}
)

// page render
func renderIndexPage(c *gin.Context) {
	render(c, gin.H{"title": "Калькулятор"}, "index.html")
}

func renderSearchPage(c *gin.Context) {
	render(c, gin.H{"title": "Пошук АРВ"}, "index.html")
}

func renderUserPage(c *gin.Context) {
	render(c, gin.H{"title": "Кабінет користувача"}, "index.html")
}

func renderMTESTPage(c *gin.Context) {
	id := c.Param("mtest_id")
	mtest, err := getMTEST(id)
	if err == nil {
		render(c, gin.H{"title": "Редагування | " + mtest.Name}, "index.html")
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}

}

func getMTESTHandler(c *gin.Context) {
	id := c.Param("mtest_id")
	mtest, err := getMTEST(id)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"mtest": mtest})
	} else {
		log.Print(err)
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func createMTESTHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	email, _ := claims["id"].(string)

	x, _ := ioutil.ReadAll(c.Request.Body)
	var m newMtest
	if err := json.Unmarshal(x, &m); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if data, err := createMTEST(m, email); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Item added", "records": data})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Print(err)
	}
}

func updateMTESTHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	email, _ := claims["id"].(string)

	x, _ := ioutil.ReadAll(c.Request.Body)
	var form map[string]interface{}
	if err := json.Unmarshal(x, &form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := updateMTEST(form, email); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Mtest updated", "data": form})
	} else {
		log.Print(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func deleteMTESTHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	email, _ := claims["id"].(string)

	x, _ := ioutil.ReadAll(c.Request.Body)
	var id deleteRequest
	if err := json.Unmarshal(x, &id); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := deleteMTEST(id.Id, email); err == nil {
		c.JSON(http.StatusOK, gin.H{"title": "Item removed"})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

//api govs
func getGovernmentsHandlers(c *gin.Context) {
	res, err := getGovernments()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"govs": res})
}

func getRegionsHandler(c *gin.Context) {
	res, err := getRegions()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"regions": res})
}

//executors and group calculations
//add executor
func createMTESTExecutorHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	email, _ := claims["id"].(string)

	x, _ := ioutil.ReadAll(c.Request.Body)
	var executor newExecutor
	if err := json.Unmarshal(x, &executor); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if mid, err := createMTESTExecutor(email, executor); err == nil {
		c.JSON(http.StatusOK, gin.H{"mid": mid})
	} else {
		log.Print(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func deleteExecutorHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	email, _ := claims["id"].(string)

	x, _ := ioutil.ReadAll(c.Request.Body)
	var delRequest delExecutorReq
	if err := json.Unmarshal(x, &delRequest); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := deleteMTESTExecutor(email, delRequest); err == nil {
		c.JSON(http.StatusOK, gin.H{"response": "ok"})
	} else {
		log.Print(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

//api administrative actions
func getAdministrativeActionsHandler(c *gin.Context) {
	res, err := getAdministrativeActions()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"actions": res})
}
