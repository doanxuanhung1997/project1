package model

// ResponseLogin Struct ResponseLogin
type ResponseLogin struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    LoginData `json:"data"`
}

// LoginData Struct DataLogin
type LoginData struct {
	Id    int    `json:"id" `
	Email string `json:"email" `
	Token string `json:"token" `
}
