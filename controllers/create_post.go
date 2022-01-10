package controllers

import (
	"strconv"
	"web1/logic"
	"web1/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 发表功能(创建帖子)
func CreatePostHandler(c *gin.Context) {
	// 1、获取参数及参数的校验
	// 1.1、创建一个模型来存放请求中的参数
	p := new(models.CreatePost)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从c 中获取登录用户的id
	userID, err := GetCurrenUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2、创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3、返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情的处理函数
func GetPostDetailHandler(c *gin.Context) {
	// 1、获取参数（从url中获取帖子的id）
	pidStr := c.Param("id")
	id, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get GetPostDetailHandler id with invalid param", zap.Error(err))
		// 参数出错
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2、根据id取出帖子详情数据（查数据库）
	data, err := logic.GetPostById(id)
	if err != nil {
		zap.L().Error("logic.GetPostById failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3、返回
	// 改造返回的data来保证返回自己所需要的
	ResponseSuccess(c, data)
}

// GetPostListHangler 获取帖子列表的处理函数
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := GetPageInfo(c)
	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() fialed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}