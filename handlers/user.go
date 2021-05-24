package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	mail "mtest.com.ua/mail"

	datapkg "mtest.com.ua/db/dataprocessor"
	hashpkg "mtest.com.ua/db/hasher"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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

	hasher interface {
		WriteHash(email string) (string, error)
		ReadHash(hash string) (hashpkg.HashData, error)
		DeleteHash(hash string) (err error)
	}
)

func (hd *Handlers) UserCabinetHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	u, ok := claims["id"]
	if !ok {
		logrus.Error("no ID in request")
		c.AbortWithStatus(400)
		return
	}
	user, ok := u.(string)
	if !ok {
		logrus.Error("no ID is not a string")
		c.AbortWithStatus(400)
		return
	}
	userData, err := hd.GetUser(user)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{"data": err})
	} else {
		c.JSON(200, gin.H{"data": userData})
	}
}

func (hd *Handlers) RegistrationHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(400)
		return
	}
	var user datapkg.User
	if err := json.Unmarshal(x, &user); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(400)
		return
	}

	if name, err := hd.CreateUser(); err == nil {
		h, err := hd.WriteHash(user.Email)
		if err != nil {
			log.Error(err)
			c.AbortWithStatus(400)
			return
		}

		if err := mail.SendEmail(user.Name, user.Email, h, "email_activate"); err != nil {
			log.Error(err)
			c.AbortWithStatus(400)
			return
		}

		c.JSON(200, gin.H{"title": "User registered", "id": name})
	} else {
		logrus.Error(err)
		c.JSON(400, gin.H{"title": err.Error()})
		c.AbortWithStatus(400)
	}
}

func (hd *Handlers) EditUserFieldHandler(c *gin.Context) {
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatus(400)
		return
	}
	var field userField
	if err := json.Unmarshal(x, &field); err != nil {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := hd.UpdateUser(field.Field, field.Data, field.Id); err == nil {
		c.JSON(200, gin.H{"title": "Field modifiyed", "data": "kek"})
	} else {
		logrus.Error(err)
		c.AbortWithStatus(400)
	}
}

func (hd *Handlers) SetNewPasswordHandler(c *gin.Context) {
	var pass map[string]string
	hsh := c.Param("hash")
	if err := c.BindJSON(&pass); err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{"title": "Невірний параметр"})
		return
	}

	defer func() {
		if err := hd.DeleteHash(hsh); err != nil {
			log.Error(err)
		}
	}()

	if err := hd.UpdatePassword(pass["password"], pass["email"], hsh); err == nil {
		c.JSON(200, gin.H{"title": "Пароль змінено"})
	} else {
		logrus.Error(err)
		c.JSON(400, gin.H{"title": "Невірне посилання"})
	}
}

func (hd *Handlers) ActivateAccountHandler(c *gin.Context) {
	if usr, err := hd.ReadHash(c.Param("hash")); err == nil {
		if err = hd.DeleteHash(c.Param("hash")); err != nil {
			logrus.Error(err)
			c.JSON(400, gin.H{"title": err})
			return
		}

		if err = hd.SetActiveField(usr.Email); err != nil {
			logrus.Error(err)
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

func (hd *Handlers) ResetPasswordHandler(c *gin.Context) {
	var user userEmail
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	hashData, err := hd.WriteHash(user.Email)
	if err != nil {
		logrus.Error(err)
		if err := c.AbortWithError(404, errors.New("код застарілий")); err != nil {
			logrus.Error(err)
		}
		return
	}

	if err = mail.SendEmail(user.Email, user.Email, hashData, "email_password"); err != nil {
		logrus.Error(err)
		if err := c.AbortWithError(404, err); err != nil {
			log.Error(err)
		}
	}
}

func (hd *Handlers) PasswordCheckHandler(c *gin.Context) {
	hashData, err := hd.ReadHash(c.Param("hash"))
	if err != nil {
		if err := c.AbortWithError(404, err); err != nil {
			logrus.Error(err)
		}
		return
	}
	if hashData.Hash == "" {
		if err := c.AbortWithError(404, errors.New("hash is empty")); err != nil {
			logrus.Error(err)
		}
	} else {
		c.Status(200)
	}
}
