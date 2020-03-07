package main

import (
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jwt "github.com/dgrijalva/jwt-go"
)

type myCustomClaims struct {
	Account string `json:"account"`
	AppName string `json:"appname"`
	jwt.StandardClaims
}

func handler(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.AuthorizationToken
	tokenSlice := strings.Split(token, " ")
	var tokenStr string
	if len(tokenSlice) > 1 {
		tokenStr = tokenSlice[len(tokenSlice)-1]
	}
	pemFile := []byte("thisisthekimatechtokenstring")
	claims := myCustomClaims{}
	if tokenStr != "" {
		log.Println(tokenStr)
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", token.Header["alg"])
			}
			return pemFile, nil
		})
		if err != nil {
			log.Println(err)
			return generatePolicy("user", "Deny", request.MethodArn), nil
		}
		log.Printf("%v", claims)
		if token.Valid {
			apiandpath := strings.Split(request.MethodArn, "/")
			resource := strings.Join(apiandpath[0:1], "/")
			resource += "/*"
			return generatePolicy("user", "Allow", resource), nil
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return generatePolicy("user", "Deny", request.MethodArn), nil
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				log.Println("Expired token")
				return generatePolicy("user", "Deny", request.MethodArn), nil
			} else {
				log.Println("Expired token,any where")
				return generatePolicy("user", "Deny", request.MethodArn), nil
			}
		} else {
			log.Println("Expired token")
			return generatePolicy("user", "Deny", request.MethodArn), nil
		}
	} else {
		log.Println("Token is empty")
		return generatePolicy("user", "deny", request.MethodArn), nil

	}
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}
	return authResponse
}

func main() {
	lambda.Start(handler)
}
