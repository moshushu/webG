package controllers

import (
	"strconv"
	"web1/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 给社区分类列表的功能
func CommunityHandler(c *gin.Context) {
	// 查找到所有的社区（community_id,community_name）以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.CommunityList() failed", zap.Error(err))
		// 不要轻易把服务器端报错暴露给外面
		ResponseError(c, CodeServerBusy)
		return
	}
	// 成功
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 根据id查找社区分类详情的功能
func CommunityDetailHandler(c *gin.Context) {
	// 1、获取社区id(参数处理),从url中获取
	communityID := c.Param("id")
	// communityID是string类型，要将其装换成int类型，10代表十进制，64代表64位
	// strconv.ParseInt 将字符串解析成数字
	id, err := strconv.ParseInt(communityID, 10, 64)
	if err != nil {
		// 参数出错
		zap.L().Error("Get CommunityDetailHandler id with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2、根据id获取社区的详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail faild", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
