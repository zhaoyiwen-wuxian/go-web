package models

type GroupMember struct {
	ID      uint
	GroupID uint
	UserID  uint
	Role    string // admin / member
}

func (table *GroupMember) TableName() string {
	return "group_member"
}
