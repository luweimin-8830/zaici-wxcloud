package service

import (
	"time"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService() *UserService {
	return &UserService{
		userDao: dao.NewUserDao(),
	}
}

func (s *UserService) GetOrCreateUser(openId string, body map[string]interface{}) (*model.User, error) {
	user, err := s.userDao.GetByOpenID(openId)
	if err != nil {
		// 新用户
		nickname, _ := body["nickname"].(string)
		userName := nickname
		if userName == "" {
			userName = "用户" + openId[len(openId)-4:]
		}

		avatar, _ := body["avatar"].(string)
		if avatar == "" {
			avatar = "https://cloud1-1gth9cum37c9015c-1380861431.tcloudbaseapp.com/logo.png?sign=44d3bfbeb3cb05b7ff6fce3460a8bcdf&t=1765181560"
		}

		newUser := &model.User{
			OpenID:      openId,
			Name:        userName,
			Avatar:      avatar,
			SuperLike:   0,
			BeLike:      0,
			BeSuperLike: 0,
			Created:     time.Now(),
			LastLogin:   time.Now(),
		}

		// 动态填充 restBody 字段
		if company, ok := body["company"].(string); ok {
			newUser.Company = company
		}
		if dept, ok := body["department"].(string); ok {
			newUser.Department = dept
		}
		if phone, ok := body["phone"].(string); ok {
			newUser.Phone = phone
		}

		if err := s.userDao.Create(newUser); err != nil {
			return nil, err
		}
		return newUser, nil
	}

	// 老用户更新资料
	updateData := make(map[string]interface{})
	updateData["last_login"] = time.Now()

	if nickname, ok := body["nickname"].(string); ok && nickname != "" {
		user.Name = nickname
	}
	if avatar, ok := body["avatar"].(string); ok && avatar != "" {
		user.Avatar = avatar
	}
	if company, ok := body["company"].(string); ok {
		user.Company = company
	}
	if dept, ok := body["department"].(string); ok {
		user.Department = dept
	}
	if phone, ok := body["phone"].(string); ok {
		user.Phone = phone
	}

	user.UpdatedAt = time.Now()
	if err := s.userDao.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) SaveUserInfo(openId string, key string, data interface{}) error {
	return s.userDao.UpdateField(openId, key, data)
}

func (s *UserService) UseSuperLike(openId string) (string, error) {
	count, err := s.userDao.GetSuperLike(openId)
	if err != nil {
		return "数据库错误", err
	}
	if count <= 0 {
		return "superLike次数已用完", nil
	}

	err = s.userDao.SetSuperLike(openId, count-1)
	if err != nil {
		return "更新失败", err
	}
	return "super Like次数-1", nil
}

type SuperLikeRequest struct {
	OpenId string `json:"openId"`
	Total  int    `json:"total"`
}

func (s *UserService) AddSuperLike(userList []SuperLikeRequest) error {
	for _, u := range userList {
		count, err := s.userDao.GetSuperLike(u.OpenId)
		if err != nil {
			continue
		}
		s.userDao.SetSuperLike(u.OpenId, count+u.Total)
	}
	return nil
}

func (s *UserService) ManageAdmin(openId string) (string, error) {
	isAdmin, err := s.userDao.CheckAdmin(openId)
	if err != nil {
		return "", err
	}
	if isAdmin {
		return "已是管理员", nil
	}
	if err := s.userDao.AddAdmin(openId); err != nil {
		return "", err
	}
	return "添加管理员成功", nil
}
