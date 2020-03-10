package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/dynamodb"
// 	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
// 	helpers "github.com/sawima/aws-sls-jwt-service/functions/layers/helpers"
// 	models "github.com/sawima/aws-sls-jwt-service/functions/layers/models"
// )

// // var dbclient *dynamodb.DynamoDB
// var defaultAppid = "kimatech"
// var defaultPasswd = "kimapasswd"
// var defaultTableName = "users"

// // func init() {
// // 	dbclient = dynamodbClient()
// // }

// // func dynamodbClient() *dynamodb.DynamoDB {
// // 	sess := session.Must(session.NewSessionWithOptions(session.Options{
// // 		SharedConfigState: session.SharedConfigEnable,
// // 	}))
// // 	svc := dynamodb.New(sess)
// // 	return svc
// // }

// //CheckDefaultAppAccount init the default account
// func checkDefaultAppAccount() error {
// 	sess := session.Must(session.NewSessionWithOptions(session.Options{
// 		SharedConfigState: session.SharedConfigEnable,
// 	}))
// 	dbclient := dynamodb.New(sess)
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
