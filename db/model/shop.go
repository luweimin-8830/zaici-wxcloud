package model

import "time"

type Shop struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ShopName  string    `json:"shopname"`
	Location  string    `json:"location"` // GeoJSON string
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// Additional fields from shop_list_demo as needed
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Image     string    `json:"image"`
	StartTime string    `json:"startTime"`
	EndTime   string    `json:"endTime"`
}
