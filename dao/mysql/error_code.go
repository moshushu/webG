package mysql

import "errors"

// 存放mysql需要的自定义错误

// 为了能在controller层进行判断，使用errors.Is()————1.13版本以后才有的
var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	ErrorInvalidID       = errors.New("无效的ID")
)
