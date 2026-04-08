package model

import "time"

// User 对应 users_demo 表
type User struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	OpenID      string    `gorm:"column:open_id;type:varchar(128);uniqueIndex" json:"openId"`
	Name        string    `gorm:"column:name" json:"name"`
	Avatar      string    `gorm:"column:avatar" json:"avatar"`
	Company     string    `gorm:"column:company" json:"company"`
	Department  string    `gorm:"column:department" json:"department"`
	Phone       string    `gorm:"column:phone" json:"phone"`
	SuperLike   int       `gorm:"column:super_like" json:"superLike"`
	BeLike      int       `gorm:"column:be_like" json:"beLike"`
	BeSuperLike int       `gorm:"column:be_super_like" json:"beSuperLike"`
	Image       string    `gorm:"column:image;type:text" json:"image"` // 存储JSON字符串或逗号分隔
	Created     time.Time `gorm:"column:created" json:"created"`
	LastLogin   time.Time `gorm:"column:last_login" json:"lastLogin"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// Admin 对应 admins 表
type Admin struct {
	ID     uint   `gorm:"primaryKey;column:id" json:"id"`
	OpenID string `gorm:"column:open_id;type:varchar(128);uniqueIndex" json:"openId"`
}
