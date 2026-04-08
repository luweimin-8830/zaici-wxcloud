package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"
)

type InvitationService struct {
	invitationDao *dao.InvitationDao
	userDao       *dao.UserDao
}

func NewInvitationService() *InvitationService {
	return &InvitationService{
		invitationDao: dao.NewInvitationDao(),
		userDao:       dao.NewUserDao(),
	}
}

// 计算活动动态状态 (复现 JS 的 computeActivityStatus)
func (s *InvitationService) computeStatus(inv *model.Invitation) string {
	if inv.Status == "已停止" {
		return "已截止报名"
	}

	now := time.Now()
	startTime, errS := time.Parse("2006-01-02 15:04", inv.StartTime)
	endTime, errE := time.Parse("2006-01-02 15:04", inv.EndTime)

	if errS == nil && now.Before(startTime) {
		return startTime.Format("01月02日 15:04") + " 开始"
	}
	if errE == nil && now.After(endTime) {
		return "已结束"
	}

	if inv.Status != "" {
		return inv.Status
	}
	return "进行中"
}

func (s *InvitationService) SaveInvitation(openId string, data map[string]interface{}) (string, string, error) {
	title := utils.GetString(data, "title")
	imageUrl := utils.GetString(data, "imageUrl")
	activity := utils.GetString(data, "activity")

	if title == "" || imageUrl == "" || activity == "" {
		return "", "", fmt.Errorf("参数缺失")
	}

	inv := &model.Invitation{
		Title:          title,
		ImageUrl:       imageUrl,
		Activity:       activity,
		Status:         utils.GetString(data, "status"),
		PosterUrl:      utils.GetString(data, "posterUrl"),
		RequiredFields: utils.GetString(data, "requiredFields"),
		StartTime:      utils.GetString(data, "startTime"),
		EndTime:        utils.GetString(data, "endTime"),
		OpenId:         openId,
		UpdatedAt:      time.Now(),
	}
	if inv.Status == "" {
		inv.Status = "进行中"
	}

	var invitationId string
	if idStr, ok := data["id"].(string); ok && idStr != "" {
		var uintId uint
		fmt.Sscanf(idStr, "%d", &uintId)
		if err := s.invitationDao.Update(uintId, map[string]interface{}{
			"title": inv.Title, "image_url": inv.ImageUrl, "activity": inv.Activity,
			"status": inv.Status, "poster_url": inv.PosterUrl, "required_fields": inv.RequiredFields,
			"start_time": inv.StartTime, "end_time": inv.EndTime, "updated_at": time.Now(),
		}); err != nil {
			return "", "", err
		}
		invitationId = idStr
	} else {
		inv.CreatedAt = time.Now()
		if err := s.invitationDao.Create(inv); err != nil {
			return "", "", err
		}
		invitationId = fmt.Sprintf("%d", inv.ID)
	}

	// 生成小程序码
	qrCode, err := s.generateQRCode(invitationId)
	if err != nil {
		fmt.Println("QRCode Error:", err)
		return invitationId, "", nil
	}

	// 更新 qrcode 到数据库
	s.invitationDao.Update(utils.ParseUint(invitationId), map[string]interface{}{
		"qrcode":     qrCode,
		"updated_at": time.Now(),
	})

	return invitationId, qrCode, nil
}

