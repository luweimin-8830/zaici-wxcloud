package handler

import (
	"net/http"
	"strconv"
	"time"
	"wxcloudrun-golang/service"

	"github.com/gin-gonic/gin"
)

type OnlineHandler struct {
	onlineService *service.OnlineService
}

func NewOnlineHandler(os *service.OnlineService) *OnlineHandler {
	return &OnlineHandler{onlineService: os}
}

func (h *OnlineHandler) GetStatus(c *gin.Context) {
	var query struct {
		OpenId string `json:"openId"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.OpenId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	res, err := h.onlineService.GetStatus(query.OpenId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": res})
}

func (h *OnlineHandler) GetNear(c *gin.Context) {
	var query struct {
		OpenId    string  `json:"openId"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		ShopId    string  `json:"shopId"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.OpenId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	list, err := h.onlineService.GetNearList(query.OpenId, query.Longitude, query.Latitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"data": list}})
}

func (h *OnlineHandler) SaveOnline(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	openId := c.GetHeader("x-wx-openid")
	if openId == "" {
		if oid, ok := body["openId"].(string); ok {
			openId = oid
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "缺少 openId"})
			return
		}
	}

	res, err := h.onlineService.SaveOnline(openId, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": res})
}

func (h *OnlineHandler) UpdateOnline(c *gin.Context) {
	var query struct {
		ID string `json:"id"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	idUint, _ := strconv.ParseUint(query.ID, 10, 64)
	now := time.Now().UnixMilli()
	err := h.onlineService.UpdateOnline(uint(idUint), now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"online": "updated", "message": "更新成功"}})
}

func (h *OnlineHandler) GetHistory(c *gin.Context) {
	var query struct {
		ShopId string `json:"shopId"`
		OpenId string `json:"openId"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.ShopId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误: 缺失 shopId"})
		return
	}

	list, err := h.onlineService.GetHistory(query.ShopId, query.OpenId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": list})
}

func (h *OnlineHandler) GetShop(c *gin.Context) {
	var query struct {
		OpenId string `json:"openId"`
		ShopId string `json:"shopId"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.OpenId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	list, err := h.onlineService.GetShopOnline(query.ShopId, query.OpenId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"data": list}})
}
