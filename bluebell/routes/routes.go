package routes

import (
	"bluebell/controller"
	_ "bluebell/docs" // 千万不要忘了导入把你上一步生成的docs
	"bluebell/logger"
	"bluebell/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/gin-contrib/pprof"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode{
		gin.SetMode(gin.ReleaseMode) // 设置成发布模式
	}

	r := gin.New()
	// 测试中间件， 这个是限流中间件
	//r.Use(logger.GinLogger(), logger.GinRecovery(true),middlewares.RateLimitMiddleware(2*time.Second,1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 前端
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})



	// 接口文档
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("api/v1")

	// 注册路由业务
	v1.POST("/signup", controller.SignUpHandler)
	// 登录
	v1.POST("/login", controller.LoginHandler)

	v1.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件

	// 根据时间或分数获取帖子列表
	v1.GET("/posts2", controller.GetPostListHandler2)
	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)
	v1.GET("/post/:id", controller.GetPostDetailHandler)
	v1.GET("/posts", controller.GetPostListHandler)

	{
		v1.POST("/post", controller.CreatePostHandler)

		// 投票
		v1.POST("/vote", controller.PostVoteController)
	}




	// 把判断的过程全部放到 JWTAuthMiddleware 这个中间件中
	r.GET("/ping", middlewares.JWTAuthMiddleware(),func(c *gin.Context) {
		// 如果是登录用户, 判断请求中是否有jwt 有效的token
		//c.Request.Header.Get("Authorization")
		//if isLogin{
		c.String(http.StatusOK, "pong")
		//}else{
		//	//否则返回登录页面
		//	c.String(http.StatusOK, "请登录")
		//}

	})

	pprof.Register(r) // 注册pprof相关路由

	r.NoRoute(func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}


