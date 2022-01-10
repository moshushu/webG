package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code":10000, //程序中的错误码
	"msg": xxx, // 提示信息
	"data": {}, // 数据
}
*/

// 错误码，封装响应方法，方便前端对接
type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// ResponseError 错误响应，返回错误信息
// code————错误码
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

// ResoponseErrorWithMsg 错误响应，返回错误信息，可自定义msg错误提示
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// RespoenseSuccess 成功响应，返回成功信息，不需要传错误码
// Code———— 固定为CodeSuccess
// Data———— 需要传入成功响应的数据，供返回
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

// func ResponseSuccessWithMsg(c *gin.Context, msg interface{}, data interface{}) {
// 	c.JSON(http.StatusOK, &ResponseData{
// 		Code: CodeSuccess,
// 		Msg:  msg,
// 		Data: data,
// 	})
// }
