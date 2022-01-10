package controllers

// ResCode 给错误码定一个类型
type ResCode int64

// 错误码常量
const (
	CodeSuccess         ResCode = 1000 + iota // success
	CodeFailure                               // failure
	CodeInvalidParam                          // 请求参数错误
	CodeUserExist                             // 用户已存在
	CodeUserNotExist                          // 用户不存在
	CodeInvalidPassword                       // 用户名或密码错误
	CodeServerBusy                            // 服务繁忙
	CodeNeedLogin                             // 需要登录
	CodeInvalidToken                          // 无效的Token
)

// codeMsgMap 创建一个code提示信息的map
var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeFailure:         "failure",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",
	CodeNeedLogin:       "需要登录",
	CodeInvalidToken:    "无效的Token",
}

// 给ResCode建立一个方法
// 当错误码调用Msg时，会返回对应的提示信息
func (c ResCode) Msg() string {
	// 判断是否存在这个错误码
	msg, ok := codeMsgMap[c]
	if !ok {
		// 如果不存在则返回服务繁忙
		msg = codeMsgMap[CodeServerBusy]
	}
	// 如果存在则返回对应key的值
	return msg
}
