package model

import "time"

// ChatHistory 对应 new_chat_history_demo 表
type ChatHistory struct {
	ID             uint       `gorm:"primaryKey;column:id" json:"id"`
	ChannelId      string     `gorm:"column:channel_id;index" json:"channelId"`
	SenderOpenID   string     `gorm:"column:sender_open_id" json:"senderOpenID"`
	ReceiverOpenID string     `gorm:"column:receiver_open_id" json:"receiverOpenID"`
	MessageContent MessageContent `gorm:"column:message_content;type:text" json:"messageContent"`
	Timestamp      int64      `gorm:"column:timestamp;index" json:"timestamp"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
}

type MessageContent struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	Type     string `json:"type"`
	ContentType string `json:"contentType"`
	Pic      string `json:"pic"`
	Name     string `json:"name"`
	State    int    `json:"state"`
}

// BlockList 对应 block_list_demo 表
type BlockList struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	OpenID    string    `gorm:"column:open_id;index" json:"openId"`
	BlockID   string    `gorm:"column:block_id;index" json:"blockId"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

// InfoMonitor 对应 information_monitor_demo 表
type InfoMonitor struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	OpenID    string    `gorm:"column:open_id;index" json:"openId"`
	Source    string    `gorm:"column:source" json:"source"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
}
