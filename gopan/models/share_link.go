package models

import "time"

// TODO: 密码
type ShareLink struct {
	Id                 int
	Identity           string
	UserIdentity       string
	RepositoryIdentity string
	Expires            int
	ClickNum           int
	CreatedAt          time.Time `xorm:"created"`
	UpdatedAt          time.Time `xorm:"updated"`
	DeletedAt          time.Time `xorm:"deleted"`
}

func (table ShareLink) TableName() string {
	return "share_link"
}
