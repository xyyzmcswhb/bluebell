package controller

import (
	"strconv"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 创建帖子
func CreatePostHandler(c *gin.Context) {
	//1.获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//从c取到当前发请求的用户ID
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedAuth)
	}
	p.AuthorID = userID
	//2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

// 获取帖子详情信息
func GetPostDetailHandler(c *gin.Context) {
	//1.获取参数，从URL中获取POSTID
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据ID取出帖子数据（查数据库）
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByID(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)
}

// 获取帖子列表
func GetPostListHandler(c *gin.Context) {
	page, size := getPageInfo(c)
	//获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error(" logic.GetPostList() failed")
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// 升级版帖子列表接口，根据前端传过来的排序类型参数，动态获取帖子列表，按创建时间排序，或者按照分数排序
// GetPostsHandler 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer jwt"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts [get]

func GetPostsHandler(c *gin.Context) {
	//1.获取参数 /api/v1/posts?page=1&... Querystring参数
	//2.去redis查询id列表
	//3.根据id去数据库查询帖子详细信息
	//初始化结构体时制定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.Orderbytime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostsHandler with invalid params",
			zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetPostListNew(p) //合二为一
	p.Page, p.Size = getPageInfo(c)
	//获取数据
	if err != nil {
		zap.L().Error(" logic.GetPostList() failed")
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// 根据社区去查询帖子列表
//func GetCommunityPostListHandler(c *gin.Context) {
//	p := &models.ParamCommunityPostList{
//		ParamPostList: &models.ParamPostList{
//			Page:  1,
//			Size:  10,
//			Order: models.Orderbytime,
//		},
//	}
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetCommunityPostListHandler with invalid params",
//			zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//	//获取数据
//
//	if err != nil {
//		zap.L().Error(" logic.GetPostList() failed")
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//	ResponseSuccess(c, data)
//}
