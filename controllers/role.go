package controllers

import (
	"github.com/aoas/server/models"
	"github.com/aoas/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type Role struct {
	Base
}

func init() {
	addPermission("role", "role.list", "角色列表")
	addPermission("role", "role.create", "新增角色")
	addPermission("role", "role.delete", "删除角色")
	addPermission("role", "role.users", "角色用户列表")
	addPermission("role", "role.adduser", "增加角色用户")
	addPermission("role", "role.deleteuser", "删除角色用户")
	addPermission("role", "role.permissions", "角色权限列表")
	addPermission("role", "role.addpermission", "增加角色权限")
	addPermission("role", "role.deletepermission", "删除角色权限")
}
func (r *Role) Create(c *gin.Context) {
	if !isGranted(c, "role.create") {
		c.JSON(403, utils.NewNoAccessPermissionError(""))
		return
	}
	cu := currentUser(c)
	var role models.Role
	if err := c.BindJSON(&role); err != nil {
		c.JSON(400, utils.NewInvalidJsonError())
		return
	}

	role.UserId = cu.Id

	if err := role.CheckValid(); err != nil {
		c.JSON(400, utils.NewError("invalid data - %s", err.Error()))
		return
	}

	if err := models.Insert(&role); err != nil {
		c.JSON(400, utils.NewError("add role failed - %s", err.Error()))
		return
	}

	c.JSON(200, role)

}
func (r *Role) List(c *gin.Context) {
	if !isGranted(c, "role.list") {
		c.JSON(403, utils.NewNoAccessPermissionError(""))
		return
	}
	//userid, _ := c.Get("userid")
	page := queryInt(c, "page", 1)
	pagesize := queryInt(c, "pagesize", 30)

	engine := models.Engine()
	session := func() *xorm.Session {
		where := engine.Where("name like ?", "%"+c.Query("name")+"%")
		return where
	}

	total, _ := session().Count(&models.Role{})
	roles := make([]models.Role, 0)

	if err := session().Limit(pagesize, pagesize*(page-1)).Asc("created_at").Find(&roles); err != nil {
		c.JSON(400, utils.NewError("find roles failed - %s", err.Error()))
		return
	}

	result := models.NewQueryResult(page, pagesize, total, roles)

	c.JSON(200, result)
}
func (r *Role) Delete(c *gin.Context) {
	if !isGranted(c, "role.delete") {
		c.JSON(403, utils.NewNoAccessPermissionError(""))
		return
	}
	id := paramInt64(c, "id")

	role := models.Role{Id: id}

	if err := role.Delete(); err != nil {
		c.JSON(400, utils.NewError("delete role failed - %s", err.Error()))
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}

func (r *Role) Users(c *gin.Context) {
	if !isGranted(c, "role.users") {
		c.JSON(403, utils.NewNoAccessPermissionError(""))
		return
	}
	id := paramInt64(c, "id")

	var role models.Role
	if err := models.GetById(id, &role); err != nil {
		c.JSON(400, utils.NewNotFoundError())
		return
	}

	users, _ := role.Users()

	c.JSON(200, users)
}

func (r *Role) AddUser(c *gin.Context) {
	if !isGranted(c, "role.adduser") {
		c.JSON(403, utils.NewNoAccessPermissionError(""))
		return
	}
	id := paramInt64(c, "id")
	var role models.Role
	if err := models.GetById(id, &role); err != nil {
		c.JSON(400, utils.NewNotFoundError())
		return
	}

	var data struct {
		UserIds []int64 `json:"user_ids"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, utils.NewInvalidJsonError())
		return
	}

	if len(data.UserIds) <= 0 {
		c.JSON(400, utils.NewError("user_ids required"))
		return
	}

	if err := role.AddUserByIds(data.UserIds...); err != nil {
		c.JSON(400, utils.NewError("add user to role failed - %s", err.Error()))
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})

}

func (r *Role) DeleteUsers(c *gin.Context) {
	if !isGranted(c, "role.deleteuser") {
		c.JSON(403, utils.NewNoAccessPermissionError(""))
		return
	}
	id := paramInt64(c, "id")
	var role models.Role

	if err := models.GetById(id, &role); err != nil {
		c.JSON(400, utils.NewNotFoundError())
		return
	}

	var data struct {
		UserIds []int64 `json:"user_ids"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, utils.NewInvalidJsonError())
		return
	}

	if len(data.UserIds) <= 0 {
		c.JSON(400, utils.NewError("user_ids required"))
		return
	}

	if err := role.DeleteUserByIds(data.UserIds...); err != nil {
		c.JSON(400, utils.NewError("delete user by role failed - %s", err.Error()))
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})

}

func (r *Role) Permissions(c *gin.Context) {
	if !isGranted(c, "role.permissions") {
		c.JSON(403, utils.NewNoAccessPermissionError(""))
		return
	}
	id := paramInt64(c, "id")
	var role models.Role

	if err := models.GetById(id, &role); err != nil {
		c.JSON(400, utils.NewNotFoundError())
		return
	}

	permissions, err := role.Permissions()
	if err != nil {
		c.JSON(400, utils.NewError("get role permissions failed - %s", err.Error()))
		return
	}

	c.JSON(200, permissions)

}

func (r *Role) AddPermissions(c *gin.Context) {
	if !isGranted(c, "role.addpermission") {
		c.JSON(403, utils.NewNoAccessPermissionError(""))
		return
	}
	id := paramInt64(c, "id")
	var role models.Role

	if err := models.GetById(id, &role); err != nil {
		c.JSON(400, utils.NewNotFoundError())
		return
	}

	var data struct {
		PermissionIds []string `json:"permission_ids"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, utils.NewInvalidJsonError())
		return
	}

	if len(data.PermissionIds) <= 0 {
		c.JSON(400, utils.NewError("permission_ids required"))
		return
	}
	if err := role.AddPermissionsByIds(data.PermissionIds...); err != nil {
		c.JSON(400, utils.NewError("add role permissions failed - %s", err.Error()))
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})

}

func (r *Role) DeletePermissions(c *gin.Context) {
	if !isGranted(c, "role.deletepermission") {
		c.JSON(403, utils.NewNoAccessPermissionError(""))
		return
	}
	id := paramInt64(c, "id")
	var role models.Role

	if err := models.GetById(id, &role); err != nil {
		c.JSON(400, utils.NewNotFoundError())
		return
	}

	var data struct {
		PermissionIds []string `json:"permission_ids"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, utils.NewInvalidJsonError())
		return
	}

	if len(data.PermissionIds) <= 0 {
		c.JSON(400, utils.NewError("permission_ids required"))
		return
	}

	if err := role.DeletePermissions(data.PermissionIds...); err != nil {
		c.JSON(400, utils.NewError("delete role permissions failed - %s", err.Error()))
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})

}
