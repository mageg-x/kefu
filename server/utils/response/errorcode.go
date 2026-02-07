package response

// ErrorCode 定义错误码类型
type ErrorCode int

// 通用错误码
const (
	ErrCodeSuccess            ErrorCode = 0    // 成功
	ErrCodeInvalidParams      ErrorCode = 1001 // 通用错误
	ErrCodeUnauthorized       ErrorCode = 1002
	ErrCodeForbidden          ErrorCode = 1003
	ErrCodeNotFound           ErrorCode = 1004
	ErrCodeInternalError      ErrorCode = 1005
	ErrCodeInvalidCredentials ErrorCode = 2001 // 登录相关错误
	ErrCodeTokenExpired       ErrorCode = 2002
	ErrCodeTokenInvalid       ErrorCode = 2003
)

// ErrorMessages 错误码到错误消息的映射
var ErrorMessages = map[ErrorCode]string{

	ErrCodeSuccess:            "success",            // 成功
	ErrCodeInvalidParams:      "invalid parameters", // 通用错误
	ErrCodeUnauthorized:       "unauthorized",
	ErrCodeForbidden:          "forbidden",
	ErrCodeNotFound:           "resource not found",
	ErrCodeInternalError:      "internal server error",
	ErrCodeInvalidCredentials: "invalid username or password", // 登录相关错误
	ErrCodeTokenExpired:       "token expired",
	ErrCodeTokenInvalid:       "invalid token",
}
