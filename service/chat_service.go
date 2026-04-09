package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

type ChatService struct {
	chatDao *dao.ChatDao
}

func NewChatService() *ChatService {
	return &ChatService{
		chatDao: dao.NewChatDao(),
	}
}

func (s *ChatService) GetMessages(channel string, skip int, limit int) ([]model.ChatHistory, error) {
	if limit == 0 {
		limit = 20
	}
	return s.chatDao.GetHistory(channel, skip, limit)
}

func (s *ChatService) SendMessage(channelId, sender, receiver string, content model.MessageContent) (*model.ChatHistory, error) {
	timestamp := int64(time.Now().UnixMilli())
	if content.ID != "" {
		// 如果 content.ID 是时间戳字符串，尝试转换
		var t int64
		if _, err := fmt.Sscanf(content.ID, "%d", &t); err == nil {
			timestamp = t
		}
	}

	msg := &model.ChatHistory{
		ChannelId:      channelId,
		SenderOpenID:   sender,
		ReceiverOpenID: receiver,
		MessageContent: content,
		Timestamp:      timestamp,
		CreatedAt:      time.Now(),
	}

	if err := s.chatDao.CreateMessage(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *ChatService) UpdateStateByReceiver(channelId, openId string, state int) error {
	return s.chatDao.UpdateMessageStateByReceiver(channelId, openId, state)
}

func (s *ChatService) UpdateStateByID(channelId, msgId string, state int) error {
	return s.chatDao.UpdateMessageStateByID(channelId, msgId, state)
}

func (s *ChatService) BlockUser(openId, blockId string) error {
	return s.chatDao.SaveBlock(openId, blockId)
}

func (s *ChatService) DeleteInfoMonitor(openId string) error {
	return s.chatDao.DeleteInfoMonitor(openId)
}

func (s *ChatService) CheckContent(content string, scene int, openId string) (map[string]interface{}, error) {
	// 实际调用微信文本检测接口
	apiUrl := "http://api.weixin.qq.com/wxa/msg_sec_check"

	// 构建请求体
	requestBody := map[string]interface{}{
		"content": content,
		"version": 2,
		"scene":   scene,
		"openid":  openId,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %v", err)
	}

	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("http post failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %v", err)
	}

	return result, nil
}

// UpdatePictureSecCheckStatus 根据 trace_id 更新图片审核状态
func (s *ChatService) UpdatePictureSecCheckStatus(traceID string, status int) error {
	// 这里需要通过 trace_id 找到对应的图片记录并更新
	// 暂时通过 DAO 直接更新，或者可以调用 mediaService
	return nil // TODO: 实现更新逻辑
}

// CreateInfoMonitor 创建消息通知记录
func (s *ChatService) CreateInfoMonitor(openId, source string) error {
	return s.chatDao.CreateInfoMonitor(openId, source)
}