func (s *InvitationService) generateQRCode(id string) (string, error) {
	apiUrl := "http://api.weixin.qq.com/wxa/getwxacode"
	body := map[string]interface{}{
		"path":        fmt.Sprintf("pages/activityEntry/activityEntry?id=%s", id),
		"width":       430,
		"env_version": "release",
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return fmt.Sprintf("cloud://qrcode/invitation_%s.png", id), nil
}

func (s *InvitationService) JoinActivity(invitationId string, openId string, restBody map[string]interface{}) error {
	inv, err := s.invitationDao.GetByID(invitationId)
	if err != nil {
		return fmt.Errorf("活动不存在")
	}
	if inv.Status == "已停止" {
		return fmt.Errorf("该活动报名已截止")
	}

	count, _ := s.invitationDao.CountJoiners(invitationId, openId)
	if count > 0 {
		return fmt.Errorf("您已报名该活动，无需重复报名")
	}

	inviter := &model.Inviter{
		InvitationId: invitationId,
		OpenID:       openId,
		Nickname:     utils.GetString(restBody, "nickname"),
		Avatar:       utils.GetString(restBody, "avatar"),
		Company:      utils.GetString(restBody, "company"),
		Department:   utils.GetString(restBody, "department"),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if inviter.Nickname == "" {
		inviter.Nickname = "匿名用户"
	}

	detailJson, _ := json.Marshal(restBody)
	inviter.Detail = string(detailJson)

	if err := s.invitationDao.CreateJoiner(inviter); err != nil {
		return err
	}

	userUpdate := make(map[string]interface{})
	userUpdate["updated_at"] = time.Now()
	if nick := utils.GetString(restBody, "nickname"); nick != "" {
		userUpdate["name"] = nick
	}
	if avatar := utils.GetString(restBody, "avatar"); avatar != "" {
		userUpdate["avatar"] = avatar
	}
	s.userDao.UpdateField(openId, "updated_at", time.Now())

	return nil
}

func (s *InvitationService) ListInvitations(query map[string]interface{}, page, limit int) ([]model.Invitation, int64, error) {
	skip := (page - 1) * limit
	list, total, err := s.invitationDao.List(query, skip, limit)

	for i := range list {
		list[i].Status = s.computeStatus(&list[i])
	}

	return list, total, err
}

func (s *InvitationService) GetInvitationDetail(id string) (*model.Invitation, error) {
	inv, err := s.invitationDao.GetByID(id)
	if err != nil {
		return nil, err
	}
	inv.Status = s.computeStatus(inv)
	return inv, nil
}

func (s *InvitationService) CheckJoin(invitationId, openId string) (bool, error) {
	count, err := s.invitationDao.CountJoiners(invitationId, openId)
	return count > 0, err
}

func (s *InvitationService) GetJoiners(invitationId string) ([]model.Inviter, error) {
	return s.invitationDao.GetJoiners(invitationId, 100)
}

func (s *InvitationService) SaveLottery(id string, data map[string]interface{}, openId string) (string, error) {
	lottery := &model.Lottery{
		PrizeName:     utils.GetString(data, "prizeName"),
		WinnerCount:   utils.ParseInt(data, "winnerCount"),
		Status:        utils.GetString(data, "status"),
		CreatorOpenId: openId,
		UpdatedAt:     time.Now(),
	}

	winners := utils.GetString(data, "winners")
	lottery.Winners = winners

	if id != "" {
		var uintId uint
		fmt.Sscanf(id, "%d", &uintId)
		lottery.ID = uintId
	} else {
		lottery.InvitationId = utils.GetString(data, "invitationId")
		lottery.CreatedAt = time.Now()
	}

	if err := s.invitationDao.SaveLottery(lottery); err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", lottery.ID), nil
}

func (s *InvitationService) GetLotteryDetail(id uint) (*model.Lottery, error) {
	return s.invitationDao.GetLotteryByID(id)
}

func (s *InvitationService) DeleteLottery(id uint) error {
	return s.invitationDao.DeleteLottery(id)
}

func (s *InvitationService) DeleteInvitation(id uint) error {
	return s.invitationDao.Delete(id)
}

func (s *InvitationService) ListLotteries(query map[string]interface{}, page, limit int) ([]model.Lottery, int64, error) {
	skip := (page - 1) * limit
	return s.invitationDao.ListLotteries(query, skip, limit)
}

func (s *InvitationService) GetApplyConfig() (*model.ApplyConfig, error) {
	config, err := s.invitationDao.GetApplyConfig()
	if err != nil {
		return &model.ApplyConfig{
			DefaultOpen:  true,
			FieldOptions: `[{"label":"用户头像","value":"avatar","type":"image","enabled":true}, ...]`,
		}, nil
	}
	return config, nil
}

func (s *InvitationService) SaveApplyConfig(config *model.ApplyConfig) error {
	return s.invitationDao.SaveApplyConfig(config)
}
