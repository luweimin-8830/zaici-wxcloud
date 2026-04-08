package handler

import (
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
