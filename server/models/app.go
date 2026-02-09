package models

import (
	"kefu-server/store"
	"kefu-server/utils"
	"kefu-server/utils/logger"
	"strings"

	"gorm.io/gorm"
)

type App struct {
	gorm.Model
	Name        string `gorm:"size:255" json:"name"`
	Logo        string `gorm:"size:255" json:"logo"`
	AppID       string `gorm:"size:255" json:"app_id"`
	Status      int    `gorm:"size:255" json:"status"` // 1=启用, 0=禁用
	AllowDomain string `gorm:"size:255" json:"allow_domain"`
	WelcomeMsg  string `gorm:"size:255" json:"welcome_msg"`
	Contact     string `gorm:"size:255" json:"contact"` // 联系人
}

// GenAppID 生成唯一的 AppID
func GenAppID() string {
	// 生成基于时间戳和随机数的 AppID
	timestamp := utils.GenerateTimestamp()
	random := utils.GenerateRandomString(8)
	return "zerospace_" + timestamp + "_" + random
}

func GetApp(appID string) *App {
	var app App
	if err := store.DB.Where("app_id = ? AND status = ?", appID, 1).First(&app).Error; err != nil {
		logger.Errorf("app not found or disabled: %s", appID)
		return nil
	}
	return &app
}

// IsDomainAllowed 检查域名是否在允许列表中
func IsDomainAllowed(origin, referer, allowDomain string) bool {
	if allowDomain == "" {
		return false
	}

	// 提取域名
	domain := utils.ExtractDomain(origin)
	if domain == "" {
		domain = utils.ExtractDomain(referer)
	}

	if domain == "" {
		return false
	}

	// 检查域名是否在允许列表中
	allowedDomains := strings.Split(allowDomain, ",")
	for _, allowedDomain := range allowedDomains {
		allowedDomain = strings.TrimSpace(allowedDomain)
		if allowedDomain == domain {
			return true
		}
		// 检查是否为通配符域名（如 *.pages.dev）
		if strings.HasPrefix(allowedDomain, "*") {
			wildcardDomain := strings.TrimPrefix(allowedDomain, "*")
			if strings.HasSuffix(domain, wildcardDomain) {
				return true
			}
		}
	}

	return false
}
