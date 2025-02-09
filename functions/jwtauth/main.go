package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/sawima/aws-sls-jwt-service/functions/layers/helpers"

	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	jwt "github.com/dgrijalva/jwt-go"

	models "github.com/sawima/aws-sls-jwt-service/functions/layers/models"
)

//APIResponse aws api gw proxy response
type APIResponse events.APIGatewayProxyResponse

//APIRequest aws api gw proxy request
type APIRequest events.APIGatewayProxyRequest
type accessToken struct {
	Token string
	Scope string
}

//ReturnToken return token to the request
type ReturnToken struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

type userAuthRequest struct {
	Account string `json:"account"`
	Passwd  string `json:"password"`
}

// //App model definition
// type App struct {
// 	Appid     string     `json:"appid"`
// 	Hashedkey string     `json:"hashedkey"`
// 	Context   Appcontext `json:"context"`
// }

// //Appcontext is the context info for applications
// type Appcontext struct {
// 	Appname   string     `json:"appname"`
// 	UUID   string `json:"uuid"`
// 	Indicatekey string `json:"indicatekey"`
// }

func dbAuth(appid, password string) (app *models.App, authIndex bool) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	tableName := "users"
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"appid": {
				S: aws.String(appid),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	tapp := models.App{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &tapp)
	if err != nil {
		log.Printf("Failed to unmarshal Record, %v", err)
	}
	if tapp.Hashedkey == "" {
		log.Println("Could not find target app")
		return nil, false
	}
	if helpers.CheckPasswordHash(password, tapp.Hashedkey) {
		return &tapp, true
	}
	return nil, false
}

func verify(account, passwd string) (*models.App, bool) {
	return dbAuth(account, passwd)
}

//UserFetchToken RPC mothod for user request new access token
// {
//     "account":"kimatech",
//     "password":"6UN6WwPhhgEz1oASmvl_0Vo3lIU="
// }
func userFetchToken(ctx context.Context, in userAuthRequest) (ReturnToken, int, error) {
	// pemFile := getPemFile()
	pemFile := []byte("thisisthekimatechtokenstring")
	if app, ok := verify(in.Account, in.Passwd); ok {
		claims := models.MyCustomClaims{
			Appid:   app.Appid,
			Context: app.Context,
			StandardClaims: jwt.StandardClaims{
				// ExpiresAt: time.Now().Add(time.Hour * 2000).Unix(),
				ExpiresAt: time.Now().Add(time.Hour * 1800).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		accessToken, err := token.SignedString(pemFile)
		if err != nil {
			log.Println(err)
			return ReturnToken{Token: "", Success: false}, 500, err
		}
		log.Println("success generate token")
		return ReturnToken{Token: accessToken, Success: true}, 200, nil
	}

	log.Println("failed ,not auth")
	return ReturnToken{Token: "user name or password is not correct", Success: false}, 403, errors.New("Not authorized")
}

//LambdaGenerateToken main entry of auth lambda funciton
func LambdaGenerateToken(ctx context.Context, request APIRequest) (APIResponse, error) {
	authRequest := userAuthRequest{}
	json.Unmarshal([]byte(request.Body), &authRequest)
	token, statusCode, err := userFetchToken(context.Background(), authRequest)
	if err != nil {
		log.Println(err)
	}
	tokenJSONStr, _ := json.Marshal(token)
	resp := APIResponse{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            string(tokenJSONStr),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "kimatech-jwt-auth",
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(LambdaGenerateToken)
}
