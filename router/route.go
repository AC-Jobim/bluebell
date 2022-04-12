package router

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"bluebell/pkg/snowflake"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //gin设置成发布模式
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("api/v1")

	//注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	// 登陆
	v1.POST("/login", controller.LoginHandler)

	//v1使用中间件,使用中间件的方法的集合...
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)

		// 查询帖子升级版
		v1.GET("/posts/", controller.GetPostListHandler2)

		// 投票功能实现
		v1.POST("/vote", controller.PostVoteController)
	}

	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {

		c.String(http.StatusOK, "success")
	})

	r.GET("/genid", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"雪花ID": snowflake.GenID(),
		})
	})

	return r
}
