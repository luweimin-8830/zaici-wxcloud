package model

import "time"

// Match 对应 match_demo 表
type Match struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	Channel   string    `gorm:"column:channel;index" json:"channel"`
	OpenId1   string    `gorm:"column:open_id_1;index" json:"openId1"`
	OpenId2   string    `gorm:"column:open_id_2;index" json:"openId2"`
	Status    int       `gorm:"column:status" json:"status"`
	LikeType  int       `gorm:"column:like_type" json:"likeType"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
