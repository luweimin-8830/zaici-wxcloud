package service

import (
	"fmt"
	"time"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"
)

type ShopService struct {
	shopDao *dao.ShopDao
}

func NewShopService() *ShopService {
	return &ShopService{
		shopDao: dao.NewShopDao(),
	}
}

func (s *ShopService) GetDetail(id string) (*model.Shop, error) {
	return s.shopDao.GetByID(id)
}

func (s *ShopService) Save(data map[string]interface{}) (*model.Shop, error) {
	// 兼容前端字段名：shopName 或 shopname
	shopName := utils.GetString(data, "shopName")
	if shopName == "" {
		shopName = utils.GetString(data, "shopname")
	}
	
	shop := &model.Shop{
		ShopName:  shopName,
		Location:  utils.GetString(data, "location"),
		Address:   utils.GetString(data, "address"),
		Phone:     utils.GetString(data, "phone"),
		Image:     utils.GetString(data, "image"),
		StartTime: utils.GetString(data, "startTime"),
		EndTime:   utils.GetString(data, "endTime"),
		Tag1:      utils.GetString(data, "tag1"),
		Tag2:      utils.GetString(data, "tag2"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.shopDao.Create(shop); err != nil {
		return nil, err
	}
	return shop, nil
}

func (s *ShopService) Update(id uint, data map[string]interface{}) error {
	// 过滤不需要更新的字段
	updates := make(map[string]interface{})
	for k, v := range data {
		if k == "createdAt" || k == "startTime" || k == "endTime" {
			continue
		}
		updates[k] = v
	}
	updates["updated_at"] = time.Now()

	return s.shopDao.Update(id, updates)
}

func (s *ShopService) Delete(id uint) error {
	return s.shopDao.Delete(id)
}

func (s *ShopService) AdminList(page, limit int, keyword string) ([]model.Shop, int64, error) {
	skip := (page - 1) * limit
	fmt.Printf("ShopService.AdminList - page: %d, limit: %d, skip: %d, keyword: %s\n", page, limit, skip, keyword)
	shops, total, err := s.shopDao.Search(keyword, skip, limit)
	if err != nil {
		fmt.Printf("ShopDao.Search error: %v\n", err)
	}
	fmt.Printf("ShopDao.Search result - shops: %d, total: %d\n", len(shops), total)
	return shops, total, err
}

func (s *ShopService) GetNearList(longitude, latitude, distance float64) ([]map[string]interface{}, error) {
	// MySQL 简化实现：实际应使用 ST_Distance_Sphere
	// 这里暂且返回所有或简单过滤，建议后续优化为空间查询
	var shops []model.Shop
	var total int64
	shops, total, err := s.shopDao.List(nil, 0, int(total))
	if err != nil {
		return nil, err
	}

	// 构建返回结果，添加距离字段
	result := make([]map[string]interface{}, 0, len(shops))
	for _, shop := range shops {
		shopMap := map[string]interface{}{
			"id":        shop.ID,
			"shopName":  shop.ShopName,
			"address":   shop.Address,
			"phone":     shop.Phone,
			"image":     shop.Image,
			"tag1":      shop.Tag1,
			"tag2":      shop.Tag2,
			"startTime": shop.StartTime,
			"endTime":   shop.EndTime,
			"distance":  0, // 简化处理，实际应根据经纬度计算
		}
		result = append(result, shopMap)
	}
	_ = total // 避免未使用变量警告
	return result, nil
}
