package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"

	"gorm.io/gorm"
)

type UserDao struct{}

func NewUserDao() *UserDao {
	return &UserDao{}
}

func (d *UserDao) GetByOpenID(openId string) (*model.User, error) {
	var user model.User
	err := db.Get().Where("open_id = ?", openId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDao) Create(user *model.User) error {
	return db.Get().Create(user).Error
}

func (d *UserDao) Update(user *model.User) error {
	return db.Get().Save(user).Error
}

func (d *UserDao) UpdateField(openId string, key string, value interface{}) error {
	return db.Get().Model(&model.User{}).Where("open_id = ?", openId).Update(key, value).Error
}

func (d *UserDao) GetSuperLike(openId string) (int, error) {
	var user model.User
	err := db.Get().Select("super_like").Where("open_id = ?", openId).First(&user).Error
	return user.SuperLike, err
}

func (d *UserDao) SetSuperLike(openId string, count int) error {
	return db.Get().Model(&model.User{}).Where("open_id = ?", openId).Update("super_like", count).Error
}

func (d *UserDao) CheckAdmin(openId string) (bool, error) {
	var admin model.Admin
	err := db.Get().Where("open_id = ?", openId).First(&admin).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *UserDao) AddAdmin(openId string) error {
	return db.Get().Create(&model.Admin{OpenID: openId}).Error
}
