package middlewares

import (
	"strings"
	"web1/controllers"
	"web1/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 形式如：Authorization: Bearer xxx.xxxx.xxx
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controllers.ResponseError(c, controllers.CodeNeedLogin)
			// 不再执行后面的中间件
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		// controllers.ContextUserIDKey,封装在了controllers的request中，避免包循环导入的问题
		c.Set(controllers.ContextUserIDKey, mc.UserID)
		// 当能走到这一步的时候，说明token是没有问题的，也代表着登录的用户是已经认证的用户
		c.Next() // 后续的处理请求的函数中可以用过c.Get(ContextUserIDKey)来获取当前请求的用户信息
	}
}
