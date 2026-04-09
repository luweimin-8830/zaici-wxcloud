package service

import (
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
	"time"
)

type DetailService struct {
	detailDao *dao.DetailDao
	userDao   *dao.UserDao
}

func NewDetailService() *DetailService {
	return &DetailService{
		detailDao: dao.NewDetailDao(),
		userDao:   dao.NewUserDao(),
	}
}

func (s *DetailService) GetMyDetail(openId string) (*model.DetailRecord, error) {
	record, err := s.detailDao.GetActiveRecord(openId)
	if err != nil {
		return nil, err
	}

	user, err := s.userDao.GetByOpenID(openId)
	if err == nil && user != nil {
		// 模拟 JS 中的: info.avatar = users.data[0].avatar; info.name = users.data[0].name;
		// 因为 model.DetailRecord 没有 Name 和 Avatar 字段，可以通过 map 返回或者扩展模型
		// 这里我们通过 model 扩展或由 handler 处理，Service 层返回 record 及其关联用户
	}

	return record, nil
}

func (s *DetailService) SaveDetail(openId string, data map[string]interface{}) (*model.DetailRecord, error) {
	record, err := s.detailDao.GetActiveRecord(openId)
	
	if err != nil {
		// 新建记录
		newRecord := &model.DetailRecord{
			OpenID:      openId,
			Plan:        CommonGetString(data, "plan"),
			Mood:        CommonGetString(data, "mood"),
			Style:       CommonGetString(data, "style"),
			Description: CommonGetString(data, "description"),
			Seat:        CommonGetString(data, "seat"),
			Image:       CommonGetString(data, "image"),
			Status:      1,
			CreatedAt:   time.Now(),
		}
		if err := s.detailDao.CreateRecord(newRecord); err != nil {
			return nil, err
		}
		return newRecord, nil
	}

	// 更新记录
	updateMap := make(map[string]interface{})
	updateMap["plan"] = CommonGetString(data, "plan")
	updateMap["mood"] = CommonGetString(data, "mood")
	updateMap["style"] = CommonGetString(data, "style")
	updateMap["description"] = CommonGetString(data, "description")
	updateMap["seat"] = CommonGetString(data, "seat")
	updateMap["image"] = CommonGetString(data, "image")
	updateMap["updated_at"] = time.Now()

	if err := s.detailDao.UpdateRecord(record.ID, updateMap); err != nil {
		return nil, err
	}
	
	// 刷新记录对象
	record.Plan = CommonGetString(data, "plan")
	record.Mood = CommonGetString(data, "mood")
	record.Style = CommonGetString(data, "style")
	record.Description = CommonGetString(data, "description")
	record.Seat = CommonGetString(data, "seat")
	record.Image = CommonGetString(data, "image")
	
	return record, nil
}

func (s *DetailService) SaveSeat(openId string, seat string) error {
	return s.detailDao.UpdateSeat(openId, seat)
}
