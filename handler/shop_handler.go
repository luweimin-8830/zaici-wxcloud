package handler

import (
	"fmt"
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

	// 从 data["id"] 获取 ID，避免在 shopList 中传递 _id
	idStr := utils.GetString(data, "id")
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
		// 如果绑定失败，尝试从 Query 参数获取
		data = make(map[string]interface{})
	}

	// 优先从 Query 参数获取（支持 GET 请求）
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	keyword := c.Query("keyword")

	page := 1
	limit := 10

	if pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &page)
	} else {
		page = utils.ParseInt(data, "page")
	}

	if limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	} else {
		limit = utils.ParseInt(data, "limit")
	}

	if keyword == "" {
		keyword = utils.GetString(data, "keyword")
	}

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	fmt.Printf("Admin API - page: %d, limit: %d, keyword: %s\n", page, limit, keyword)

	shops, total, err := h.shopService.AdminList(page, limit, keyword)
	if err != nil {
		fmt.Printf("AdminList error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败", "detail": err.Error()})
		return
	}

	fmt.Printf("AdminList result - shops count: %d, total: %d\n", len(shops), total)

	response := map[string]interface{}{
		"code":  0,
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
