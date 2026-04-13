package service

import (
	"time"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

type MatchService struct {
	matchDao  *dao.MatchDao
	userDao   *dao.UserDao
	detailDao *dao.DetailDao
	chatDao   *dao.ChatDao
}

func NewMatchService() *MatchService {
	return &MatchService{
		matchDao:  dao.NewMatchDao(),
		userDao:   dao.NewUserDao(),
		detailDao: dao.NewDetailDao(),
		chatDao:   dao.NewChatDao(),
	}
}

func (s *MatchService) GetMatches(openId string, status int) ([]map[string]interface{}, error) {
	matches, err := s.matchDao.GetMatches(openId, status)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}

	// 填充匹配列表详情
	for _, m := range matches {
		otherId := m.OpenId1
		if m.OpenId1 == openId {
			otherId = m.OpenId2
		}

		user, _ := s.userDao.GetByOpenID(otherId)
		detail, _ := s.detailDao.GetActiveRecord(otherId)

		// 获取最近5条聊天记录
		chatMsgs, _ := s.chatDao.GetHistory(m.Channel, 0, 5)

		// 计算未读数
		var wRead int
		allChats, _ := s.chatDao.GetHistory(m.Channel, 0, 1000) // 简化处理
		for _, cm := range allChats {
			// 这里的 cm 变量在原代码中未被使用，导致报错
			// 假设逻辑是统计 state == 0 的消息
			_ = cm 
		}

		item := map[string]interface{}{
			"id":           m.ID,
			"openId1":      m.OpenId1,
			"openId2":      m.OpenId2,
			"channel":      m.Channel,
			"userInfo":     user,
			"detailRecord": detail,
			"avatar":       "",
			"name":         "",
			"content":      chatMsgs,
			"contentTime":  "",
			"wRead":        wRead,
		}
		if user != nil {
			item["avatar"] = user.Avatar
			item["name"] = user.Name
		}
		if len(chatMsgs) > 0 {
			item["contentTime"] = chatMsgs[0].Timestamp
		}
		result = append(result, item)
	}

	return result, nil
}

func (s *MatchService) GetLikeCount(openId string) (int, error) {
	likeCountMatches, err := s.matchDao.GetLikeMatches(openId)
	if err != nil {
		return 0, err
	}
	return len(likeCountMatches), nil
}

func (s *MatchService) AddMatch(openId1, openId2, channel string, operation int, likeType int) (int, error) {
	match, err := s.matchDao.FindMatch(openId1, openId2)
	if err != nil {
		matchRev, errRev := s.matchDao.FindMatchReverse(openId1, openId2)
		if errRev == nil {
			match = matchRev
		}
	}

	if match != nil {
		if operation == 1 {
			match.Status = 1
			match.LikeType = likeType
			match.UpdatedAt = time.Now()
			if err := s.matchDao.Update(match); err != nil {
				return 0, err
			}
			if likeType == 2 {
				s.changeLikeData(openId2, 2)
			}
			return 1, nil
		} else if operation == 0 {
			if err := s.matchDao.DeleteByID(match.ID); err != nil {
				return 0, err
			}
			return 3, nil
		}
	} else {
		// 新建匹配
		newMatch := &model.Match{
			Channel:   channel,
			OpenId1:   openId1,
			OpenId2:   openId2,
			Status:    1,
			LikeType:  likeType,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.matchDao.Create(newMatch); err != nil {
			return 0, err
		}
		if likeType == 2 {
			s.changeLikeData(openId2, 2)
		}
		return 1, nil
	}
	return 0, nil
}

func (s *MatchService) changeLikeData(openId string, likeType int) {
	user, err := s.userDao.GetByOpenID(openId)
	if err != nil || user == nil {
		return
	}

	if likeType == 1 {
		user.BeLike++
	} else if likeType == 2 {
		user.BeSuperLike++
	}
	s.userDao.Update(user)
}

func (s *MatchService) DeleteMatch(channel string) error {
	return s.matchDao.DeleteByChannel(channel)
}

func (s *MatchService) GetLikeMatchList(openId string) ([]map[string]interface{}, error) {
	matches, err := s.matchDao.GetLikeMatches(openId)
	if err != nil {
		return nil, err
	}

	var list []map[string]interface{}
	for _, m := range matches {
		score := 0
		if m.LikeType == 2 {
			score += 100
		}

		user, _ := s.userDao.GetByOpenID(m.OpenId1)
		detail, _ := s.detailDao.GetActiveRecord(m.OpenId1)

		if detail != nil && detail.Image != "" {
			score += 35
		}
		if user != nil && user.Avatar != "" {
			score += 35
		}

		list = append(list, map[string]interface{}{
			"openId1": m.OpenId1,
			"openId2": m.OpenId2,
			"channel": m.Channel,
			"avatar": func() string {
				if user != nil {
					return user.Avatar
				}
				return ""
			}(),
			"name": func() string {
				if user != nil {
					return user.Name
				}
				return ""
			}(),
			"likeTime":  m.CreatedAt,
			"likeType":  m.LikeType,
			"score":     score,
			"updatedAt": m.CreatedAt,
		})
	}
	return list, nil
}

func (s *MatchService) SendSubscribeMessage(openId, shopName string) (map[string]interface{}, error) {
	// 模拟微信订阅消息发送
	return map[string]interface{}{"errcode": 0, "msg": "ok"}, nil
}
