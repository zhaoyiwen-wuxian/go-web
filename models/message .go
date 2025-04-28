package models

import "gorm.io/gorm"

type MessageType string

const (
	TextMessage     MessageType = "text"
	ImageMessage    MessageType = "image"
	EmojiMessage    MessageType = "emoji"
	AudioMessage    MessageType = "audio"
	VideoCallInvite MessageType = "video_invite"
	VoiceCallInvite MessageType = "voice_invite"
	CallAccept      MessageType = "call_accept"
	CallReject      MessageType = "call_reject"
	CallCandidate   MessageType = "call_candidate"
	CallOffer       MessageType = "call_offer"
	CallAnswer      MessageType = "call_answer"
	CallLeave       MessageType = "call_leave"
)

type Message struct {
	gorm.Model
	Type      string `json:"type"`      // 消息类型
	From      string `json:"from"`      // 发送者ID
	To        string `json:"to"`        // 接收者ID或群ID
	Content   string `json:"content"`   // 内容
	Timestamp int64  `json:"timestamp"` // 时间戳
	Extra     string `json:"extra"`     // 拓展字段
}

func (table *Message) TableName() string {
	return "message"
}
