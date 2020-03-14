package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jwt "github.com/dgrijalva/jwt-go"

	models "github.com/sawima/aws-sls-jwt-service/functions/layers/models"
)

func handler(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.AuthorizationToken
	tokenSlice := strings.Split(token, " ")
	var tokenStr string
	if len(tokenSlice) > 1 {
		tokenStr = tokenSlice[len(tokenSlice)-1]
	}
	pemFile := []byte("thisisthekimatechtokenstring")
	claims := models.MyCustomClaims{}
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
			return generatePolicy("user", "Deny", request.MethodArn, nil), nil
		}
		log.Println("$$$$$$$$$$$$$$indicate claims$$$$$$$$$$$$$$$$$$$")
		log.Printf("%v", claims)
		log.Printf("%v", claims.Context)
		if token.Valid {
			apiandpath := strings.Split(request.MethodArn, "/")
			resource := strings.Join(apiandpath[0:1], "/")
			resource += "/*"
			return generatePolicy("user", "Allow", resource, &claims), nil
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return generatePolicy("user", "Deny", request.MethodArn, nil), nil
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				log.Println("Expired token")
				return generatePolicy("user", "Deny", request.MethodArn, nil), nil
			} else {
				log.Println("Expired token,any where")
				return generatePolicy("user", "Deny", request.MethodArn, nil), nil
			}
		} else {
			log.Println("Expired token")
			return generatePolicy("user", "Deny", request.MethodArn, nil), nil
		}
	} else {
		log.Println("Token is empty")
		return generatePolicy("user", "Deny", request.MethodArn, nil), nil

	}
}

func generatePolicy(principalID, effect, resource string, claims *models.MyCustomClaims) events.APIGatewayCustomAuthorizerResponse {
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
	// Optional output with custom properties of the String, Number or Boolean type.
	// authResponse.Context = map[string]interface{}{
	// 	"appname":       "appname",
	// 	"uuid":          "123",
	// 	"indicatevalue": "value",
	// }
	var respContext map[string]interface{}
	inrec, _ := json.Marshal(claims.Context)
	json.Unmarshal(inrec, &respContext)

	log.Println("############claim context####################")
	log.Printf("%v", respContext)
	authResponse.Context = respContext
	return authResponse
}

func main() {
	lambda.Start(handler)
}
