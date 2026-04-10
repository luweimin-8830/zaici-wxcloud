package service

import (
	"encoding/json"
	"fmt"
	"time"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"
)

type OnlineService struct {
	onlineDao  *dao.OnlineDao
	userDao    *dao.UserDao
	detailDao  *dao.DetailDao
	matchDao   *dao.MatchDao
}

func NewOnlineService() *OnlineService {
	return &OnlineService{
		onlineDao:  dao.NewOnlineDao(),
		userDao:    dao.NewUserDao(),
		detailDao:  dao.NewDetailDao(),
		matchDao:   dao.NewMatchDao(),
	}
}

func (s *OnlineService) GetStatus(openId string) (map[string]interface{}, error) {
	now := time.Now().UnixMilli()
	record, err := s.onlineDao.GetActiveRecord(openId, now)
	if err != nil {
		return map[string]interface{}{"status": "未出门"}, nil
	}
	return map[string]interface{}{
		"position": record,
		"status":   "已出门",
	}, nil
}

func (s *OnlineService) GetNearList(openId string, longitude, latitude float64) ([]map[string]interface{}, error) {
	now := time.Now().UnixMilli()
	
	// 注意：MySQL 并不原生支持像 TCB 那样的 .geoNear() 聚合
	// 在生产环境下，通常使用 ST_Distance_Sphere 或第三方地理索引 (如 Redis Geo)
	// 这里实现一个基于经纬度过滤的简化逻辑，实际部署建议使用 MySQL 空间索引
	records, err := s.onlineDao.GetShopOnlineUsers("", now, openId) // 简化：暂不传 shopId
	if err != nil {
		return nil, err
	}

	var userList []map[string]interface{}
	for _, rec := range records {
		// 1. 检查匹配状态
		match, _ := s.matchDao.FindMatch(openId, rec.OpenID)
		if match == nil {
			matchRev, _ := s.matchDao.FindMatchReverse(openId, rec.OpenID)
			match = matchRev
		}

		state, likeType, channel := 0, 0, ""
		if match != nil {
			if match.Status == 2 {
				state = 1
			} else if match.Status == 1 {
				state = 0
			}
			if match.OpenId1 == openId {
				likeType = match.LikeType
				channel = match.Channel
			} else {
				channel = match.Channel
			}
		}

		// 2. 补齐用户信息
		user, _ := s.userDao.GetByOpenID(rec.OpenID)
		detail, _ := s.detailDao.GetActiveRecord(rec.OpenID)
		
		score := 0
		if likeType == 2 { score += 100 }
		if user != nil && user.Avatar != "" { score += 35 }
		if detail != nil && detail.Image != "" { score += 35 }

		item := map[string]interface{}{
			"openId":       rec.OpenID,
			"state":        state,
			"likeType":     likeType,
			"channel":      channel,
			"score":        score,
			"userInfo":     user,
			"detailRecord": detail,
			"avatar":       func() string { if user != nil { return user.Avatar }; return "" }(),
			"name":         func() string { if user != nil { return user.Name }; return "" }(),
		}
		userList = append(userList, item)
	}

	return userList, nil
}

func (s *OnlineService) SaveOnline(openId string, data map[string]interface{}) (*model.OnlineRecord, error) {
	// 计算过期时间：次日 05:00 UTC
	now := time.Now().UTC()
	due := time.Date(now.Year(), now.Month(), now.Day()+1, 5, 0, 0, 0, time.UTC)
	
	// 1. 将之前的在线记录设为结束
	s.onlineDao.UpdateDueTime(openId, now.UnixMilli())

	// 处理 shopId，支持字符串和数字类型
	shopId := utils.GetString(data, "shopId")
	if shopId == "" {
		// 尝试从其他类型转换
		if v, ok := data["shopId"].(float64); ok {
			shopId = fmt.Sprintf("%.0f", v)
		} else if v, ok := data["shopId"].(int); ok {
			shopId = fmt.Sprintf("%d", v)
		}
	}

	fmt.Printf("SaveOnline - shopId: %s\n", shopId)

	// 2. 构造新记录
	rec := &model.OnlineRecord{
		OpenID:    openId,
		ShopID:    shopId,
		ShopName:  utils.GetString(data, "shopName"),
		Status:    "在线",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DueTime:   due.UnixMilli(),
		Flag:      utils.ParseInt(data, "flag"),
	}

	// 填充用户信息
	user, _ := s.userDao.GetByOpenID(openId)
	if user != nil {
		rec.Name = user.Name
		rec.Avatar = user.Avatar
	}

	// 处理 Location (GeoJSON)
	if loc, ok := data["location"].(map[string]interface{}); ok {
		b, _ := json.Marshal(loc)
		rec.Location = string(b)
	}

	if err := s.onlineDao.CreateRecord(rec); err != nil {
		return nil, err
	}
	return rec, nil
}

func (s *OnlineService) UpdateOnline(id uint, dueTime int64) error {
	return s.onlineDao.UpdateDueTime(fmt.Sprintf("%d", id), dueTime) // 简化处理
}

// ... existing code ...
func (s *OnlineService) GetHistory(shopId string, openId string) ([]model.OnlineRecord, error) {
	now := time.Now().UnixMilli()
	return s.onlineDao.GetHistory(shopId, now, openId)
}

func (s *OnlineService) GetShopOnline(shopId string, openId string) ([]map[string]interface{}, error) {
	now := time.Now().UnixMilli()

	// 获取门店在线用户
	records, err := s.onlineDao.GetShopOnlineUsers(shopId, now, openId)
	if err != nil {
		return nil, err
	}

	var userList []map[string]interface{}
	for _, rec := range records {
		// 检查匹配状态
		match, _ := s.matchDao.FindMatch(openId, rec.OpenID)
		if match == nil {
			matchRev, _ := s.matchDao.FindMatchReverse(openId, rec.OpenID)
			match = matchRev
		}

		state, likeType, channel := 0, 0, ""
		if match != nil {
			if match.Status == 2 {
				state = 1
			} else if match.Status == 1 {
				state = 0
			}
			if match.OpenId1 == openId {
				likeType = match.LikeType
				channel = match.Channel
			} else {
				channel = match.Channel
			}
		}

		// 补齐用户信息
		user, _ := s.userDao.GetByOpenID(rec.OpenID)
		detail, _ := s.detailDao.GetActiveRecord(rec.OpenID)

		score := 0
		if likeType == 2 { score += 100 }
		if user != nil && user.Avatar != "" { score += 35 }
		if detail != nil && detail.Image != "" { score += 35 }

		item := map[string]interface{}{
			"openId":       rec.OpenID,
			"state":        state,
			"likeType":     likeType,
			"channel":      channel,
			"score":        score,
			"userInfo":     user,
			"detailRecord": detail,
			"avatar":       func() string { if user != nil { return user.Avatar }; return "" }(),
			"name":         func() string { if user != nil { return user.Name }; return "" }(),
		}
		userList = append(userList, item)
	}

	return userList, nil
}
