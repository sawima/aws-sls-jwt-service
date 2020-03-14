package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (Response, error) {
	log.Println("***********auth**************")
	log.Println(req.RequestContext.Authorizer)
	log.Printf("%v", req.RequestContext.Authorizer)
	authContext := req.RequestContext.Authorizer
	// {
	// 	Appname:
	// 	indicatevalue:
	// 	uuid

	// 	integrationLatency:
	// 	principalId:
	// }

	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message":       "api demo!",
		"uuid":          authContext["uuid"],
		"appname":       authContext["appname"],
		"indicatevalue": authContext["indicatevalue"],
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
