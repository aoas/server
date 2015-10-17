package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Role struct {
	Id          int64
	Name        string `xorm:"UNIQUE NOT NULL"`
	UserId      int64  `json:"user_id"`
	Description string
	CreatedAt   time.Time `xorm:"created" json:"created_at"`
	UpdatedAt   time.Time `xorm:"updated" json:"updated_at"`
}

type RoleUser struct {
	RoleId int64 `xorm:"pk"`
	UserId int64 `xorm:"pk"`
}

func (r *Role) CheckValid() error {
	if strings.Trim(r.Name, " ") == "" {
		return errors.New("name required")
	}
	if r.UserId <= 0 {
		return errors.New("create user required")
	}

	return nil
}

// 删除role及相关数据
func (r *Role) Delete() error {
	s := x.NewSession()

	// delete permission
	if _, err := s.Where("role_id = ?", r.Id).Delete(&PermissionRole{}); err != nil {
		s.Rollback()
		return err
	}

	// delete user role list
	if _, err := s.Where("role_id = ?", r.Id).Delete(&RoleUser{}); err != nil {
		s.Rollback()
		return err
	}

	// delete role
	if _, err := s.Where("id = ?", r.Id).Delete(&Role{}); err != nil {
		s.Rollback()
		return err
	}

	s.Commit()

	return nil

}

// ---- users ----
func (r *Role) AddUser(user *User) error {
	obj := &RoleUser{RoleId: r.Id, UserId: user.Id}
	return Insert(&obj)
}

func (r *Role) AddUserByIds(ids ...int64) error {
	if len(ids) == 0 {
		return errors.New("required ids")
	}

	data := make([]RoleUser, len(ids))
	for i, c := range ids {
		data[i].RoleId = r.Id
		data[i].UserId = c
	}

	return Insert(&data)
}

func (r *Role) DeleteUserByIds(ids ...int64) error {
	if len(ids) == 0 {
		return errors.New("required ids")
	}
	_, err := x.In("user_id", ids).Delete(&RoleUser{})
	return err
}

func (r *Role) Users() (*[]User, error) {
	var users []User
	sql := fmt.Sprintf("select u.* from user u, role_user r where u.id = r.user_id and r.role_id = %d", r.Id)
	err := x.Sql(sql).Find(&users)
	return &users, err
}

// ---- permissions ----
func (r *Role) Permissions() (*[]Permission, error) {
	permissions := make([]Permission, 0)
	sql := fmt.Sprintf("select u.* from permission u, permission_role r where u.id = r.permission_id and r.role_id = %d", r.Id)
	err := x.Sql(sql).Find(&permissions)
	return &permissions, err
}

func (r *Role) AddPermissionsByIds(ids ...string) error {
	if len(ids) == 0 {
		return errors.New("ids required")
	}

	data := make([]PermissionRole, len(ids))
	for i, c := range ids {
		data[i].RoleId = r.Id
		data[i].PermissionId = c
	}

	return Insert(&data)
}
func (r *Role) DeletePermissions(ids ...string) error {
	if len(ids) == 0 {
		return errors.New("ids required")
	}
	_, err := x.In("permission_id", ids).Delete(&PermissionRole{})
	return err
}

func (r *Role) IsGranted(permission *Permission) bool {
	return r.IsGrantedById(permission.Id)
}

func (r *Role) IsGrantedById(permissionId string) bool {
	all := strings.Split(permissionId, ".")
	has, _ := x.Where("role_id = ? and ( permission_id = ? or permission_id = ? )",
		r.Id,
		permissionId,
		all[0]).Get(&PermissionRole{})
	return has
}
