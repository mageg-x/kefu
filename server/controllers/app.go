package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"kefu-server/models"
	"kefu-server/store"
	"kefu-server/utils/logger"
	"kefu-server/utils/response"
)

type AppController struct{}

type AppRequest struct {
	Name        string `json:"name" binding:"required"`
	AppID       string `json:"app_id"`
	Logo        string `json:"logo"`
	AllowDomain string `json:"allow_domain"`
	WelcomeMsg  string `json:"welcome_msg"`
	Contact     string `json:"contact"`
	Status      int    `json:"status" binding:"required,oneof=0 1"`
}

// GetApps 获取应用列表
func (ac *AppController) GetApps(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	statusStr := c.Query("status")

	// 构建查询
	query := store.DB.Model(&models.App{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR app_id LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err == nil {
			query = query.Where("status = ?", status)
		}
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Errorf("count apps failed: %v", err)
		response.ResponseError(c, http.StatusInternalServerError, response.ErrCodeInternalError)
		return
	}

	// 分页查询
	offset := (page - 1) * pageSize
	var apps []models.App
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&apps).Error; err != nil {
		logger.Errorf("get apps failed: %v", err)
		response.ResponseError(c, http.StatusInternalServerError, response.ErrCodeInternalError)
		return
	}

	logger.Infof("get apps successful, page: %d, page_size: %d, total: %d", page, pageSize, total)
	response.ResponseSuccess(c, gin.H{
		"data":  apps,
		"total": total,
	})
}

// CreateApp 创建应用
func (ac *AppController) CreateApp(c *gin.Context) {
	// 检查管理员权限
	if !IsAdmin(c) {
		logger.Errorf("permission denied, not admin")
		response.ResponseError(c, http.StatusForbidden, response.ErrCodeForbidden)
		return
	}

	var req AppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("create app request parameter error: %v", err)
		response.ResponseError(c, http.StatusBadRequest, response.ErrCodeInvalidParams)
		return
	}

	// 生成 AppID（如果未提供）
	appID := req.AppID
	if appID == "" {
		appID = models.GenAppID()
	}

	// 检查 AppID 是否已存在
	var existingApp models.App
	if err := store.DB.Where("app_id = ?", appID).First(&existingApp).Error; err == nil {
		logger.Errorf("app id already exists: %s", appID)
		response.ResponseError(c, http.StatusBadRequest, response.ErrCodeInvalidParams)
		return
	}

	// 创建应用
	app := models.App{
		Name:        req.Name,
		AppID:       appID,
		Logo:        req.Logo,
		AllowDomain: req.AllowDomain,
		WelcomeMsg:  req.WelcomeMsg,
		Contact:     req.Contact,
		Status:      req.Status,
	}

	if err := store.DB.Create(&app).Error; err != nil {
		logger.Errorf("create app failed: %v", err)
		response.ResponseError(c, http.StatusInternalServerError, response.ErrCodeInternalError)
		return
	}

	logger.Infof("create app successful: %s", app.Name)
	response.ResponseSuccess(c, app)
}

// UpdateApp 更新应用
func (ac *AppController) UpdateApp(c *gin.Context) {
	// 检查管理员权限
	if !IsAdmin(c) {
		logger.Errorf("permission denied, not admin")
		response.ResponseError(c, http.StatusForbidden, response.ErrCodeForbidden)
		return
	}

	var req struct {
		AppID       string `json:"app_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Logo        string `json:"logo"`
		AllowDomain string `json:"allow_domain"`
		WelcomeMsg  string `json:"welcome_msg"`
		Contact     string `json:"contact"`
		Status      int    `json:"status" binding:"required,oneof=0 1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("update app request parameter error: %v", err)
		response.ResponseError(c, http.StatusBadRequest, response.ErrCodeInvalidParams)
		return
	}

	// 检查应用是否存在
	var app models.App
	if err := store.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		logger.Errorf("app not found: %s", req.AppID)
		response.ResponseError(c, http.StatusNotFound, response.ErrCodeNotFound)
		return
	}

	// 更新应用
	updates := map[string]interface{}{
		"Name":        req.Name,
		"Logo":        req.Logo,
		"AllowDomain": req.AllowDomain,
		"WelcomeMsg":  req.WelcomeMsg,
		"Contact":     req.Contact,
		"Status":      req.Status,
	}

	if err := store.DB.Model(&app).Updates(updates).Error; err != nil {
		logger.Errorf("update app failed: %v", err)
		response.ResponseError(c, http.StatusInternalServerError, response.ErrCodeInternalError)
		return
	}

	// 重新获取更新后的数据
	if err := store.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		logger.Errorf("get updated app failed: %v", err)
		response.ResponseError(c, http.StatusInternalServerError, response.ErrCodeInternalError)
		return
	}

	logger.Infof("update app successful: %s", req.AppID)
	response.ResponseSuccess(c, app)
}

// DeleteApp 删除应用
func (ac *AppController) DeleteApp(c *gin.Context) {
	// 检查管理员权限
	if !IsAdmin(c) {
		logger.Errorf("permission denied, not admin")
		response.ResponseError(c, http.StatusForbidden, response.ErrCodeForbidden)
		return
	}

	// 获取查询参数
	appID := c.Query("app_id")
	if appID == "" {
		logger.Errorf("app_id is required")
		response.ResponseError(c, http.StatusBadRequest, response.ErrCodeInvalidParams)
		return
	}

	// 检查应用是否存在
	var app models.App
	if err := store.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		logger.Errorf("app not found: %s", appID)
		response.ResponseError(c, http.StatusNotFound, response.ErrCodeNotFound)
		return
	}

	// 删除应用
	if err := store.DB.Unscoped().Delete(&app).Error; err != nil {
		logger.Errorf("delete app failed: %v", err)
		response.ResponseError(c, http.StatusInternalServerError, response.ErrCodeInternalError)
		return
	}

	logger.Infof("delete app successful: %s", appID)
	response.ResponseSuccess(c, gin.H{"message": "delete successful"})
}

// GetConfig 获取应用配置（前端 widget 接入接口）
func (ac *AppController) GetConfig(c *gin.Context) {
	// 获取请求参数
	appID := c.Query("appid")
	if appID == "" {
		logger.Errorf("appid is required")
		response.ResponseError(c, http.StatusBadRequest, response.ErrCodeInvalidParams)
		return
	}

	// 获取 Origin 和 Referer 头
	origin := c.GetHeader("Origin")
	referer := c.GetHeader("Referer")

	// 检查请求来源
	if origin == "" && referer == "" {
		logger.Errorf("origin or referer is required")
		response.ResponseError(c, http.StatusForbidden, response.ErrCodeForbidden)
		return
	}

	// 获取应用信息
	var app models.App
	if err := store.DB.Where("app_id = ? AND status = ?", appID, 1).First(&app).Error; err != nil {
		logger.Errorf("app not found or disabled: %s", appID)
		response.ResponseError(c, http.StatusForbidden, response.ErrCodeForbidden)
		return
	}

	// 检查域名是否在允许列表中
	if !models.IsDomainAllowed(origin, referer, app.AllowDomain) {
		logger.Errorf("domain not allowed: origin=%s, referer=%s", origin, referer)
		response.ResponseError(c, http.StatusForbidden, response.ErrCodeForbidden)
		return
	}

	// 返回配置信息
	response.ResponseSuccess(c, gin.H{
		"name":        app.Name,
		"logo":        app.Logo,
		"welcome_msg": app.WelcomeMsg,
	})
}
