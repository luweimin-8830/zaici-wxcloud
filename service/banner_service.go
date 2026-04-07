package service

import (
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
	"encoding/json"
	"fmt"
	"time"
)

type BannerService struct {
	bannerDao *dao.BannerDao
}

func NewBannerService() *BannerService {
	return &BannerService{
		bannerDao: dao.NewBannerDao(),
	}
}

func (s *BannerService) ListBanners() ([]model.Banner, error) {
	return s.bannerDao.GetAll()
}

func (s *BannerService) GetBannerDetail(id uint) (*model.Banner, error) {
	return s.bannerDao.GetByID(id)
}

func (s *BannerService) SaveBanner(bannerData map[string]interface{}) (string, error) {
	// 处理 detail 逻辑
	detailMap := make(map[string]interface{})
	if d, ok := bannerData["detail"].(map[string]interface{}); ok {
		detailMap = d
	}

	// 确保 type 和 url 存在
	if detailMap["type"] == nil {
		if t, ok := bannerData["type"].(string); ok {
			detailMap["type"] = t
		} else {
			detailMap["type"] = "image"
		}
	}
	if detailMap["url"] == nil {
		if u, ok := bannerData["url"].(string); ok {
			detailMap["url"] = u
		} else {
			detailMap["url"] = ""
		}
	}

	detailJson, _ := json.Marshal(detailMap)

	// 构造 Banner 模型
	banner := &model.Banner{
		Title:    getString(bannerData, "title"),
		ImageUrl: getString(bannerData, "imageUrl"),
		LinkUrl:  getString(bannerData, "linkUrl"),
		Type:     getString(bannerData, "type"),
		Url:      getString(bannerData, "url"),
		Detail:   string(detailJson),
		UpdatedAt: time.Now(),
	}

	if interval, ok := bannerData["interval"].(float64); ok {
		banner.Interval = int(interval)
	}

	// 检查是否是更新 (判断 _id)
	if idStr, ok := bannerData["_id"].(string); ok && idStr != "" {
		var id uint
		fmt.Sscanf(idStr, "%d", &id)
		banner.ID = id
		if err := s.bannerDao.Update(banner); err != nil {
			return "", err
		}
		return "更新成功", nil
	}

	// 新增
	banner.CreatedAt = time.Now()
	if err := s.bannerDao.Create(banner); err != nil {
		return "", err
	}
	return "新增成功", nil
}

func (s *BannerService) DeleteBanner(id uint) error {
	return s.bannerDao.Delete(id)
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
