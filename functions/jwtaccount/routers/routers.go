/*
Package routers all site routers are here
*/
package routers

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	controllers "github.com/sawima/aws-sls-jwt-service/functions/jwtaccount/controllers"
)

//SetupRouters export all route to server using gin Router
func SetupRouters() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	dnsStageRoute := r.Group("/auth/account")
	dnsStageRoute.GET("", func(c *gin.Context) { c.String(http.StatusOK, "go aws wildcard lambda") })

	// dnsStageRoute.POST("/", controllers.InitDefaultAppAccount)
	dnsStageRoute.POST("/update", controllers.ResetDefaultAppAccount)
	dnsStageRoute.POST("/updateapp", controllers.ResetAppAccount)
	dnsStageRoute.POST("/new", controllers.RegisterNewApp)

	return r
}
