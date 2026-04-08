package service

import (
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

type ConfigService struct {
	configDao *dao.ConfigDao
}

func NewConfigService() *ConfigService {
	return &ConfigService{
		configDao: dao.NewConfigDao(),
	}
}

func (s *ConfigService) GetDistance() (*model.Config, error) {
	config, err := s.configDao.GetConfig()
	if err != nil {
		// 如果没找到配置，则创建默认配置 (distance = 1)
		defaultConfig := &model.Config{
			Distance: 1,
		}
		if errCreate := s.configDao.UpdateDistance(1); errCreate != nil {
			return nil, errCreate
		}
		return defaultConfig, nil
	}
	return config, nil
}

func (s *ConfigService) SaveDistance(distance int) error {
	return s.configDao.UpdateDistance(distance)
}
