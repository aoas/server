package models

import (
	"fmt"
	"time"
)

type Permission struct {
	Id        string    `xorm:"pk" json:"id"`
	Name      string    `json:"name"`
	Group     string    `json:"group"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

type PermissionRole struct {
	RoleId       int64  `xorm:"pk"`
	PermissionId string `xorm:"pk"`
}

func (p *Permission) AddRoles(roles []*Role) error {
	ids := make([]int64, len(roles))
	for _, c := range roles {
		ids[0] = c.Id
	}

	return p.AddRolesByIds(ids)
}

func (p *Permission) AddRolesByIds(ids []int64) error {
	data := make([]PermissionRole, len(ids))
	for i, c := range ids {
		data[i].PermissionId = p.Id
		data[i].RoleId = c
	}

	return Insert(&data)
}

func (p *Permission) Roles() (roles []*Role, err error) {
	roles = make([]*Role, 0)
	sql := fmt.Sprintf("select u.* from role u, permission_role r where u.id = r.role_id and r.permission_id= %s", p.Id)
	err = x.Sql(sql).Find(roles)
	return
}

func AddPermissionsByList(permissions []Permission) error {
	for _, c := range permissions {
		// check exist
		if IsExist(&c) {
			continue
		}
		if err := Insert(&c); err != nil {
			return err
		}
	}

	return nil

}
