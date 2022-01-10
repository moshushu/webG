package controllers

import (
	"errors"
	"fmt"
	"web1/dao/mysql"
	"web1/logic"
	"web1/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SignUpHandler 注册业务
func SignUpHandler(c *gin.Context) {
	// 1、请求参数处理或校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//请求参数有误
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": err.Error(),
		// })
		return
	}
	// 2、注册业务处理
	if err := logic.SignUp(p); err != nil {
		// 记录注册错误的日志
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		// 判断err是不是与mysql.ErrorUserExist相同
		if errors.Is(err, mysql.ErrorUserExist) {
			// 用户已存在
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": "注册失败",
		// })
		return
	}
	// 3、返回响应
	ResponseSuccess(c, nil)
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": "注册成功",
	// })
}

// LoginHandler 登录业务
func LoginHandler(c *gin.Context) {
	// 1、请求参数处理或校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with  invalid param", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": err.Error(),
		// })
		return
	}
	// 2、登录业务的处理
	user, err := logic.Login(p)
	if err != nil {
		// 记录登录错误的日志
		zap.L().Error("logic.Login failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			// 用户不存在
			ResponseError(c, CodeUserNotExist)
			return
		}
		// 密码错误
		ResponseError(c, CodeInvalidPassword)
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": "用户名或密码错误",
		// })
		return
	}
	// 3、返回响应
	ResponseSuccess(c, gin.H{
		// id值大于1<<53-1  int64类型的最大值是1<<63-1
		// 当id的值超过这个范围时，就会失真
		// 解决办法：将int64类型的数值以string类型的格式传递给前端
		"user_id":   fmt.Sprintf("%d", user.UserID),
		"user_name": user.Username,
		"token":     user.Token,
	})
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": "登录成功",
	// })
}
