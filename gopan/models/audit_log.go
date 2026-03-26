package models

import "time"

type AuditLog struct {
	Id             uint64
	Identity       string
	ActorIdentity  string
	ActorName      string
	ActorRole      int
	Action         string
	TargetType     string
	TargetIdentity string
	Detail         string
	CreatedAt      time.Time `xorm:"created"`
	UpdatedAt      time.Time `xorm:"updated"`
	DeletedAt      time.Time `xorm:"deleted"`
}

func (table AuditLog) TableName() string {
	return "audit_log"
}
