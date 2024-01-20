package mysql

import (
	"time"
)

// UserMsg represents the user_msg table
type UserMsg struct {
	ID         uint64    `gorm:"column:id;primaryKey"`
	MsgID      string    `gorm:"column:msg_id;uniqueIndex"`
	ChatID     string    `gorm:"column:chat_id"`
	ChatName   string    `gorm:"column:chat_name"`
	SenderID   string    `gorm:"column:sender_id"`
	SenderName string    `gorm:"column:sender_name"`
	MsgType    string    `gorm:"column:msg_type"`
	Content    string    `gorm:"column:content"`
	Interact   bool      `gorm:"column:interact"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

// RepliedMsg represents the replied_msg table
type RepliedMsg struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	MsgID     string    `gorm:"column:msg_id;uniqueIndex"`
	UserMsgID string    `gorm:"column:user_msg_id"`
	Platform  string    `gorm:"column:platform"`
	Model     string    `gorm:"column:model"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (UserMsg) TableName() string {
	return "user_msg"
}

func (RepliedMsg) TableName() string {
	return "replied_msg"
}
