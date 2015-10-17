package controllers

import (
	"fmt"
	"strconv"

	"github.com/aoas/server/config"
	"github.com/aoas/server/models"
	"github.com/aoas/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type Base struct {
	Engine *xorm.Engine
	Config config.Config
	Logger utils.ILogger
}

var permissions = []models.Permission{}

func PermissionList() []models.Permission {
	return permissions
}

func addPermission(group string, id string, name string) {
	permissions = append(permissions, models.Permission{
		Group: group,
		Id:    id,
		Name:  name,
	})
}

func isGranted(c *gin.Context, permissionId string) bool {
	validPermission := false
	for _, p := range permissions {
		if p.Id == permissionId {
			validPermission = true
		}
	}
	if !validPermission {
		return false
	}

	cu := currentUser(c)
	if cu == nil {
		return false
	}

	return cu.IsGranted(permissionId)
}

func currentUser(c *gin.Context) *models.User {
	sid, _ := c.Get("userid")
	if sid == "" {
		return nil
	}
	id, _ := strconv.ParseInt(fmt.Sprintf("%v", sid), 10, 64)
	var user models.User
	if err := models.GetById(id, &user); err != nil {
		fmt.Println("get user error:", err.Error())
		return nil
	}

	return &user
}

func queryInt(c *gin.Context, name string, def ...int) int {
	sv := c.Query(name)
	if sv == "" {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}

	i, _ := strconv.ParseInt(sv, 10, 32)
	return int(i)
}

func queryInt64(c *gin.Context, name string, def ...int64) int64 {
	sv := c.Query(name)
	if sv == "" {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}

	i, _ := strconv.ParseInt(sv, 10, 64)
	return i
}

func paramInt64(c *gin.Context, name string, def ...int64) int64 {
	sv := c.Param(name)
	if sv == "" {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}

	i, _ := strconv.ParseInt(sv, 10, 64)
	return i
}
