package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	myrouters "github.com/sawima/aws-sls-jwt-service/functions/api-template/routers"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// stdout and stderr are sent to AWS CloudWatch Logs
	router := myrouters.SetupRouters()

	ginLambda = ginadapter.New(router)
}

//Handler is the export func for lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}

// import (
// 	myrouters "github.com/sawima/aws-sls-gin-template/functions/api-template/routers"
// )

// func main() {
// 	router := myrouters.SetupRouters()
// 	router.Run(":8090") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
// }
