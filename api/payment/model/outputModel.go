package model

type ResponsePaymentCallNow struct {
	Code    int                `json:"code" example:"200"`
	Message string             `json:"message"`
	Data    DataPaymentCallNow `json:"data"`
}

type DataPaymentCallNow struct {
	OrderId     string `json:"order_id" example:"bc72fg9esvbsjd240v2ihidhv9"`
	CallId      string `json:"call_id" example:" bc72fg9esvbsjd240v2ihidhv9"`
	Token       string `json:"token" example:"data token agora"`
	Channel     string `json:"channel" example:"123456"`
	Uid         string `json:"uid" example:"343f9n39fj23fj40"`
	RingingTime int64  `json:"ringing_time" example:"5"`
}
