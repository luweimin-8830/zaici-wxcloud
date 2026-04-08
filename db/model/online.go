package model

import "time"

// OnlineRecord 对应 online_demo 表
type OnlineRecord struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	OpenID    string    `gorm:"column:open_id;index" json:"openId"`
	ShopID    string    `gorm:"column:shop_id;index" json:"shopId"`
	ShopName  string    `gorm:"column:shop_name" json:"shopName"`
	Name      string    `gorm:"column:name" json:"name"`
	Avatar    string    `gorm:"column:avatar" json:"avatar"`
	Status    string    `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
	DueTime   int64     `gorm:"column:due_time;index" json:"dueTime"`
	Flag      int       `gorm:"column:flag" json:"flag"`
	Location  string    `gorm:"column:location;type:text" json:"location"` // 存储 GeoJSON 字符串
}
