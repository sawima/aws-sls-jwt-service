package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jwtaccount "github.com/sawima/aws-sls-jwt-service/functions/jwtaccount/account"
	models "github.com/sawima/aws-sls-jwt-service/functions/layers/models"
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

//ResetAppAccount recieve post request to reset none-default app password and return password to front
func ResetAppAccount(c *gin.Context) {
	var app models.App
	c.ShouldBindJSON(&app)
	if app.Appid != "" {
		success, newpwd, err := jwtaccount.UpdateTargetAppSecurityKey(app.Appid)
		if err != nil {
			c.Status(http.StatusInternalServerError)
		}
		if newpwd != "" {
			c.JSON(http.StatusOK, gin.H{"result": success, "newpassword": newpwd})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "application is not exists"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"result": "failed", "message": "wrong appid"})
	}
}

//RegisterNewApp register new app account
// {
// 	"appname":"",
// 	"context":{
// 		"org":"",
// 		"orgid":""
// 	}
// }
func RegisterNewApp(c *gin.Context) {
	app := models.App{}
	c.ShouldBindJSON(&app)

	success, newappid, newpwd, err := jwtaccount.AddNewItemInAccountTable(&app)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, gin.H{"result": success, "newpassword": newpwd, "newappid": newappid})
}
