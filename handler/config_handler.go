package handler

import (
	"wxcloudrun-golang/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ConfigHandler struct {
	configService *service.ConfigService
}

func NewConfigHandler(cs *service.ConfigService) *ConfigHandler {
	return &ConfigHandler{configService: cs}
}

func (h *ConfigHandler) GetDistance(c *gin.Context) {
	config, err := h.configService.GetDistance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": config})
}

func (h *ConfigHandler) SaveDistance(c *gin.Context) {
	var body struct {
		Distance interface{} `json:"distance"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var dist int
	switch v := body.Distance.(type) {
	case float64:
		dist = int(v)
	case string:
		d, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "distance 必须为数字"})
			return
		}
		dist = d
	case int:
		dist = v
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "distance 类型错误"})
		return
	}

	if err := h.configService.SaveDistance(dist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "更新成功"})
}
