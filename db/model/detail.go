package model

import "time"

// DetailRecord 对应 detail_record_demo 表
type DetailRecord struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	OpenID      string    `gorm:"column:open_id;index" json:"openId"`
	Plan        string    `gorm:"column:plan" json:"plan"`
	Mood        string    `gorm:"column:mood" json:"mood"`
	Style       string    `gorm:"column:style" json:"style"`
	Description string    `gorm:"column:description" json:"description"`
	Seat        string    `gorm:"column:seat" json:"seat"`
	Image       string    `gorm:"column:image" json:"image"`
	Status      int       `gorm:"column:status;index" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// PictureList 对应 picture_list_demo 表
type PictureList struct {
	ID             uint   `gorm:"primaryKey;column:id" json:"id"`
	UserPicURL     string `gorm:"column:user_pic_url" json:"userPicUrl"`
	SecCheckStatus int    `gorm:"column:sec_check_status" json:"secCheckStatus"`
	TraceID        string `gorm:"column:trace_id" json:"traceId"`
	PicHash        string `gorm:"column:pic_hash;index" json:"picHash"`
}
