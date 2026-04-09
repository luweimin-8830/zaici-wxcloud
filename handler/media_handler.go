package handler

import (
	"net/http"
	"strconv"
	"wxcloudrun-golang/service"

	"github.com/gin-gonic/gin"
)

// MediaHandler 媒体处理器
type MediaHandler struct {
	mediaService *service.MediaService
}

// NewMediaHandler 创建媒体处理器
func NewMediaHandler() *MediaHandler {
	return &MediaHandler{
		mediaService: service.NewMediaService(),
	}
}

// CreatePictureRequest 创建图片请求
type CreatePictureRequest struct {
	UserPicURL string `json:"userPicUrl" binding:"required"`
	PicHash    string `json:"picHash" binding:"required"`
	TraceID    string `json:"traceId"`
}

// UpdatePictureRequest 更新图片请求
type UpdatePictureRequest struct {
	UserPicURL string `json:"userPicUrl"`
	PicHash    string `json:"picHash"`
	TraceID    string `json:"traceId"`
}

// UpdateSecCheckStatusRequest 更新审核状态请求
type UpdateSecCheckStatusRequest struct {
	SecCheckStatus int `json:"secCheckStatus" binding:"required"`
}

// CreatePicture 创建图片
func (h *MediaHandler) CreatePicture(c *gin.Context) {
	var req CreatePictureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	picture, err := h.mediaService.CreatePicture(req.UserPicURL, req.PicHash, req.TraceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": picture})
}

// GetPictureByID 根据ID获取图片
func (h *MediaHandler) GetPictureByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "ID参数错误"})
		return
	}

	picture, err := h.mediaService.GetPictureByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "图片不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "获取成功", "data": picture})
}

// GetPictureByHash 根据哈希获取图片
func (h *MediaHandler) GetPictureByHash(c *gin.Context) {
	picHash := c.Query("picHash")
	if picHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "picHash参数不能为空"})
		return
	}

	picture, err := h.mediaService.GetPictureByHash(picHash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "图片不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "获取成功", "data": picture})
}

// UpdatePicture 更新图片
func (h *MediaHandler) UpdatePicture(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "ID参数错误"})
		return
	}

	var req UpdatePictureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 先获取原记录
	picture, err := h.mediaService.GetPictureByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "图片不存在"})
		return
	}

	// 更新字段
	if req.UserPicURL != "" {
		picture.UserPicURL = req.UserPicURL
	}
	if req.PicHash != "" {
		picture.PicHash = req.PicHash
	}
	if req.TraceID != "" {
		picture.TraceID = req.TraceID
	}

	if err := h.mediaService.UpdatePicture(picture); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": picture})
}

// UpdatePictureSecCheckStatus 更新图片审核状态
func (h *MediaHandler) UpdatePictureSecCheckStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "ID参数错误"})
		return
	}

	var req UpdateSecCheckStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	if err := h.mediaService.UpdatePictureSecCheckStatus(uint(id), req.SecCheckStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

// DeletePicture 删除图片
func (h *MediaHandler) DeletePicture(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "ID参数错误"})
		return
	}

	if err := h.mediaService.DeletePicture(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ListPictures 获取图片列表
func (h *MediaHandler) ListPictures(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	pictures, err := h.mediaService.ListPictures(pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取失败", "error": err.Error()})
		return
	}

	total, _ := h.mediaService.CountPictures()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data": gin.H{
			"list":     pictures,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}
