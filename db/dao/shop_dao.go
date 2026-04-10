package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type ShopDao struct{}

func NewShopDao() *ShopDao {
	return &ShopDao{}
}

func (dao *ShopDao) GetByID(id string) (*model.Shop, error) {
	var shop model.Shop
	if err := db.Get().Where("id = ?", id).First(&shop).Error; err != nil {
		return nil, err
	}
	return &shop, nil
}

func (dao *ShopDao) Create(shop *model.Shop) error {
	return db.Get().Create(shop).Error
}

func (dao *ShopDao) Update(id uint, updates map[string]interface{}) error {
	return db.Get().Model(&model.Shop{}).Where("id = ?", id).Updates(updates).Error
}

func (dao *ShopDao) Delete(id uint) error {
	return db.Get().Delete(&model.Shop{}, id).Error
}

func (dao *ShopDao) List(query map[string]interface{}, skip, limit int) ([]model.Shop, int64, error) {
	var shops []model.Shop
	var total int64

	dbQuery := db.Get().Model(&model.Shop{})
	if query != nil {
		dbQuery = dbQuery.Where(query)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.Offset(skip).Limit(limit).Order("created_at DESC").Find(&shops).Error; err != nil {
		return nil, 0, err
	}

	return shops, total, nil
}

func (dao *ShopDao) Search(keyword string, skip, limit int) ([]model.Shop, int64, error) {
	var shops []model.Shop
	var total int64

	dbQuery := db.Get().Model(&model.Shop{})
	
	// 如果 keyword 不为空，添加模糊查询条件
	if keyword != "" {
		dbQuery = dbQuery.Where("shopname LIKE ?", "%"+keyword+"%")
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.Offset(skip).Limit(limit).Order("created_at DESC").Find(&shops).Error; err != nil {
		return nil, 0, err
	}

	return shops, total, nil
}
