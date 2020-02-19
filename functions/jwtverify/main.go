package main

import (
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jwt "github.com/dgrijalva/jwt-go"
)

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
	Token string `json:"token"`
}

type verifyInfo struct {
	IsValid bool   `json:"isvalid"`
	Account string `json:"account"`
	AppName string `json:"appname"`
}

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
		log.Println("validate any where")
		if err != nil {
			log.Println("(()(*)(*&*(*&*(*&")
			log.Println(err)
			log.Println("error printout and return custome response")
			return generatePolicy("user", "Deny", request.MethodArn), nil
		}
		log.Printf("%v", claims)
		if token.Valid {
			apiandpath := strings.Split(request.MethodArn, "/")
			resource := strings.Join(apiandpath[0:1], "/")
			resource += "/*"
			return generatePolicy("user", "Allow", resource), nil
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			log.Println("validation error**********")
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				log.Println("That's not even a token")
				return generatePolicy("user", "Deny", request.MethodArn), nil
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
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
