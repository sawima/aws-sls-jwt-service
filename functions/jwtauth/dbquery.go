package main

// // snippet-start:[dynamodb.go.read_item.imports]
// import (
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/dynamodb"
// 	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

// 	"fmt"
// )

// // snippet-end:[dynamodb.go.read_item.imports]

// // App is the dynanodb table
// // snippet-start:[dynamodb.go.read_item.struct]
// // Create struct to hold info about new item
// // type App struct {
// // 	Appid       string `json:"appid"`
// // 	Securitykey string `json:"securitykey"`
// // 	Hashedkey   string `json:"hashedkey"`
// // 	Appname     string `json:"appname"`
// // }

// // snippet-end:[dynamodb.go.read_item.struct]

// func dbAuth(appid, securitykey string) (app *App, authIndex bool) {
// 	// snippet-start:[dynamodb.go.read_item.session]
// 	// Initialize a session that the SDK will use to load
// 	// credentials from the shared credentials file ~/.aws/credentials
// 	// and region from the shared configuration file ~/.aws/config.
// 	sess := session.Must(session.NewSessionWithOptions(session.Options{
// 		SharedConfigState: session.SharedConfigEnable,
// 	}))

// 	// Create DynamoDB client
// 	svc := dynamodb.New(sess)
// 	// snippet-end:[dynamodb.go.read_item.session]

// 	// snippet-start:[dynamodb.go.read_item.call]
// 	tableName := "applications"
// 	// movieName := "The Big New Movie"
// 	// movieYear := "2015"

// 	result, err := svc.GetItem(&dynamodb.GetItemInput{
// 		TableName: aws.String(tableName),
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"appid": {
// 				S: aws.String(appid),
// 			},
// 		},
// 	})
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil, false
// 	}

// 	fmt.Printf("%v", result.Item)
// 	// snippet-end:[dynamodb.go.read_item.call]

// 	// snippet-start:[dynamodb.go.read_item.unmarshall]
// 	tapp := App{}

// 	err = dynamodbattribute.UnmarshalMap(result.Item, &tapp)

// 	if err != nil {
// 		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
// 	}

// 	fmt.Printf("%v", tapp)

// 	if tapp.Securitykey == "" {
// 		fmt.Println("Could not find target app")
// 		return nil, false
// 	}

// 	if tapp.Securitykey == securitykey {
// 		return &tapp, true
// 	}

// 	return nil, false
// }
