package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"

	"gorm.io/gorm"
)

type ChatDao struct{}

func NewChatDao() *ChatDao {
	return &ChatDao{}
}

func (d *ChatDao) GetHistory(channelId string, skip int, limit int) ([]model.ChatHistory, error) {
	var histories []model.ChatHistory
	err := db.Get().Where("channel_id = ?", channelId).
		Order("timestamp desc").
		Offset(skip).
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

func (d *ChatDao) CreateMessage(msg *model.ChatHistory) error {
	return db.Get().Create(msg).Error
}

func (d *ChatDao) UpdateMessageStateByReceiver(channelId string, receiverOpenId string, state int) error {
	// MySQL JSON update: UPDATE table SET col = JSON_SET(col, '$.state', val)
	return db.Get().Model(&model.ChatHistory{}).
		Where("channel_id = ? AND receiver_open_id = ?", channelId, receiverOpenId).
		UpdateColumn("message_content", gorm.Expr("JSON_SET(message_content, '$.state', ?)", state)).Error
}

func (d *ChatDao) UpdateMessageStateByID(channelId string, msgId string, state int) error {
	// JSON path update for specific message content ID
	return db.Get().Model(&model.ChatHistory{}).
		Where("channel_id = ? AND JSON_EXTRACT(message_content, '$.id') = ?", channelId, msgId).
		UpdateColumn("message_content", gorm.Expr("JSON_SET(message_content, '$.state', ?)", state)).Error
}

func (d *ChatDao) SaveBlock(openId string, blockId string) error {
	return db.Get().Create(&model.BlockList{
		OpenID:  openId,
		BlockID: blockId,
	}).Error
}

func (d *ChatDao) DeleteInfoMonitor(openId string) error {
	return db.Get().Where("open_id = ?", openId).Delete(&model.InfoMonitor{}).Error
}
