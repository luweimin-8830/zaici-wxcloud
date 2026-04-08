package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type ConfigDao struct{}

func NewConfigDao() *ConfigDao {
	return &ConfigDao{}
}

func (d *ConfigDao) GetConfig() (*model.Config, error) {
	var config model.Config
	// 假设只维护一条全局配置，取第一条
	err := db.Get().First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (d *ConfigDao) UpdateDistance(distance int) error {
	var config model.Config
	err := db.Get().First(&config).Error
	if err != nil {
		// 如果不存在则创建
		return db.Get().Create(&model.Config{
			Distance: distance,
		}).Error
	}
	return db.Get().Model(&config).Update("distance", distance).Error
}
