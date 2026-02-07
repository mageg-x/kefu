package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"kefu-server/models"
	"kefu-server/utils"
	"kefu-server/utils/logger"
	"kefu-server/utils/response"
)

type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"` // 2次hash后的密码
	Timestamp string `json:"timestamp" binding:"required"`
	Nonce     string `json:"nonce" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type UserController struct{}

// hashPassword 使用SHA256 hash密码
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// verifySignature 验证签名（使用密码的1次hash作为密钥）
func verifySignature(username, passwordHash1, passwordHash2, timestamp, nonce, signature string) bool {
	// 组合数据（使用2次hash后的密码）
	data := username + passwordHash2 + timestamp + nonce
	// 使用密码的1次hash作为密钥生成HMAC
	h := hmac.New(sha256.New, []byte(passwordHash1))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	// 比较签名
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// verifyTimestamp 验证时间戳，防止回放攻击
func verifyTimestamp(timestampStr string) bool {
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false
	}
	// 时间戳有效期为5分钟
	return time.Now().Unix()-timestamp < 300
}

func (uc *UserController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("login request parameter error: %v", err)
		response.ResponseError(c, http.StatusBadRequest, response.ErrCodeInvalidParams)
		return
	}

	// 验证时间戳
	if !verifyTimestamp(req.Timestamp) {
		logger.Errorf("invalid timestamp: %s", req.Timestamp)
		response.ResponseError(c, http.StatusBadRequest, response.ErrCodeInvalidParams)
		return
	}

	var user models.User
	if err := utils.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		logger.Errorf("user does not exist: %s", req.Username)
		response.ResponseError(c, http.StatusUnauthorized, response.ErrCodeInvalidCredentials)
		return
	}

	// 对数据库中的明文密码进行1次hash，得到密码的1次hash
	passwordHash1 := hashPassword(user.Password)

	// 验证签名（使用密码的1次hash作为密钥）
	if !verifySignature(req.Username, passwordHash1, req.Password, req.Timestamp, req.Nonce, req.Signature) {
		logger.Errorf("invalid signature for user: %s", req.Username)
		response.ResponseError(c, http.StatusBadRequest, response.ErrCodeInvalidParams)
		return
	}

	// 对密码的1次hash再进行1次hash，得到密码的2次hash
	passwordHash2 := hashPassword(passwordHash1)

	// 比较2次hash后的密码
	if passwordHash2 != req.Password {
		logger.Errorf("password error: %s", req.Username)
		response.ResponseError(c, http.StatusUnauthorized, response.ErrCodeInvalidCredentials)
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		logger.Errorf("generate token failed: %v", err)
		response.ResponseError(c, http.StatusInternalServerError, response.ErrCodeInternalError)
		return
	}

	logger.Infof("user login successful: %s", req.Username)
	response.ResponseSuccess(c, LoginResponse{
		Token: token,
		User:  user,
	})
}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.ResponseError(c, http.StatusUnauthorized, response.ErrCodeUnauthorized)
		return
	}

	var user models.User
	if err := utils.DB.First(&user, userID).Error; err != nil {
		logger.Errorf("get user info failed: %v", err)
		response.ResponseError(c, http.StatusNotFound, response.ErrCodeNotFound)
		return
	}

	response.ResponseSuccess(c, gin.H{"user": user})
}

func (uc *UserController) Logout(c *gin.Context) {
	logger.Infof("user logout")
	response.ResponseSuccess(c, gin.H{"message": "logout successful"})
}
