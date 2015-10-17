package controllers

import (
	"fmt"
	"strings"

	"github.com/aoas/server/models"
	"github.com/aoas/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type User struct {
	Base
}

func init() {
	addPermission("user", "user.list", "查询用户列表")
	addPermission("user", "user.get", "获取用户")
	addPermission("user", "user.delete", "删除用户")
	addPermission("user", "user.active", "启用/禁用用户账号")
	addPermission("user", "user.roles", "获取用户所在组别列表")
}

// liist all user
// is_active 1/0
func (u *User) List(c *gin.Context) {
	if !isGranted(c, "user.list") {
		c.JSON(403, utils.NewNoAccessPermissionError(""))
		return
	}

	page := queryInt(c, "page", 1)
	pagesize := queryInt(c, "pagesize", 30)
	isActive := strings.ToLower(c.Query("is_active"))

	engine := models.Engine()
	session := func() *xorm.Session {
		where := engine.Where("nickname like ? and username like ? and email like ? ",
			"%"+c.Query("nickname")+"%", "%"+c.Query("username")+"%", "%"+c.Query("email")+"%")

		if isActive != "" {
			if isActive == "1" || isActive == "true" {
				where = where.Where("is_active = ?", true)
			} else if isActive == "0" || isActive == "false" {
				where = where.Where("is_active = ?", false)
			}
		}

		return where
	}

	total, _ := session().Count(&models.User{})
	users := make([]models.User, 0)

	if err := session().Limit(pagesize, pagesize*(page-1)).Asc("created_at").Find(&users); err != nil {
		c.JSON(400, utils.NewError("find users failed - %s", err.Error()))
		return
	}

	result := models.NewQueryResult(page, pagesize, total, users)

	c.JSON(200, result)

}

func (u *User) Get(c *gin.Context) {
	//cu := currentUser(c)
	id := paramInt64(c, "id")

	var user models.User
	if err := models.GetById(id, &user); err != nil {
		c.JSON(400, utils.NewNotFoundError())
		return
	}

	c.JSON(200, user)

}

// disable user
func (u *User) Active(c *gin.Context) {
	if !isGranted(c, "user.active") {
		c.JSON(403, utils.NewNoAccessPermissionError(""))
		return
	}

	id := paramInt64(c, "id")
	var user models.User
	var data models.User
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, utils.NewInvalidJsonError())
		return
	}

	if err := models.GetById(id, &user); err != nil {
		c.JSON(400, utils.NewNotFoundError())
		return
	}

	fmt.Println("is active", data.IsActive)
	user.IsActive = data.IsActive
	if err := models.UpdateById(id, &user, "is_active"); err != nil {
		c.JSON(400, utils.NewError("update database failed - %s", err.Error()))
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})

}

func (u *User) Roles(c *gin.Context) {
	//cu := currentUser(c)
	id := paramInt64(c, "id")

	user := models.User{Id: id}
	roles, err := user.Roles()
	if err != nil {
		c.JSON(400, utils.NewError("get user roles failed - %s", err.Error()))
		return
	}

	c.JSON(200, roles)

}
