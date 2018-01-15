package main

import (
	"github.com/gin-gonic/gin"

	"net/http"

	"io/ioutil"
	"encoding/json"
	"github.com/appleboy/gin-jwt"
	"log"
)


type newMtest struct {
	Name string 	`json:"name"`
	Region int		`json:"region"`
	Government int 	`json:"government"`
	CalcType int `json:"calc_type"`
}

type deleteRequest struct {
	Id string `json:"id"`
}

type delExecutorReq struct {
	ExEmail string `json:"ex_email"`
	ExMtestId string `json:"ex_mtest_id"`
	DevMtestId string `json:"dev_mtest_id"`
}

type newExecutor struct {
	Title string `json:"title"`
	Email string `json:"email"`
	Region int `json:"region"`
	Government int `json:"government"`
	DevMid string `json:"dev_mid"`
}

// page render
func showIndexPage(c *gin.Context) {
	render(c, gin.H{
		"title":   "Калькулятор",
	}, "index.html")
}

func showUserPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Кабінет користувача",
	}, "index.html")
}

func showMtestPage(c *gin.Context) {
	id := c.Param("mtest_id")
	mtest, err := readMtest(id)
	if err == nil {
		render(c, gin.H{
			"title": "Редагування | " + mtest.Name,
		}, "index.html")
	} else {
		c.AbortWithError(http.StatusNotFound, err)
	}

}

// api mtest
func getReadMtest(c *gin.Context) {
	id := c.Param("mtest_id")
	mtest, err := readMtest(id)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"mtest": mtest,
		})
	} else {
		log.Print(err)
		c.AbortWithError(http.StatusNotFound, err)
	}
}

func postCreateMtest(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	email, _ := claims["id"].(string)

	x, _ := ioutil.ReadAll(c.Request.Body)
	var m newMtest
	json.Unmarshal([]byte(x), &m)

	if data, err := createNewMtest(m, email); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"title": "Item added",
			"records": data,
		})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Print(err)
	}
}

func postUpdateMtest(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	email, _ := claims["id"].(string)

	x, _ := ioutil.ReadAll(c.Request.Body)
	var form map[string]interface{}
	json.Unmarshal([]byte(x), &form)

	if err := updateMtest(form, email); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"title": "Mtest updated",
			"data": form,
		})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func postDeleteMtest(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	email, _ := claims["id"].(string)

	x, _ := ioutil.ReadAll(c.Request.Body)
	var id deleteRequest
	json.Unmarshal([]byte(x), &id)

	log.Print(id.Id)

	if err := deleteMtest(id.Id, email); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"title": "Item removed",
		})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

//api govs
func getGovernments(c *gin.Context) {
	res, err := getGovs()
	check(err)
	c.JSON(http.StatusOK, gin.H{
		"govs": res,
	})
}


func getRegions(c *gin.Context) {
	res, err := getRegs()
	check(err)
	c.JSON(http.StatusOK, gin.H{
		"regions": res,
	})
}


//executors and group calculations
//add executor
func postCreateMtestExecutor(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	email, _ := claims["id"].(string)

	x, _ := ioutil.ReadAll(c.Request.Body)
	var executor newExecutor
	json.Unmarshal([]byte(x), &executor)

	if mid, err := createMtestExecutor(email, executor); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"mid": mid,
		})
	} else {
		log.Print(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}


func postDeleteExecutor(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	email, _ := claims["id"].(string)

	x, _ := ioutil.ReadAll(c.Request.Body)
	var delRequest delExecutorReq
	json.Unmarshal([]byte(x), &delRequest)

	if  err := deleteMtestExecutor(email, delRequest); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"response": "ok",
		})
	} else {
		log.Print(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
//api administrative actions
func getAdmActions(c *gin.Context) {
	res, err := getAdmactions()
	check(err)
	c.JSON(http.StatusOK, gin.H{
		"actions": res,
	})
}