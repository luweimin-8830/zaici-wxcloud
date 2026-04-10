package model

import "time"

type Shop struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ShopName  string    `gorm:"column:shopname" json:"shopname"`
	Location  string    `gorm:"column:location" json:"location"` // GeoJSON string
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
	// Additional fields from shop_list_demo as needed
	Address   string    `gorm:"column:address" json:"address"`
	Phone     string    `gorm:"column:phone" json:"phone"`
	Image     string    `gorm:"column:image" json:"image"`
	StartTime string    `gorm:"column:start_time" json:"startTime"`
	EndTime   string    `gorm:"column:end_time" json:"endTime"`
	Tag1      string    `gorm:"column:tag1" json:"tag1"`
	Tag2      string    `gorm:"column:tag2" json:"tag2"`
}
