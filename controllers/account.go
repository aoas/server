package controllers

import (
	"strings"

	"github.com/aoas/server/models"
	"github.com/aoas/server/utils"
	"github.com/gin-gonic/gin"
)

type loginUser struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (l *loginUser) CheckValid() error {
	if strings.Trim(l.UserName, " ") == "" {
		return utils.NewError("username required")
	}

	if strings.Trim(l.Password, " ") == "" {
		return utils.NewError("password required")
	}

	return nil
}

type Account struct {
	Base
}

func (a *Account) Login(c *gin.Context) {
	var param loginUser
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(400, utils.NewInvalidJsonError())
		return
	}

	if err := param.CheckValid(); err != nil {
		c.JSON(400, err)
		return
	}

	user := models.GetUserByUserName(param.UserName)
	if user == nil {
		c.JSON(400, utils.NewError("user not exist"))
		return
	}

	if !user.IsValidPassword(param.Password) {
		c.JSON(400, utils.NewError("invalid password"))
		return
	}

	login := models.NewLogin(a.Config.TokenSecret, a.Config.TokenExpiredIn)
	token, err := login.GetToken(user)
	if err != nil {
		c.JSON(400, utils.NewError("gen token failed - %d -%s", user.Id, err.Error()))
		return
	}

	user.Token = token
	user.ExpiredIn = a.Config.TokenExpiredIn

	c.JSON(200, user)

}

func (a *Account) Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, utils.NewError(err.Error()))
		return
	}

	if err := models.CreateUser(&user); err != nil {
		c.JSON(400, utils.NewError(err.Error()))
		return
	}

	c.JSON(201, user)

}
