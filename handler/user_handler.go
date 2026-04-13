package handler

import (
	"net/http"
	"wxcloudrun-golang/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

func (h *UserHandler) GetOrCreateUser(c *gin.Context) {
	openId := c.GetHeader("x-wx-openid")
	if openId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未获取到openId,请确认"})
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		body = make(map[string]interface{})
	}

	user, err := h.userService.GetOrCreateUser(openId, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
}

func (h *UserHandler) SaveUser(c *gin.Context) {
	var query struct {
		OpenId string      `json:"openId"`
		Key    string      `json:"key"`
		Data   interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.OpenId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	if err := h.userService.SaveUserInfo(query.OpenId, query.Key, query.Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "更新成功"})
}

func (h *UserHandler) UseSuperLike(c *gin.Context) {
	var query struct {
		OpenId string `json:"openId"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.OpenId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	msg, err := h.userService.UseSuperLike(query.OpenId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": msg})
}

func (h *UserHandler) AddSuperLike(c *gin.Context) {
	var query struct {
		UserList []struct {
			OpenId string `json:"openId"`
			Total  int    `json:"total"`
		} `json:"userList"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || len(query.UserList) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	// 转换为 service.SuperLikeRequest 切片
	var reqList []service.SuperLikeRequest
	for _, u := range query.UserList {
		reqList = append(reqList, service.SuperLikeRequest{
			OpenId: u.OpenId,
			Total:  u.Total,
		})
	}

	if err := h.userService.AddSuperLike(reqList); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "发放成功"})
}

func (h *UserHandler) AddAdmin(c *gin.Context) {
	openId := c.GetHeader("x-wx-openid")
	if openId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未获取到openId"})
		return
	}

	msg, err := h.userService.ManageAdmin(openId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": msg})
}

func (h *UserHandler) GetOtherUserInfo(c *gin.Context) {
	var query struct {
		Id string `json:"id"`
	}
	if err := c.ShouldBindJSON(&query); err != nil || query.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 401, "message": "参数错误"})
		return
	}

	user, err := h.userService.GetOtherUserInfo(query.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
}

// GetUserInfo 获取当前登录用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	openId := c.GetHeader("x-wx-openid")
	if openId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未获取到openId,请确认"})
		return
	}

	user, err := h.userService.GetUserInfo(openId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
}
