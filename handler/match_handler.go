package handler

import (
	"net/http"
	"wxcloudrun-golang/service"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	matchService *service.MatchService
}

func NewMatchHandler(ms *service.MatchService) *MatchHandler {
	return &MatchHandler{matchService: ms}
}

func (h *MatchHandler) GetMatches(c *gin.Context) {
	var query struct {
		OpenId string `json:"openId"`
		Status int    `json:"status"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.OpenId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	list, err := h.matchService.GetMatches(query.OpenId, query.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": list})
}

func (h *MatchHandler) AddMatch(c *gin.Context) {
	var query struct {
		OpenId1   string `json:"openId1"`
		OpenId2   string `json:"openId2"`
		Operation int    `json:"operation"`
		Channel   string  `json:"channel"`
		LikeType  int    `json:"likeType"`
	}
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	status, err := h.matchService.AddMatch(query.OpenId1, query.OpenId2, query.Channel, query.Operation, query.LikeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": status})
}

func (h *MatchHandler) DeleteMatch(c *gin.Context) {
	var query struct {
		Channel string `json:"channel"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.Channel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	if err := h.matchService.DeleteMatch(query.Channel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "解除匹配成功"})
}

func (h *MatchHandler) GetLikeMatch(c *gin.Context) {
	var query struct {
		OpenId string `json:"openId"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.OpenId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	list, err := h.matchService.GetLikeMatchList(query.OpenId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": list})
}

func (h *MatchHandler) SendMessage(c *gin.Context) {
	var query struct {
		OpenId   string `json:"openId"`
		ShopName string `json:"shopName"`
	}
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	res, err := h.matchService.SendSubscribeMessage(query.OpenId, query.ShopName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": res})
}
