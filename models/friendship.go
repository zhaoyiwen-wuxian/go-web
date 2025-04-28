package models

import "time"

type Friendship struct {
	ID        uint
	UserID    uint   // 主用户
	FriendID  uint   // 好友用户
	Status    string // apply / accepted / blocked / deleted
	CreatedAt time.Time
}

func (table *Friendship) TableName() string {
	return "friendship"
}
