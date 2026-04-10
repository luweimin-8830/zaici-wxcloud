package dao

import (
	"fmt"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type OnlineDao struct{}

func NewOnlineDao() *OnlineDao {
	return &OnlineDao{}
}

func (d *OnlineDao) GetActiveRecord(openId string, now int64) (*model.OnlineRecord, error) {
	var record model.OnlineRecord
	err := db.Get().Where("open_id = ? AND due_time >= ?", openId, now).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (d *OnlineDao) GetShopOnlineUsers(shopId string, now int64, excludeOpenId string) ([]model.OnlineRecord, error) {
	var list []model.OnlineRecord
	query := db.Get().Where("shop_id = ? AND status = ? AND due_time >= ? AND flag = ?", shopId, "在线", now, 1)
	if excludeOpenId != "" {
		query = query.Where("open_id != ?", excludeOpenId)
	}
	err := query.Find(&list).Error
	fmt.Printf("GetShopOnlineUsers - shopId: %s, now: %d, count: %d, err: %v\n", shopId, now, len(list), err)
	return list, err
}

func (d *OnlineDao) CreateRecord(record *model.OnlineRecord) error {
	return db.Get().Create(record).Error
}

func (d *OnlineDao) UpdateDueTime(openId string, now int64) error {
	return db.Get().Model(&model.OnlineRecord{}).
		Where("open_id = ? AND due_time > ?", openId, now).
		Update("due_time", now).Error
}

func (d *OnlineDao) UpdateRecord(id uint, data map[string]interface{}) error {
	return db.Get().Model(&model.OnlineRecord{}).Where("id = ?", id).Updates(data).Error
}

func (d *OnlineDao) GetHistory(shopId string, now int64, excludeOpenId string) ([]model.OnlineRecord, error) {
	var list []model.OnlineRecord
	err := db.Get().Where("shop_id = ? AND due_time <= ? AND open_id != ?", shopId, now, excludeOpenId).
		Order("due_time desc").
		Limit(100).
		Find(&list).Error
	return list, err
}

func (d *OnlineDao) DeleteByOpenId(openId string) error {
	return db.Get().Where("open_id = ?", openId).Delete(&model.OnlineRecord{}).Error
}
