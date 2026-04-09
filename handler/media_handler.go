package handler

import (
	"bytes"
	"encoding/json"
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

// GetPictureByHash 根据哈希获取图片（新接口，返回完整数据）
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

// GetHash 根据哈希获取图片URL（兼容旧接口，从detail_handler迁移）
func (h *MediaHandler) GetHash(c *gin.Context) {
	var query struct {
		PicHash string `json:"picHash"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.PicHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	pic, err := h.mediaService.GetPictureByHash(query.PicHash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": "无相应图片"})
		return
	}

	if pic.SecCheckStatus == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "图片存在违规行为,禁止发布"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": pic.UserPicURL})
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

// StartCensorRequest 启动内容审核请求
type StartCensorRequest struct {
	FileID string `json:"fileid" binding:"required"`
	Digest string `json:"digest" binding:"required"`
}

// StartCensor 启动图片内容审核（从 index.js 迁移）
func (h *MediaHandler) StartCensor(c *gin.Context) {
	var req StartCensorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	openId := c.GetHeader("x-wx-openid")
	if openId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少 openid"})
		return
	}

	// 1. 先查询数据库是否已有相同 hash 的图片
	existingPic, err := h.mediaService.GetPictureByHash(req.Digest)
	if err == nil && existingPic != nil {
		// 已有相同文件上传
		c.JSON(http.StatusOK, gin.H{
			"code":           200,
			"message":        "the same file has been uploaded already",
			"secCheckStatus": existingPic.SecCheckStatus,
		})
		return
	}

	// 2. 没有相同文件，返回状态 2（待审核）
	c.JSON(http.StatusOK, gin.H{
		"code":           200,
		"message":        "there is no same file uploaded",
		"secCheckStatus": 2,
	})

	// 3. 获取临时文件 URL（这里简化处理，实际需要调用微信云托管 API）
	// 注意：在 Go 云托管中，fileid 就是 cloud:// 开头的文件 ID
	// 需要调用微信接口获取临时 URL，然后调用内容安全检测

	// 4. 调用微信内容安全检测 API
	wxReqBody := map[string]interface{}{
		"media_url":  req.FileID, // 这里应该是临时 URL，需要额外处理
		"media_type": 2,
		"openid":     openId,
		"version":    2,
		"scene":      1,
	}

	wxReqJSON, _ := json.Marshal(wxReqBody)
	resp, err := http.Post(
		"http://api.weixin.qq.com/wxa/media_check_async",
		"application/json",
		bytes.NewBuffer(wxReqJSON),
	)
	if err != nil {
		// 记录错误但不返回，因为已经返回了 200
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var wxResp struct {
		TraceID string `json:"trace_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wxResp); err != nil {
		return
	}

	// 5. 创建图片记录到数据库
	h.mediaService.CreatePictureWithStatus(req.FileID, req.Digest, wxResp.TraceID, 2)
}
