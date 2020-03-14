package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	helpers "github.com/sawima/aws-sls-jwt-service/functions/layers/helpers"
	models "github.com/sawima/aws-sls-jwt-service/functions/layers/models"
)

// var dbclient *dynamodb.DynamoDB
var defaultAppid = "kimatech"
var defaultPasswd = "kimapasswd"
var defaultTableName = "users"

//checkDefaultAppAccount init the default account
func checkDefaultAppAccount() error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dbclient := dynamodb.New(sess)
	result, err := dbclient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(defaultTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"appid": {
				S: aws.String(defaultAppid),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	// fmt.Printf("%v", result.Item)

	tapp := models.App{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &tapp)

	if err != nil {
		return err
	}

	if tapp.Appid == "" {
		newItem := models.App{
			Appid:     defaultAppid,
			Hashedkey: helpers.GenerateHashPassword(defaultPasswd),
			Context: models.Appcontext{
				Appname:       defaultAppid,
				UUID:          defaultAppid,
				Indicatevalue: defaultAppid,
			},
		}
		newDefaultAppItem, _ := dynamodbattribute.MarshalMap(newItem)

		_, err = dbclient.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(defaultTableName),
			Item:      newDefaultAppItem,
		})
		if err != nil {
			log.Println(err.Error())
			return err
		}
		log.Println("default app created")
	}

	return nil
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	err := checkDefaultAppAccount()
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	var buf bytes.Buffer

	body, _ := json.Marshal(map[string]interface{}{
		"message": "Default application initialized",
	})
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
