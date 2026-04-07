package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type BannerDao struct{}

func NewBannerDao() *BannerDao {
	return &BannerDao{}
}

func (d *BannerDao) GetAll() ([]model.Banner, error) {
	var banners []model.Banner
	err := db.Get().Find(&banners).Error
	return banners, err
}

func (d *BannerDao) GetByID(id uint) (*model.Banner, error) {
	var banner model.Banner
	err := db.Get().First(&banner, id).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

func (d *BannerDao) Create(banner *model.Banner) error {
	return db.Get().Create(banner).Error
}

func (d *BannerDao) Update(banner *model.Banner) error {
	return db.Get().Save(banner).Error
}

func (d *BannerDao) Delete(id uint) error {
	return db.Get().Delete(&model.Banner{}, id).Error
}
