package account

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	helpers "github.com/sawima/aws-sls-jwt-service/functions/layers/helpers"
	models "github.com/sawima/aws-sls-jwt-service/functions/layers/models"
)

var dbclient *dynamodb.DynamoDB
var defaultAppid = "kimatech"
var defaultPasswd = "kimapasswd"
var defaultTableName = "users"

func init() {
	dbclient = dynamodbClient()
}

func dynamodbClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)
	return svc
}

//UpdateDefaultSecurityKey receive post request to reset default password of default appid
func UpdateDefaultSecurityKey() (bool, string, error) {
	appid := defaultAppid
	return updateAppAccountTable(appid)
}

func updateAppAccountTable(appid string) (bool, string, error) {
	randPwd, _ := helpers.GenerateRandomString(20)
	hashedkey := helpers.GenerateHashPassword(randPwd)
	_, err := dbclient.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pwd": {
				S: aws.String(hashedkey),
			},
		},
		ConditionExpression: aws.String("attribute_exists(hashedkey)"),
		TableName:           aws.String(defaultTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"appid": {
				S: aws.String(appid),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set hashedkey=:pwd"),
	})

	if err != nil {
		log.Println(err.Error())
		return false, "", err
	}
	log.Println("update password")
	return true, randPwd, nil
}

//UpdateTargetAppSecurityKey update target app security key
func UpdateTargetAppSecurityKey(appid string) (bool, string, error) {
	return updateAppAccountTable(appid)
}

//AddNewItemInAccountTable add new app to dynamodb
func AddNewItemInAccountTable(app *models.App) (bool, string, string, error) {
	newppid := helpers.GenerateRandAppID(20)
	app.Appid = newppid
	log.Println(app.Appid)
	apppasswd, _ := helpers.GenerateRandomString(20)
	app.Hashedkey = helpers.GenerateHashPassword(apppasswd)

	newapp, err := dynamodbattribute.MarshalMap(app)

	if err != nil {
		log.Println("unable to mashal new app")
		return false, "", "", err
	}
	_, err = dbclient.PutItem(&dynamodb.PutItemInput{
		Item:      newapp,
		TableName: aws.String(defaultTableName),
	})

	if err != nil {
		log.Println("unable save new app to dynamodb")
		return false, "", "", err
	}

	return true, newppid, apppasswd, nil
}
