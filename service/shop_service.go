package service

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
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
	// 获取所有门店（使用较大的 limit 值来获取所有门店）
	var shops []model.Shop
	var total int64
	shops, total, err := s.shopDao.List(nil, 0, 1000)
	if err != nil {
		fmt.Printf("GetNearList - shopDao.List error: %v\n", err)
		return nil, err
	}
	fmt.Printf("GetNearList - shops count: %d, total: %d\n", len(shops), total)
	_ = total

	// 构建返回结果，计算距离并过滤
	result := make([]map[string]interface{}, 0, len(shops))
	for _, shop := range shops {
		// 解析 location 字段（GeoJSON 格式）
		shopLongitude, shopLatitude := parseLocation(shop.Location)
		
		// 计算距离（单位：公里）
		dist := haversineDistance(longitude, latitude, shopLongitude, shopLatitude)
		
		// 只返回在指定距离范围内的门店
		if dist <= distance/1000 { // distance 参数是米，转换为公里比较
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
				"distance":  dist, // 距离（公里）
			}
			result = append(result, shopMap)
		}
	}
	
	// 按距离排序
	sort.Slice(result, func(i, j int) bool {
		return result[i]["distance"].(float64) < result[j]["distance"].(float64)
	})
	
	return result, nil
}

// parseLocation 解析 GeoJSON 格式的 location 字符串
func parseLocation(location string) (longitude, latitude float64) {
	if location == "" {
		return 0, 0
	}
	
	// 尝试解析 {"type":"Point","coordinates":[longitude,latitude]}
	var geo struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	}
	
	if err := json.Unmarshal([]byte(location), &geo); err == nil && len(geo.Coordinates) >= 2 {
		return geo.Coordinates[0], geo.Coordinates[1] // longitude, latitude
	}
	
	return 0, 0
}

// haversineDistance 计算两点之间的球面距离（单位：公里）
func haversineDistance(lon1, lat1, lon2, lat2 float64) float64 {
	const R = 6371 // 地球半径（公里）
	
	// 转换为弧度
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180
	
	// Haversine 公式
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return R * c
}
