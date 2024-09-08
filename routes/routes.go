package routes

import (
	"time"
	_ "web_app/docs" // 千万不要忘了导入把你上一步生成的docs

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"net/http"
	"web_app/controller"
	"web_app/logger"
	"web_app/midware"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	//注册业务逻辑路由
	v1.POST("/signup", controller.SignupHandler)

	v1.POST("/login", controller.LoginHandler)
	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)
	v1.GET("/post/:id", controller.GetPostDetailHandler)
	v1.GET("/postlist", controller.GetPostListHandler)
	//根据时间或分数获取帖子列表
	v1.GET("/posts", controller.GetPostsHandler)

	v1.Use(midware.JWTAuthMiddleware(), midware.RateLimitMiddleware(2*time.Second, 1)) //应用JWT认证中间件

	{
		//v1.POST("/comment", controller.CommentHandler)    // 评论
		//v1.GET("/comment", controller.CommentListHandler) // 评论列表
		v1.POST("/post", controller.CreatePostHandler)

		//投票
		v1.POST("/vote", controller.PostVoteController)
		v1.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	}

	// v1.GET("/", midware.JWTAuthMiddleware(), func(c *gin.Context) {

	// 	//如果是登录的用户，判断请求头中是否包含有效的JWT
	// 	c.Request.Header.Get("Authorization")
	// 	//如果是登录的用户
	// 	c.String(http.StatusOK, "pong")
	// 	//否则直接返回请登录
	// 	c.String(http.StatusOK, "请登录")
	// })
	pprof.Register(r) ///注册pprof相关路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": 404,
		})
	})

	return r
}
