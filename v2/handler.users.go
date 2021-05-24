// handlers.user.go

package v2

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type (
	userField struct {
		Field string `json:"field"`
		Data  string `json:"data"`
		Id    int    `json:"id"`
	}
	userEmail struct {
		Email string `json:"email" binding:"required"`
	}
)

func userCabinetHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := claims["id"].(string)
	userData, err := getUser(user)
	if err != nil {
		c.JSON(200, gin.H{"data": err})
	} else {
		c.JSON(200, gin.H{"data": userData})
	}
}

func registrationHandler(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var user User
	if err := json.Unmarshal(x, &user); err != nil {
		c.AbortWithStatus(400)
		return
	}

	if name, err := user.createUser(); err == nil {
		c.JSON(200, gin.H{"title": "User registered", "id": name})
	} else {
		c.JSON(400, gin.H{"title": err.Error()})
		c.AbortWithStatus(400)
	}
}

func editUserFieldHandler(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var field userField
	if err := json.Unmarshal(x, &field); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := updateUser(field.Field, field.Data, field.Id); err == nil {
		c.JSON(200, gin.H{"title": "Field modifiyed", "data": "kek"})
	} else {
		c.AbortWithStatus(400)
	}
}

func setNewPasswordHandler(c *gin.Context) {
	var pass map[string]string
	email := c.Param("hash")
	if err := c.BindJSON(&pass); err != nil {
		c.JSON(400, gin.H{"title": "Невірний параметр"})
		return
	}

	if err := passwordResetter(pass["password"], email); err == nil {
		c.JSON(200, gin.H{"title": "Пароль змінено"})
	} else {
		c.JSON(400, gin.H{"title": "Невірне посилання"})
	}
}

func activateAccountHandler(c *gin.Context) {
	hash := c.Param("hash")
	var user User

	if usr, err := user.readHash(hash); err == nil {
		if err = user.deleteHash(hash); err != nil {
			c.JSON(400, gin.H{"title": err})
			return
		}

		if err = setActiveField(usr.Email); err != nil {
			c.JSON(400, gin.H{"title": err})
			return
		}
		c.JSON(200, gin.H{"title": "ok"})
	} else {
		if err := c.AbortWithError(404, errors.New("код застарілий")); err != nil {
			log.Error(err)
		}
	}
}

func resetPasswordHandler(c *gin.Context) {
	var user userEmail
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var u User
	u.Email = user.Email
	u.Name = user.Email
	hash, err := u.writeHash()
	if err != nil {
		if err := c.AbortWithError(404, errors.New("код застарілий")); err != nil {
			log.Error(err)
		}
		return
	}

	if err = doSendEmail(u, hash, "email_password"); err != nil {
		if err := c.AbortWithError(404, err); err != nil {
			log.Error(err)
		}
	}
}

func passwordCheckHandler(c *gin.Context) {
	var u User
	hash, err := u.readHash(c.Param("hash"))
	if hash == nil {
		if err := c.AbortWithError(404, err); err != nil {
			log.Error(err)
		}
	} else {
		c.Status(200)
	}
}
