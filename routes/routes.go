package routes

import (
	"web1/controllers"
	"web1/logger"
	"web1/middlewares"

	"github.com/gin-gonic/gin"
)

func Routes(mode string) *gin.Engine {
	// 判断是否的发布模式
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	// 不使用gin默认的中间件，使用定制logger中间件
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	v1 := r.Group("/api/v1")
	// 注册业务路由
	v1.POST("/signup", controllers.SignUpHandler)

	// 登录业务路由
	v1.POST("/login", controllers.LoginHandler)

	// 注册认证中间件
	v1.Use(middlewares.JWTAuthMiddleware())

	// 业务功能
	{
		v1.GET("/community", controllers.CommunityHandler)
		// :id 路径参数
		v1.GET("/community/:id", controllers.CommunityDetailHandler)
		// 创建帖子
		v1.POST("/post", controllers.CreatePostHandler)
		// 点击帖子，获取帖子的详细信息
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
		// 获取帖子列表
		v1.GET("/posts/",controllers.GetPostListHandler)
	}

	// 只有登录的用户才能访问，游客是不许访问的
	// middlewares.JWTAuthMiddleware()判断请求头中是否为有效的JWT？
	// 当middlewares.JWTAuthMiddleware()能通过时，即表明是已经登录的用户
	// r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	// 	c.String(http.StatusOK, "Pong")
	// })
	return r
}
