package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type MatchDao struct{}

func NewMatchDao() *MatchDao {
	return &MatchDao{}
}

func (d *MatchDao) GetMatches(openId string, status int) ([]model.Match, error) {
	var list []model.Match
	err := db.Get().Where("status = ? AND (open_id_1 = ? OR open_id_2 = ?)", status, openId, openId).Find(&list).Error
	return list, err
}

func (d *MatchDao) GetLikeMatches(openId string) ([]model.Match, error) {
	var list []model.Match
	err := db.Get().Where("status = 1 AND open_id_2 = ?", openId).Order("created_at desc").Find(&list).Error
	return list, err
}

func (d *MatchDao) FindMatch(openId1, openId2 string) (*model.Match, error) {
	var match model.Match
	err := db.Get().Where("open_id_1 = ? AND open_id_2 = ?", openId1, openId2).First(&match).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (d *MatchDao) FindMatchReverse(openId1, openId2 string) (*model.Match, error) {
	var match model.Match
	err := db.Get().Where("open_id_1 = ? AND open_id_2 = ?", openId2, openId1).First(&match).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (d *MatchDao) Create(match *model.Match) error {
	return db.Get().Create(match).Error
}

func (d *MatchDao) Update(match *model.Match) error {
	return db.Get().Save(match).Error
}

func (d *MatchDao) DeleteByChannel(channel string) error {
	return db.Get().Where("channel = ?", channel).Delete(&model.Match{}).Error
}

func (d *MatchDao) DeleteByID(id uint) error {
	return db.Get().Delete(&model.Match{}, id).Error
}
