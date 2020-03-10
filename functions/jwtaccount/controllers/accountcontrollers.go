package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jwtaccount "github.com/sawima/aws-sls-jwt-service/functions/jwtaccount/account"
)

//InitDefaultAppAccount receive gin post request to init default app account
// func InitDefaultAppAccount(c *gin.Context) {
// 	err := jwtaccount.CheckDefaultAppAccount()
// 	if err != nil {
// 		c.Status(http.StatusInternalServerError)
// 	}
// 	c.JSON(http.StatusOK, gin.H{"result": "success init default account"})
// }

//ResetDefaultAppAccount recieve post request to reset default app password and return password to front
func ResetDefaultAppAccount(c *gin.Context) {
	success, newpwd, err := jwtaccount.UpdateDefaultSecurityKey()
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, gin.H{"result": success, "newpassword": newpwd})
}
