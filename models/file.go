package models

import (
	"errors"
	"time"
)

type File struct {
	Id     int64  `json:"id"`
	Key    string `json:key xorm:"UNIQUE NOT NULL"`
	Path   string `json:"path"`
	Name   string `json:"name"`
	Ext    string `json:"ext"`
	Size   int    `json:"size"`
	UserId int64  `json:"user_id"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`

	URL string `xorm:"-" json:"url"`
}

func (f *File) CheckValid() error {
	if f.UserId == 0 {
		return errors.New("user_id required")
	}

	return nil
}
