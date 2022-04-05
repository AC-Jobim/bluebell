package routers

import (
	"bluebell/controllers"
	"bluebell/logger"
	"bluebell/middlewares"
	"bluebell/pkg/snowflake"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //gin设置成发布模式
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//注册业务路由
	r.POST("/signup", controllers.SignUpHandler)
	r.POST("/login", controllers.LoginHandler)

	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {

		fmt.Println(controllers.GetCurrentUser(c))
		c.String(http.StatusOK, "success")
	})

	r.GET("/genid", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"雪花ID": snowflake.GenID(),
		})
	})

	return r
}
