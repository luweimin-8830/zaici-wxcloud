package dao

import (
	"fmt"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type InvitationDao struct{}

func NewInvitationDao() *InvitationDao {
	return &InvitationDao{}
}

func (d *InvitationDao) GetByID(id string) (*model.Invitation, error) {
	var invitation model.Invitation
	var uintId uint
	fmt.Sscanf(id, "%d", &uintId)
	err := db.Get().First(&invitation, uintId).Error
	return &invitation, err
}

func (d *InvitationDao) Create(inv *model.Invitation) error {
	return db.Get().Create(inv).Error
}

func (d *InvitationDao) Update(id uint, data map[string]interface{}) error {
	return db.Get().Model(&model.Invitation{}).Where("id = ?", id).Updates(data).Error
}

func (d *InvitationDao) List(query map[string]interface{}, skip, limit int) ([]model.Invitation, int64, error) {
	var list []model.Invitation
	var total int64
	
	dbQuery := db.Get().Model(&model.Invitation{}).Where(query)
	dbQuery.Count(&total)
	err := dbQuery.Order("updated_at desc").Offset(skip).Limit(limit).Find(&list).Error
	
	return list, total, err
}

func (d *InvitationDao) Delete(id uint) error {
	return db.Get().Delete(&model.Invitation{}, id).Error
}

// --- Inviter (报名者) ---

func (d *InvitationDao) CountJoiners(invitationId string, openId string) (int64, error) {
	var count int64
	err := db.Get().Model(&model.Inviter{}).Where("invitation_id = ? AND open_id = ?", invitationId, openId).Count(&count).Error
	return count, err
}

func (d *InvitationDao) CreateJoiner(inviter *model.Inviter) error {
	return db.Get().Create(inviter).Error
}

func (d *InvitationDao) GetJoiners(invitationId string, limit int) ([]model.Inviter, error) {
	var list []model.Inviter
	err := db.Get().Where("invitation_id = ?", invitationId).Order("created_at desc").Limit(limit).Find(&list).Error
	return list, err
}

// --- Lottery (抽奖) ---

func (d *InvitationDao) SaveLottery(lottery *model.Lottery) error {
	if lottery.ID != 0 {
		return db.Get().Save(lottery).Error
	}
	return db.Get().Create(lottery).Error
}

func (d *InvitationDao) ListLotteries(query map[string]interface{}, skip, limit int) ([]model.Lottery, int64, error) {
	var list []model.Lottery
	var total int64
	
	dbQuery := db.Get().Model(&model.Lottery{}).Where(query)
	dbQuery.Count(&total)
	err := dbQuery.Order("created_at desc").Offset(skip).Limit(limit).Find(&list).Error
	
	return list, total, err
}

func (d *InvitationDao) GetLotteryByID(id uint) (*model.Lottery, error) {
	var lottery model.Lottery
	err := db.Get().First(&lottery, id).Error
	return &lottery, err
}

func (d *InvitationDao) DeleteLottery(id uint) error {
	return db.Get().Delete(&model.Lottery{}, id).Error
}

// --- Config (报名配置) ---

func (d *InvitationDao) GetApplyConfig() (*model.ApplyConfig, error) {
	var config model.ApplyConfig
	err := db.Get().First(&config).Error
	return &config, err
}

func (d *InvitationDao) SaveApplyConfig(config *model.ApplyConfig) error {
	return db.Get().Save(config).Error
}

func (d *InvitationDao) UpdateConfigFields(id uint, fieldOptions string) error {
	return db.Get().Model(&model.ApplyConfig{}).Where("id = ?", id).Update("field_options", fieldOptions).Error
}
