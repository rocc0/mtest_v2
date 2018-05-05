// handlers.user.go

package main

import (
	"io/ioutil"
	"encoding/json"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"errors"
)

type userField struct {
	Field string `json:"field"`
	Data string `json:"data"`
	Id int `json:"id"`
}

type userEmail struct {
	Email string `json:"email" binding:"required"`
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
	var user User
	err := json.Unmarshal([]byte(x), &user)
	check(err)

	if name, err := user.createUser(); err == nil {
		c.JSON(200, gin.H{
			"title": "User registered",
			"id": name,
		})
	} else {
		c.JSON(400, gin.H{
			"title": err.Error(),
		})
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

func resetPassword(c *gin.Context) {
	var pass map[string]string
	email := c.Param("hash")
	err := c.BindJSON(&pass)
	if err != nil {
		c.JSON(400, gin.H{
			"title": "Невірний параметр",
		})
	}

	if err := passwordResetter(pass["password"], email); err == nil {
		c.JSON(200, gin.H{
			"title": "Пароль змінено",
		})
	} else {
		c.JSON(400, gin.H{
			"title": "Невірне посилання",
		})
	}
}

func activateAccount(c *gin.Context) {
	hash := c.Param("hash")
	var user User

	if usr, err := user.readHash(hash); err == nil {
		err = user.deleteHash(hash)
		if err != nil {
			c.JSON(400, gin.H{
				"title": err,
			})
		}
		err = setActiveField(usr.Email)
		if err != nil {
			c.JSON(400, gin.H{
				"title": err,
			})
		}
		c.JSON(200, gin.H{
			"title": "ok",
		})
	} else {
		c.AbortWithError(404, errors.New("Код застарілий"))
	}
}

func sendResetPasswordLink(c *gin.Context) {

	var user userEmail
	err := c.BindJSON(&user)

	var u User
	u.Email = user.Email
	u.Name = user.Email
	res, err := u.writeHash()
	if err != nil {
		c.AbortWithError(404, errors.New("Код застарілий"))
	}
	err = doSendEmail(u, *res, "email_password")
	if err != nil {
		c.AbortWithError(404, err)
	}
}

func checkPasswordLink(c *gin.Context) {
	var u User
	hash, err := u.readHash(c.Param("hash"))
	if hash == nil {
		c.AbortWithError(404, err)
	} else {
		c.Status(200)
	}
}