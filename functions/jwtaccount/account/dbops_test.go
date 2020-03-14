package account

import (
	"testing"

	models "github.com/sawima/aws-sls-jwt-service/functions/layers/models"
)

// func TestInitDefaultAccount(t *testing.T) {
// 	err := CheckDefaultAppAccount()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log("success init default account")
// }

func TestUpdateDefaultHashedPasswd(t *testing.T) {
	_, passwd, err := UpdateDefaultSecurityKey()
	if err != nil {
		t.Error(err)
	}
	t.Log("success update default password", passwd)
}

func TestAddNewItemToDynamodb(t *testing.T) {
	app := models.App{
		Appname: "testapp2",
		Context: models.Appcontext{
			Org:   "kimacloud",
			Orgid: "kimacloudapp2",
		},
	}
	_, appid, passwd, err := AddNewItemInAccountTable(&app)

	if err != nil {
		t.Error(err)
	}
	t.Log("saved new app id", appid)
	t.Log("saved new app password, ", passwd)
}
