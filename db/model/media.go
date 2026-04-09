package model

import "time"

// PictureList 对应 picture_list_demo 表
type PictureList struct {
	ID             uint      `gorm:"primaryKey;column:id" json:"id"`
	UserPicURL     string    `gorm:"column:user_pic_url" json:"userPicUrl"`
	SecCheckStatus int       `gorm:"column:sec_check_status" json:"secCheckStatus"`
	TraceID        string    `gorm:"column:trace_id" json:"traceId"`
	PicHash        string    `gorm:"column:pic_hash;index" json:"picHash"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName 指定表名
func (PictureList) TableName() string {
	return "picture_list_demo"
}
