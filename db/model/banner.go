package model

import "time"

// Banner 对应 banner_demo 表
type Banner struct {
	ID        uint       `gorm:"primaryKey;column:id" json:"id"`
	Title     string     `gorm:"column:title" json:"title"`
	ImageUrl  string     `gorm:"column:image_url" json:"imageUrl"`
	LinkUrl   string     `gorm:"column:link_url" json:"linkUrl"`
	Type      string     `gorm:"column:type" json:"type"`
	Url       string     `gorm:"column:url" json:"url"`
	Interval  int        `gorm:"column:interval" json:"interval"`
	Detail    string     `gorm:"column:detail;type:text" json:"detail"` // 存储 JSON 字符串
	CreatedAt time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}
