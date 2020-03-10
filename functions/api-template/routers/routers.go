/*
Package routers all site routers are here
*/
package routers

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//SetupRouters export all route to server using gin Router
func SetupRouters() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	//CORS config
	r.Use(cors.Default())

	// r.GET("/gin/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	r.GET("/dev/gin", func(c *gin.Context) { c.String(http.StatusOK, "go aws wildcard lambda") })
	r.GET("/dev/gin/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	r.GET("/dev/gin/pong", func(c *gin.Context) { c.String(http.StatusOK, "world") })
	r.POST("/dev/gin/push", func(c *gin.Context) { c.String(http.StatusOK, "post data") })
	r.POST("/dev/gin/test", func(c *gin.Context) { c.String(http.StatusOK, "post data") })
	// r.GET("/gin/hello", func(c *gin.Context) { c.String(http.StatusOK, "world") })
	// r.POST("/auth", middleware.FetchAccessTokenFromRPCServer)
	// r.GET("/", middleware.AccessTokenVerifyRPCRequest(), controllers.TestMiddleware)
	// r.GET("/", controllers.TestMiddleware)
	//User related route

	// user := controllers.UserInstance()
	// r.GET("/alluser", user.All)
	// r.POST("/query/email", user.ByEmail)

	// userRouter := r.Group("/user")
	// // userRouter.Use(middleware.AccessTokenVerifyRPCRequest())
	// userRouter.POST("/", user.New)
	// userRouter.GET("/", user.All)
	// userRouter.POST("/email", user.ByEmail)
	// userRouter.PUT("/", user.Update)
	// userRouter.DELETE("/", user.Remove)

	return r
}
