package models

import "time"

type User struct {
	Id                 uint64
	Identity           string
	Name               string
	Password           string
	Email              string
	Status             int
	Role               int
	UploadPermission   int
	DownloadPermission int
	SharePermission    int
	LastLoginAt        time.Time
	CreatedAt          time.Time `xorm:"created"`
	UpdatedAt          time.Time `xorm:"updated"`
	DeletedAt          time.Time `xorm:"deleted"`
}

func (table User) TableName() string {
	return "user_basic"
}
