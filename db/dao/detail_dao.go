package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type DetailDao struct{}

func NewDetailDao() *DetailDao {
	return &DetailDao{}
}

func (d *DetailDao) GetActiveRecord(openId string) (*model.DetailRecord, error) {
	var record model.DetailRecord
	err := db.Get().Where("open_id = ? AND status = ?", openId, 1).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (d *DetailDao) CreateRecord(record *model.DetailRecord) error {
	return db.Get().Create(record).Error
}

func (d *DetailDao) UpdateRecord(id uint, data map[string]interface{}) error {
	return db.Get().Model(&model.DetailRecord{}).Where("id = ?", id).Updates(data).Error
}

func (d *DetailDao) UpdateSeat(openId string, seat string) error {
	return db.Get().Model(&model.DetailRecord{}).
		Where("open_id = ? AND status = ?", openId, 1).
		Update("seat", seat).Error
}

func (d *DetailDao) GetPictureByHash(picHash string) (*model.PictureList, error) {
	var pic model.PictureList
	err := db.Get().Where("pic_hash = ?", picHash).First(&pic).Error
	if err != nil {
		return nil, err
	}
	return &pic, nil
}
