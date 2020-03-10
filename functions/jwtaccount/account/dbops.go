package account

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	helpers "github.com/sawima/aws-sls-jwt-service/functions/layers/helpers"
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

//CheckDefaultAppAccount init the default account
// func CheckDefaultAppAccount() error {
// 	result, err := dbclient.GetItem(&dynamodb.GetItemInput{
// 		TableName: aws.String(defaultTableName),
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"appid": {
// 				S: aws.String(defaultAppid),
// 			},
// 		},
// 	})
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return err
// 	}
// 	// fmt.Printf("%v", result.Item)

// 	tapp := models.App{}
// 	err = dynamodbattribute.UnmarshalMap(result.Item, &tapp)

// 	if err != nil {
// 		return err
// 	}

// 	if tapp.Appid == "" {
// 		newItem := models.App{
// 			Appid:     defaultAppid,
// 			Hashedkey: helpers.GenerateHashPassword(defaultPasswd),
// 			Appname:   defaultAppid,
// 			Context: models.Appcontext{
// 				Org:   defaultAppid,
// 				Orgid: defaultAppid,
// 			},
// 		}
// 		newDefaultAppItem, _ := dynamodbattribute.MarshalMap(newItem)

// 		_, err = dbclient.PutItem(&dynamodb.PutItemInput{
// 			TableName: aws.String(defaultTableName),
// 			Item:      newDefaultAppItem,
// 		})
// 		if err != nil {
// 			log.Println(err.Error())
// 			return err
// 		}
// 		log.Println("default app created")
// 	}

// 	return nil
// }

//UpdateDefaultSecurityKey receive post request to reset default password of default appid
func UpdateDefaultSecurityKey() (bool, string, error) {
	randPwd, _ := helpers.GenerateRandomString(20)
	hashedkey := helpers.GenerateHashPassword(randPwd)
	_, err := dbclient.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pwd": {
				S: aws.String(hashedkey),
			},
		},
		TableName: aws.String(defaultTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"appid": {
				S: aws.String(defaultAppid),
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
	//todo:export to http service,post request
	return true, randPwd, nil
}
