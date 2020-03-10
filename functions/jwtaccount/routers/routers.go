/*
Package routers all site routers are here
*/
package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	controllers "github.com/sawima/aws-sls-jwt-service/functions/jwtaccount/controllers"
)

//SetupRouters export all route to server using gin Router
func SetupRouters() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	dnsStageRoute := r.Group("/auth")
	// dnsStageRoute.POST("/", controllers.InitDefaultAppAccount)
	dnsStageRoute.POST("/update", controllers.ResetDefaultAppAccount)

	return r
}
