package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

// GetCurrenUser 获取当前登录用户的ID，因为在middlewares中，通过c.Set()将userID保存到上下文中
func GetCurrenUser(c *gin.Context) (userID int64, err error) {
	// 获取userID
	uid, ok := c.Get(ContextUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	// 将userID转成int64类型
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

// GetPageInfo 帖子分页展示功能
func GetPageInfo(c *gin.Context) (int64, int64) {
	// 如:http://127.0.0.1:8081/api/v1/posts/?page=1&size=2
	// 获取请求的url中的page和size参数的内容
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	var (
		page int64
		size int64
		err  error
	)
	// 将pageStr转化成int64，十进制
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	// 将sizeStr转化成int64，十进制
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	// 当url中没用参数时，会返回page=1，size=10
	return page, size
}
