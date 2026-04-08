package handler

import (
	"net/http"
	"wxcloudrun-golang/service"

	"github.com/gin-gonic/gin"
)

type DetailHandler struct {
	detailService *service.DetailService
	userService   *service.UserService
}

func NewDetailHandler(ds *service.DetailService, us *service.UserService) *DetailHandler {
	return &DetailHandler{
		detailService: ds,
		userService:   us,
	}
}

func (h *DetailHandler) GetDetail(c *gin.Context) {
	var query struct {
		OpenId string `json:"openId"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.OpenId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	record, err := h.detailService.GetMyDetail(query.OpenId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": "未填写今日的我"})
		return
	}

	// 补齐用户信息 (avatar, name)
	user, err := h.userService.GetOrCreateUser(query.OpenId, make(map[string]interface{}))
	if err == nil && user != nil {
		// 将 user 字段动态添加到响应中
		res := gin.H{
			"id":          record.ID,
			"openId":      record.OpenID,
			"plan":        record.Plan,
			"mood":        record.Mood,
			"style":       record.Style,
			"description": record.Description,
			"seat":        record.Seat,
			"image":       record.Image,
			"status":      record.Status,
			"avatar":      user.Avatar,
			"name":        user.Name,
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": res})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": record})
}

func (h *DetailHandler) SaveDetail(c *gin.Context) {
	// 直接使用 map 接收请求体
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	openId, ok := body["openId"].(string)
	if !ok || openId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	record, err := h.detailService.SaveDetail(openId, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": record})
}

func (h *DetailHandler) SaveSeat(c *gin.Context) {
	var query struct {
		OpenId string `json:"openId"`
		Seat   string `json:"seat"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.OpenId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	if err := h.detailService.SaveSeat(query.OpenId, query.Seat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "更新成功"})
}

func (h *DetailHandler) GetHash(c *gin.Context) {
	var query struct {
		PicHash string `json:"picHash"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.PicHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	pic, err := h.detailService.GetPictureHash(query.PicHash)
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
