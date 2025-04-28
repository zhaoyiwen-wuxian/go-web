package models

import (
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string
	Email         string
	Identity      string
	ClentIp       string
	ClentPort     string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time
	DeviceInfo    string
	IsOnline      bool
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
