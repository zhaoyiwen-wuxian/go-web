package models

import "time"

type Group struct {
	ID        uint
	Name      string
	OwnerID   uint
	CreatedAt time.Time
}

func (table *Group) TableName() string {
	return "group"
}
