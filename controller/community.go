package controller

import (
	"strconv"
	"web_app/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 社区
func CommunityHandler(c *gin.Context) {
	//查询到所有的社区(community_id, community_name)以列表的权限返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端报错暴露给外部
		return
	}
	ResponseSuccess(c, data)
}

// 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	//1. 获取社区ID
	idStr := c.Param("id")                     //获取URL参数
	id, err := strconv.ParseInt(idStr, 10, 64) //字符转换
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2. 根据id获取社区详情
	detaildata, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务端报错暴露给外部
		return
	}
	ResponseSuccess(c, detaildata)
}
