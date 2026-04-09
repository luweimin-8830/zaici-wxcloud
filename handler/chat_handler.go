package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/service"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(cs *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: cs}
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	var query struct {
		Channel string `json:"channel"`
		Length  int    `json:"length"`
		Cont    int    `json:"cont"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.Channel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少频道ID"})
		return
	}

	list, err := h.chatService.GetMessages(query.Channel, query.Length, query.Cont)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": list})
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	var query struct {
		ChannelId      string `json:"channelId"`
		SenderOpenID   string `json:"senderOpenID"`
		ReceiverOpenID string `json:"receiverOpenID"`
		MessageContent struct {
			ID          string `json:"id"`
			Content     string `json:"content"`
			Type        string `json:"type"`
			ContentType string `json:"contentType"`
			Pic         string `json:"pic"`
			Name        string `json:"name"`
			State       int    `json:"state"`
		} `json:"messageContent"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.ChannelId == "" || query.SenderOpenID == "" || query.ReceiverOpenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 转换为 model.MessageContent
	content := model.MessageContent{
		ID:          query.MessageContent.ID,
		Content:     query.MessageContent.Content,
		Type:        query.MessageContent.Type,
		ContentType: query.MessageContent.ContentType,
		Pic:         query.MessageContent.Pic,
		Name:        query.MessageContent.Name,
		State:       query.MessageContent.State,
	}

	result, err := h.chatService.SendMessage(query.ChannelId, query.SenderOpenID, query.ReceiverOpenID, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
}

func (h *ChatHandler) UpdateState(c *gin.Context) {
	var query struct {
		ChannelId string `json:"channelId"`
		OpenId    string `json:"openId"`
		State     int    `json:"state"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.ChannelId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	if err := h.chatService.UpdateStateByReceiver(query.ChannelId, query.OpenId, query.State); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "更新成功"})
}

func (h *ChatHandler) UpdateStateByID(c *gin.Context) {
	var query struct {
		ChannelId string `json:"channelId"`
		ID        string `json:"id"`
		State     int    `json:"state"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.ChannelId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	if err := h.chatService.UpdateStateByID(query.ChannelId, query.ID, query.State); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "更新成功"})
}

func (h *ChatHandler) SaveBlock(c *gin.Context) {
	var query struct {
		OpenId  string `json:"openId"`
		BlockId string `json:"blockId"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.OpenId == "" || query.BlockId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	if err := h.chatService.BlockUser(query.OpenId, query.BlockId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "屏蔽成功"})
}

func (h *ChatHandler) DeleteInfo(c *gin.Context) {
	openId := c.GetHeader("x-wx-openid")
	if openId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	if err := h.chatService.DeleteInfoMonitor(openId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "删除成功"})
}

func (h *ChatHandler) CheckContent(c *gin.Context) {
	var query struct {
		Content string `json:"content"`
		Scene   int    `json:"scene"`
		OpenId  string `json:"openId"`
	}
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	res, err := h.chatService.CheckContent(query.Content, query.Scene, query.OpenId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": res})
}

// CensorRequest 微信内容安全回调请求
type CensorRequest struct {
	Action    string `json:"action"`
	MsgType   string `json:"MsgType"`
	Event     string `json:"Event"`
	TraceID   string `json:"trace_id"`
	Result    struct {
		Suggest string `json:"suggest"`
	} `json:"result"`
	FromUserName string `json:"FromUserName"`
	ToUserName   string `json:"ToUserName"`
}

// Censor 微信内容安全回调（从 index.js 迁移）
func (h *ChatHandler) Censor(c *gin.Context) {
	var req CensorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{}) // 返回空响应避免微信重试
		return
	}

	// 1. 微信服务健康检查
	if req.Action == "CheckContainerPath" {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	// 2. 微信内容安全检查回调
	if req.MsgType == "event" && req.Event == "wxa_media_check" {
		// 异步更新审核状态
		status := 0
		if req.Result.Suggest == "pass" {
			status = 1
		}
		h.chatService.UpdatePictureSecCheckStatus(req.TraceID, status)
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	// 3. 消息转发到客服
	if req.MsgType == "text" || req.MsgType == "image" || (req.Event == "user_enter_tempsession" && req.MsgType == "event") {
		c.JSON(http.StatusOK, gin.H{
			"ToUserName":   req.FromUserName,
			"FromUserName": req.ToUserName,
			"CreateTime":   0,
			"MsgType":      "transfer_customer_service",
		})
		return
	}

	// 默认回复
	c.JSON(http.StatusOK, gin.H{})
}

// WebhookRequest GoEasy 回调请求
type WebhookRequest struct {
	Content string `json:"content"`
}

// Webhook GoEasy 消息回调（从 index.js 迁移）
func (h *ChatHandler) Webhook(c *gin.Context) {
	// 获取签名
	receivedSignature := c.GetHeader("x-goeasy-signature")
	if receivedSignature == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing signature"})
		return
	}

	// 读取 body
	body, _ := c.GetRawData()

	// 验证签名（HMAC-SHA1 + Base64）
	const secretKey = "e42487a150c54d35"
	hmac := hmac.New(sha1.New, []byte(secretKey))
	hmac.Write(body)
	expectedSignature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))

	// 安全比较（防时序攻击）
	if !bytes.Equal([]byte(receivedSignature), []byte(expectedSignature)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"content": "success",
	})

	// 解析消息内容
	var req WebhookRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return
	}

	var messages []struct {
		Channel   string `json:"channel"`
		Timestamp int64  `json:"timestamp"`
		Content   string `json:"content"`
	}
	if err := json.Unmarshal([]byte(req.Content), &messages); err != nil {
		return
	}

	// 处理每条消息
	for _, msg := range messages {
		var msgContent struct {
			ID           string `json:"id"`
			Content      string `json:"content"`
			Type         string `json:"type"`
			ContentType  string `json:"contentType"`
			Pic          string `json:"pic"`
			Name         string `json:"name"`
			State        int    `json:"state"`
			OpenID       string `json:"openID"`
			ReceiverOpenID string `json:"receiverOpenID"`
		}
		if err := json.Unmarshal([]byte(msg.Content), &msgContent); err != nil {
			continue
		}

		// 保存聊天记录
		content := model.MessageContent{
			ID:          msgContent.ID,
			Content:     msgContent.Content,
			Type:        msgContent.Type,
			ContentType: msgContent.ContentType,
			Pic:         msgContent.Pic,
			Name:        msgContent.Name,
			State:       msgContent.State,
		}
		h.chatService.SendMessage(msg.Channel, msgContent.OpenID, msgContent.ReceiverOpenID, content)

		// 创建消息通知
		h.chatService.CreateInfoMonitor(msgContent.ReceiverOpenID, "chat")
	}
}
