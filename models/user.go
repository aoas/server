package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aoas/server/utils"
)

var (
	ErrPasswordNotMatch = errors.New("The password not matched")
)

type User struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Username string `xorm:"UNIQUE NOT NULL" json:"username"`
	Password string `xorm:"NOT NULL" json:"-"`
	Email    string `xorm:"UNIQUE NOT NULL" json:"email"`

	// Permissions
	IsActive bool `json:"is_active"`
	IsAdmin  bool `json:"is_admin"`

	// 头像
	Avatar string `xorm:"VARCHAR(1024)" json:"avatar"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`

	Token     string `xorm:"-" json:"token,omitempty"`
	ExpiredIn int64  `xorm:"-" json:"expired_in,omitempty"`
}

func (u *User) CheckValid() error {
	if strings.Trim(u.Username, " ") == "" {
		return errors.New("username required")
	}
	if strings.Trim(u.Email, " ") == "" {
		return errors.New("email required")
	}

	return nil
}

func (u *User) EncodePassword() error {
	u.Password = utils.EncodeMd5(u.Password)
	return nil
}

//  compare user password
func (u *User) IsValidPassword(pass string) bool {
	newUser := &User{Password: pass}
	newUser.EncodePassword()

	return u.Password == newUser.Password
}

// 更新用户密码
func (u *User) UpdatePassword(old string, pass string) error {
	if !u.IsValidPassword(old) {
		return ErrPasswordNotMatch
	}
	user := &User{Password: pass}
	user.EncodePassword()
	return UpdateById(u.Id, user, "password")
}

// --- role / permission related
func (u *User) Roles() ([]Role, error) {
	var roles []Role
	sql := fmt.Sprintf("select u.* from role u, role_user r where u.id = r.role_id and r.user_id = %d", u.Id)
	err := x.Sql(sql).Find(&roles)
	return roles, err
}

// 检查用户是否有对某个功能有权限
func (u *User) IsGranted(permissionId string) bool {
	roles, _ := u.Roles()

	for _, c := range roles {
		if c.IsGrantedById(permissionId) {
			return true
		}
	}

	return false
}

// ----------
func GetUserByUserName(username string) *User {
	var user User
	has, err := x.Where("username = ?", strings.ToLower(username)).Get(&user)
	if !has || err != nil {
		return nil
	}
	return &user
}

func CreateUser(user *User) error {
	if err := user.CheckValid(); err != nil {
		return err
	}
	if IsUserNameExist(user.Username) {
		return fmt.Errorf("username is exists: [username: %s]", user.Username)
	}
	if IsEmailUsed(user.Email) {
		return fmt.Errorf("email has been used: [email: %s]", user.Email)
	}
	if user.Nickname == "" {
		user.Nickname = user.Username
	}
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	user.IsActive = true

	user.EncodePassword()

	if err := Insert(user); err != nil {
		return err
	}

	return nil
}

// 检查用户名是否已被注册
func IsUserNameExist(username string) bool {
	username = strings.ToLower(username)
	has, _ := x.Get(&User{Username: username})
	return has
}

// 邮箱地址是否被使用
func IsEmailUsed(email string) bool {
	if strings.Trim(email, " ") == "" {
		return false
	}
	email = strings.ToLower(email)
	has, _ := x.Where("email = ?", strings.ToLower(email)).Get(&User{})
	return has
}

// 用户名是否正确
func IsValidUserName(name string) error {
	return nil
}
