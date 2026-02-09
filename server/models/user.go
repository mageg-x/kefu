package models

import (
	"crypto/sha256"
	"encoding/hex"

	"gorm.io/gorm"

	"kefu-server/utils/logger"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`     // 明文密码
	Avatar   string `gorm:"size:255" json:"avatar"`         // 头像
	Role     string `gorm:"size:50;not null" json:"role"`   // agent or admin
	Status   int    `gorm:"size:50;not null" json:"status"` // 1、在席 2、离席
	Active   bool   `gorm:"default:true" json:"active"`     // 1、激活 0、禁用
	Apps     string `gorm:"type:text" json:"apps"`          // 客服负责的业务, 格式位json字符串数组， 范围 缺省 ["all"]
}

// hashPassword 使用SHA256 hash密码
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func CreateDefaultUsers(db *gorm.DB) error {
	var count int64
	db.Model(&User{}).Count(&count)

	if count == 0 {
		users := []User{
			{
				Username: "admin",
				Password: "12345678", // 存储明文密码
				Role:     "admin",
				Avatar:   "https://api.dicebear.com/7.x/avataaars/svg?seed=admin",
			},
			{
				Username: "agent",
				Password: "12345678", // 存储明文密码
				Role:     "agent",
				Avatar:   "https://api.dicebear.com/7.x/avataaars/svg?seed=agent",
			},
		}

		if err := db.Create(&users).Error; err != nil {
			logger.Errorf("create default users failed: %v", err)
			return err
		}

		logger.Infof("default users created successfully")
	}

	return nil
}
