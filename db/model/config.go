package model

import "time"

// Config 对应 config 表
type Config struct {
	ID        uint       `gorm:"primaryKey;column:id" json:"id"`
	Distance  int        `gorm:"column:distance" json:"distance"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}
