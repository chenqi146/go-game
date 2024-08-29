package qerr

import "fmt"

type CodeError struct {
	Code    uint32
	Message string
}

const (
	// ServerError 1开头为服务级别错误码
	ServerError        = 10001
	ParamError         = 10002
	UnsupportedError   = 10003
	DbError            = 10004
	ConvertError       = 10005
	ExistError         = 10006
	JsonUnmarshalError = 10007
	RpcInvokeError     = 10008
	
	// FileParseError 文件解析错误
	FileParseError  = 10901
	FileOpenError   = 10902
	FileRenameError = 10903
	
	// 2开头为业务错误码
	
	// GenerateTokenError 202为用户模块
	GenerateTokenError   = 20201
	TokenExpireError     = 20202
	JwtTokenParseError   = 20203
	JwtInvalidTokenError = 20204
	UserRpcInvokeError   = 20205
	
	// WxMiniAuthError 203为微信错误码
	WxMiniAuthError = 20301
	
	// AgentInvokeError 205为agent相关错误码
	AgentInvokeError = 20501
	
	// DetectLimitError 206为fire相关错误码
	DetectLimitError          = 20601
	DayShareAddFireLimitError = 20602
	DetectError               = 20603
	// RedisError 207为redis相关错误码
	RedisError = 20701
)

var errMsgMap = map[uint32]string{
	ServerError:               "服务器错误",
	ParamError:                "参数错误",
	UnsupportedError:          "不支持的方式",
	DbError:                   "Db错误",
	ConvertError:              "转换错误",
	GenerateTokenError:        "token生成错误",
	TokenExpireError:          "token过期",
	ExistError:                "已存在",
	JsonUnmarshalError:        "json解码错误",
	JwtTokenParseError:        "jwtToken解析错误",
	JwtInvalidTokenError:      "非法token",
	FileParseError:            "文件解析错误",
	FileOpenError:             "文件打开错误",
	FileRenameError:           "文件重命名错误",
	AgentInvokeError:          "代理调用错误",
	RpcInvokeError:            "rpc调用错误",
	UserRpcInvokeError:        "用户rpc调用失败",
	WxMiniAuthError:           "微信小程序授权失败",
	DetectLimitError:          "识别限制错误",
	DetectError:               "识别错误",
	DayShareAddFireLimitError: "当日分享增加次数已到达上限",
	RedisError:                "redis错误",
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}
func NewServerError() *CodeError {
	return &CodeError{
		Code:    ServerError,
		Message: errMsgMap[ServerError],
	}
}
func NewServerMessageError(message string) *CodeError {
	if len(message) == 0 {
		return &CodeError{
			Code:    ServerError,
			Message: errMsgMap[ServerError],
		}
	}
	return &CodeError{
		Code:    ServerError,
		Message: message,
	}
}

func NewErrCodeMsg(code uint32, message string) *CodeError {
	return &CodeError{Code: code, Message: message}
}

func NewErrCode(code uint32) *CodeError {
	s := errMsgMap[code]
	if len(s) > 0 {
		return &CodeError{Code: code, Message: s}
	}
	return &CodeError{Code: code, Message: ""}
}
