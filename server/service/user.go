package service

import (
	"fmt"
	"strings"

	"github.com/golang-infrastructure/go-shuffle"

	"kefu-server/models"
	"kefu-server/store"
	"kefu-server/utils/logger"
)

type UserService struct {
}

var (
	instUserService *UserService
)

func GetUserService() *UserService {
	if instUserService == nil {
		instUserService = &UserService{}
	}
	return instUserService
}

func (us *UserService) GetUser(username string) (*models.User, error) {
	var user models.User
	if err := store.DB.Where("username = ?", username).First(&user).Error; err != nil {
		logger.Errorf("user does not exist: %s", username)
		return nil, fmt.Errorf("user does not exist: %s", username)
	}
	return &user, nil
}

func (us *UserService) CreateUser(username, password, avatar, role string, active bool) (*models.User, error) {
	user := models.User{
		Username: username,
		Password: password,
		Role:     role,
		Active:   active,
		Avatar:   avatar,
	}
	if err := store.DB.Create(&user).Error; err != nil {
		logger.Errorf("create user failed: %s", username)
		return nil, fmt.Errorf("create user failed: %s", username)
	}
	return us.GetUser(username)
}

func (us *UserService) UpdateUser(username, password, avatar, role string, active bool) (*models.User, error) {
	user := models.User{
		Username: username,
		Password: password,
		Role:     role,
		Active:   active,
		Avatar:   avatar,
	}
	if err := store.DB.Model(&models.User{}).Where("username = ?", username).Updates(&user).Error; err != nil {
		logger.Errorf("update user failed: %s", username)
		return nil, fmt.Errorf("update user failed: %s", username)
	}
	return us.GetUser(username)

}

func (us *UserService) SetUserStatus(username string, status int) error {
	if err := store.DB.Model(&models.User{}).Where("username = ?", username).Update("status", status).Error; err != nil {
		logger.Errorf("failed to set user status: %v, username: %s, status: %d", err, username, status)
		return fmt.Errorf("failed to set user status: %v", err)
	}
	return nil
}

func (us *UserService) SetUserActive(username string, active bool) error {
	if err := store.DB.Model(&models.User{}).Where("username = ?", username).Update("active", active).Error; err != nil {
		logger.Errorf("failed to set user active: %v, username: %s, active: %t", err, username, active)
		return fmt.Errorf("failed to set user active: %v", err)
	}
	return nil
}

func (us *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := store.DB.First(&user, id).Error; err != nil {
		logger.Errorf("user does not exist: %d", id)
		return nil, fmt.Errorf("user does not exist: %d", id)
	}
	return &user, nil
}

// 查找一个能处理此业务的客服
func (us *UserService) FindAgent(appID string) (*models.User, error) {
	var users []models.User
	// 将 appID 转换为小写
	lowerAppID := strings.ToLower(appID)

	// 先获取所有角色为 agent、状态为在席（1）、激活状态为 true 的用户
	if err := store.DB.Where("role = ? AND status = ? AND active = ?", "agent", 1, true).Find(&users).Error; err != nil {
		logger.Errorf("failed to get agents: %v", err)
		return nil, fmt.Errorf("failed to get agents")
	}

	// 把 users 列表 顺序随机排序
	shuffle.Shuffle(users)

	// 优先查找精确匹配指定 appID 的客服
	for _, user := range users {
		// 将 Apps 字段转换为小写
		lowerApps := strings.ToLower(user.Apps)

		// 检查是否精确匹配指定的 appID
		if strings.Contains(lowerApps, fmt.Sprintf("\"%s\"", lowerAppID)) {
			return &user, nil
		}
	}

	// 如果没有找到精确匹配的客服，查找包含 "all" 的客服
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Apps), "all") {
			return &user, nil
		}
	}

	// 如果没有找到符合条件的客服，返回错误
	logger.Errorf("no available agent found for appID: %s", appID)
	return nil, fmt.Errorf("no available agent found")
}
