package handler

import (
	"net/http"
	"strconv"
	"wxcloudrun-golang/service"

	"github.com/gin-gonic/gin"
)

type BannerHandler struct {
	bannerService *service.BannerService
}

func NewBannerHandler(bs *service.BannerService) *BannerHandler {
	return &BannerHandler{bannerService: bs}
}

func (h *BannerHandler) ListBanners(c *gin.Context) {
	banners, err := h.bannerService.ListBanners()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": banners})
}

func (h *BannerHandler) GetBannerDetail(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "缺少参数ID"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "ID格式错误"})
		return
	}

	banner, err := h.bannerService.GetBannerDetail(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "未找到该广告详情"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": banner})
}

func (h *BannerHandler) SaveBanner(c *gin.Context) {
	var body struct {
		Banner map[string]interface{} `json:"banner"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Banner == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	msg, err := h.bannerService.SaveBanner(body.Banner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": msg})
}

func (h *BannerHandler) DeleteBanner(c *gin.Context) {
	var query struct {
		ID string `json:"id"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	id, err := strconv.Atoi(query.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "ID格式错误"})
		return
	}

	if err := h.bannerService.DeleteBanner(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "删除成功"})
}
