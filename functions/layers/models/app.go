package models

//App is the jwt applicaiton schema
type App struct {
	Appid     string     `json:"appid"`
	Hashedkey string     `json:"hashedkey"`
	Appname   string     `json:"appname"`
	Context   Appcontext `json:"context"`
}

//Appcontext is the context info for applications
type Appcontext struct {
	Org   string `json:"org"`
	Orgid string `json:"orgid"`
}
