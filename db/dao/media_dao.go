package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// PictureListDao 图片数据访问层
type PictureListDao struct{}

// Create 创建图片记录
func (d *PictureListDao) Create(picture *model.PictureList) error {
	return db.Get().Create(picture).Error
}

// GetByID 根据ID获取图片
func (d *PictureListDao) GetByID(id uint) (*model.PictureList, error) {
	var picture model.PictureList
	err := db.Get().First(&picture, id).Error
	if err != nil {
		return nil, err
	}
	return &picture, nil
}

// GetByPicHash 根据图片哈希获取图片
func (d *PictureListDao) GetByPicHash(picHash string) (*model.PictureList, error) {
	var picture model.PictureList
	err := db.Get().Where("pic_hash = ?", picHash).First(&picture).Error
	if err != nil {
		return nil, err
	}
	return &picture, nil
}

// Update 更新图片记录
func (d *PictureListDao) Update(picture *model.PictureList) error {
	return db.Get().Save(picture).Error
}

// UpdateSecCheckStatus 更新图片审核状态
func (d *PictureListDao) UpdateSecCheckStatus(id uint, status int) error {
	return db.Get().Model(&model.PictureList{}).Where("id = ?", id).Update("sec_check_status", status).Error
}

// Delete 删除图片记录
func (d *PictureListDao) Delete(id uint) error {
	return db.Get().Delete(&model.PictureList{}, id).Error
}

// List 获取图片列表
func (d *PictureListDao) List(limit, offset int) ([]model.PictureList, error) {
	var pictures []model.PictureList
	err := db.Get().Order("created_at DESC").Limit(limit).Offset(offset).Find(&pictures).Error
	return pictures, err
}

// Count 获取图片总数
func (d *PictureListDao) Count() (int64, error) {
	var count int64
	err := db.Get().Model(&model.PictureList{}).Count(&count).Error
	return count, err
}
