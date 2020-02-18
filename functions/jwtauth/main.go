package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jwt "github.com/dgrijalva/jwt-go"

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/ObjectId"

	// "github.com/mongodb/mongo-go-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"fmt"
)

//APIResponse aws api gw proxy response
type APIResponse events.APIGatewayProxyResponse

//APIRequest aws api gw proxy request
type APIRequest events.APIGatewayProxyRequest

type config struct {
	PemFilePath      string
	SimpleDefaultPem string
	ConnStr          string
	DBDriver         string
}

var serverconfig = &config{
	PemFilePath:      "/.resource/jwtapp.rsa",
	SimpleDefaultPem: "thisisthesimplesecretkey",
	// ConnStr:          "mongodb://authuser:foodunion2019@dds-uf6489fcd43470841611-pub.mongodb.rds.aliyuncs.com:3717,dds-uf6489fcd43470842606-pub.mongodb.rds.aliyuncs.com:3717/platform?replicaSet=mgset-14040555",
	ConnStr:  "mongodb://authuser:foodunion2019@dds-uf6489fcd43470841611-pub.mongodb.rds.aliyuncs.com:3717,dds-uf6489fcd43470842606-pub.mongodb.rds.aliyuncs.com:3717/platform?replicaSet=mgset-14040555",
	DBDriver: "mongo",
}

type accessToken struct {
	Token string
	Scope string
}

type verifyInfo struct {
	IsValid bool
}

type myCustomClaims struct {
	Account string `json:"account"`
	AppName string `json:"appname"`
	jwt.StandardClaims
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

//App model definition
// type App struct {
// 	// ID             objectid.ObjectID `json:"id" bson:"_id"`
// 	Account        string `json:"account" bson:"account,omitempty"`
// 	HashedPassword string `json:"password" bson:"hashedpassword,omitempty"`
// 	AppName        string `json:"appname" bson:"appname,omitempty"`
// }

//App model definition
type App struct {
	Appid       string `json:"appid"`
	Securitykey string `json:"securitykey"`
	Hashedkey   string `json:"hashedkey"`
	Appname     string `json:"appname"`
}

func dbAuth(appid, securitykey string) (app *App, authIndex bool) {

	log.Println("appid", "securitykey")
	log.Println(appid, securitykey)
	// snippet-start:[dynamodb.go.read_item.session]
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	log.Println("dynamodb connection status")
	log.Println(svc)
	// snippet-end:[dynamodb.go.read_item.session]

	// snippet-start:[dynamodb.go.read_item.call]
	tableName := "applications"
	// movieName := "The Big New Movie"
	// movieYear := "2015"

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

	fmt.Printf("%v", result.Item)
	// snippet-end:[dynamodb.go.read_item.call]

	// snippet-start:[dynamodb.go.read_item.unmarshall]
	tapp := App{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &tapp)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	fmt.Printf("%v", tapp)

	if tapp.Securitykey == "" {
		fmt.Println("Could not find target app")
		return nil, false
	}

	if tapp.Securitykey == securitykey {
		return &tapp, true
	}

	return nil, false
}

func verify(account, passwd string) (*App, bool) {
	return dbAuth(account, passwd)
	// if account != "" && passwd != "" {
	// 	client, err := mongo.Connect(context.TODO(), serverconfig.ConnStr)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	collection := client.Database("platform").Collection("applications")
	// 	app := &App{}
	// 	err = collection.FindOne(context.Background(), bson.D{{"account", account}}).Decode(app)
	// 	// err = collection.FindOne(context.Background(), bson.D{}).Decode(app)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	// if checkPasswordHash(app.HashedPassword, passwd) {
	// 	// 	return app, true
	// 	// }
	// 	if app != nil {
	// 		if app.HashedPassword == passwd {
	// 			return app, true
	// 		}
	// 	} else {
	// 		log.Println("app is not found")
	// 		return nil, false
	// 	}
	// }

	// return nil, false
}

// func HashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return string(bytes), err
// }

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//UserFetchToken RPC mothod for user request new access token
func userFetchToken(ctx context.Context, in userAuthRequest) (ReturnToken, error) {
	// pemFile := getPemFile()
	log.Println("start fetch token")
	pemFile := []byte("thisisthefoodunionencrptstring")
	log.Println("input value")
	log.Println(in.Account)
	log.Println(in.Passwd)
	if user, ok := verify(in.Account, in.Passwd); ok {
		claims := myCustomClaims{
			Account: user.Appid,
			AppName: user.Appname,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 2000).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		accessToken, err := token.SignedString(pemFile)
		if err != nil {
			log.Println(err)
			return ReturnToken{Token: "", Success: false}, err
		}
		log.Println("success generate token")
		return ReturnToken{Token: accessToken, Success: true}, nil
	}

	log.Println("failed ,not auth")
	return ReturnToken{Token: "", Success: false}, errors.New("Not authorized")
}

//LambdaGenerateToken main entry of auth lambda funciton
func LambdaGenerateToken(ctx context.Context, request APIRequest) (APIResponse, error) {
	authRequest := userAuthRequest{}
	json.Unmarshal([]byte(request.Body), &authRequest)
	token, err := userFetchToken(context.Background(), authRequest)
	if err != nil {
		log.Println(err)
	}
	log.Println("function run to end")
	tokenJSONStr, _ := json.Marshal(token)
	resp := APIResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(tokenJSONStr),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "foodunion-jwt-auth",
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(LambdaGenerateToken)
}
