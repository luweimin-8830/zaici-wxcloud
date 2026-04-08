package handler

import (
	"net/http"
	"strconv"
	"wxcloudrun-golang/service"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

type ShopHandler struct {
	shopService *service.ShopService
}

func NewShopHandler() *ShopHandler {
	return &ShopHandler{
		shopService: service.NewShopService(),
	}
}

func (h *ShopHandler) GetDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	shop, err := h.shopService.GetDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该门店"})
		return
	}

	c.JSON(http.StatusOK, shop)
}

func (h *ShopHandler) Save(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	shopList, ok := data["shopList"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	shop, err := h.shopService.Save(shopList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存门店失败"})
		return
	}

	c.JSON(http.StatusOK, shop)
}

func (h *ShopHandler) Update(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	shopList, ok := data["shopList"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	idStr := utils.GetString(shopList, "_id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID 格式错误"})
		return
	}

	if err := h.shopService.Update(uint(id), shopList); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新门店失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ShopHandler) Delete(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	idStr := utils.GetString(data, "id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.shopService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除门店失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ShopHandler) Admin(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	page := utils.ParseInt(data, "page")
	limit := utils.ParseInt(data, "limit")
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	keyword := utils.GetString(data, "keyword")

	shops, total, err := h.shopService.AdminList(page, limit, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}

	response := map[string]interface{}{
		"list":  shops,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ShopHandler) GetNearList(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	longitude := 0.0
	latitude := 0.0
	distance := 1000.0

	if lons, ok := data["longitude"].(float64); ok {
		longitude = lons
	}
	if lats, ok := data["latitude"].(float64); ok {
		latitude = lats
	}
	if dists, ok := data["distance"].(float64); ok {
		distance = dists
	}

	shops, err := h.shopService.GetNearList(longitude, latitude, distance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取附近门店失败"})
		return
	}

	c.JSON(http.StatusOK, shops)
}
