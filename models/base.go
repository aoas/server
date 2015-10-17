package models

import (
	"errors"

	"github.com/go-xorm/xorm"
)

var (
	ErrNotExist = errors.New("not exist")
)

func Engine() *xorm.Engine {
	return x
}

// 打开缓存, 默认只缓存User, Role, Task相关表
func EnableCache() error {
	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	x.MapCacher(&User{}, cacher)
	x.MapCacher(&Role{}, cacher)
	x.MapCacher(&RoleUser{}, cacher)
	x.MapCacher(&Permission{}, cacher)
	x.MapCacher(&PermissionRole{}, cacher)

	return nil
}

// 同步表结构到数据库
func SyncTables() error {
	err := x.Sync2(
		new(User),
		new(Role),
		new(Permission),
		new(RoleUser),
		new(PermissionRole),
		new(File),
	)

	return err

}

func InitDefaultData() error {

	return nil
}
