package models

import (
	jwt "github.com/dgrijalva/jwt-go"
)

//App is the jwt applicaiton schema
type App struct {
	Appid     string     `json:"appid"`
	Hashedkey string     `json:"hashedkey"`
	Context   Appcontext `json:"context"`
}

//Appcontext is the context info for applications
type Appcontext struct {
	Appname       string `json:"appname"`
	UUID          string `json:"uuid"`
	Indicatevalue string `json:"indicatevalue"`
}

//MyCustomClaims defined the claims data schema
type MyCustomClaims struct {
	Appid   string `json:"appid"`
	Context Appcontext
	jwt.StandardClaims
}
