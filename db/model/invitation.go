package model

import "time"

// Invitation 对应 invitation_demo 表
type Invitation struct {
	ID             uint       `gorm:"primaryKey;column:id" json:"id"`
	Title          string     `gorm:"column:title" json:"title"`
	ImageUrl       string     `gorm:"column:image_url" json:"imageUrl"`
	Activity       string     `gorm:"column:activity" json:"activity"`
	Status         string     `gorm:"column:status" json:"status"`
	PosterUrl      string     `gorm:"column:poster_url" json:"posterUrl"`
	RequiredFields string     `gorm:"column:required_fields;type:text" json:"requiredFields"` // 存储 JSON 字符串
	StartTime      string     `gorm:"column:start_time" json:"startTime"`
	EndTime        string     `gorm:"column:end_time" json:"endTime"`
	OpenId         string     `gorm:"column:open_id;index" json:"openId"`
	Qrcode         string     `gorm:"column:qrcode" json:"qrcode"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

// Inviter 对应 inviter_demo 表
type Inviter struct {
	ID             uint       `gorm:"primaryKey;column:id" json:"id"`
	InvitationId   string     `gorm:"column:invitation_id;index" json:"invitationId"`
	OpenID         string     `gorm:"column:open_id;index" json:"openId"`
	Nickname       string     `gorm:"column:nickname" json:"nickname"`
	Avatar         string     `gorm:"column:avatar" json:"avatar"`
	Company        string     `gorm:"column:company" json:"company"`
	Department     string     `gorm:"column:department" json:"department"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	// 动态字段将通过 JSON 存储或在具体业务中处理，MySQL 建议增加一个 Detail 字段
	Detail         string     `gorm:"column:detail;type:text" json:"detail"` 
}

// Lottery 对应 lottery_demo 表
type Lottery struct {
	ID             uint       `gorm:"primaryKey;column:id" json:"id"`
	InvitationId   string     `gorm:"column:invitation_id;index" json:"invitationId"`
	PrizeName      string     `gorm:"column:prize_name" json:"prizeName"`
	WinnerCount    int        `gorm:"column:winner_count" json:"winnerCount"`
	Winners        string     `gorm:"column:winners;type:text" json:"winners"` // 存储 JSON 数组
	Status         string     `gorm:"column:status" json:"status"`
	CreatorOpenId  string     `gorm:"column:creator_open_id" json:"creatorOpenId"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

// ApplyConfig 对应 apply_config_demo 表
type ApplyConfig struct {
	ID             uint       `gorm:"primaryKey;column:id" json:"id"`
	DefaultOpen    bool       `gorm:"column:default_open" json:"defaultOpen"`
	MaxParticipants string    `gorm:"column:max_participants" json:"maxParticipants"`
	FieldOptions   string     `gorm:"column:field_options;type:text" json:"fieldOptions"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
}
