// handlers.user.go

package main

import (
	"io/ioutil"
	"encoding/json"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

)

type newUser struct {
	Name, Surename, Email, Password string
}
type userField struct {
	Field string `json:"field"`
	Data string `json:"data"`
	Id int `json:"id"`
}


func cabinetHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := claims["id"].(string)

	userdata, err := readUser(user)
	if err != nil {
		c.JSON(200, gin.H{
			"data": err,
		})
	} else {
		c.JSON(200, gin.H{
			"data": userdata,
		})
	}
}

func registerHandler(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var user newUser
	err := json.Unmarshal([]byte(x), &user)
	check(err)

	if name, err := createUser(user.Name,user.Surename,user.Email,user.Password); err == nil {
		c.JSON(200, gin.H{
			"title": "User registered",
			"id": name,
		})
	} else {
		c.AbortWithStatus(400)
	}
}

func editUserField(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var field userField
	err := json.Unmarshal([]byte(x), &field)
	check(err)
	if err := updateUser(field.Field, field.Data, field.Id); err == nil {
		c.JSON(200, gin.H{
			"title": "Field modifiyed",
			"data": "kek",
		})

	} else {
		c.AbortWithStatus(400)
	}
}