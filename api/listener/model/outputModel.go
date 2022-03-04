package model

/*Struct ResponseDataLogin*/
type ResponseLogin struct {
	Code    int               `json:"code" example:"200"`
	Message string            `json:"message" example:"Success"`
	Data    LoginData `json:"data"`
}

/*Data of Struct ResponseDataLogin*/
type LoginData struct {
	Id           string `json:"id" `
	PhoneNumber  string `json:"phone_number" `
	Email        string `json:"email" `
	Token        string `json:"token" `
	RefreshToken string `json:"refresh_token" `
}

type AggregateRevenue struct {
	CallNumber  int   `json:"call_number" bson:"call_number"`
	TotalStart  int   `json:"total_star" bson:"total_star"`
	CountRating int   `json:"count_rating" bson:"count_rating"`
	Diamond     int64 `json:"diamond" bson:"diamond"`
}

type TotalWithdrawalDiamond struct {
	Diamond int64 `json:"diamond" bson:"diamond"`
}

/*Struct ResponseGetInfoListener*/
type ResponseGetInfoListener struct {
	Code    int      `json:"code" example:"200"`
	Message string   `json:"message" example:"Success"`
	Data    Listener `json:"data"`
}

/*Struct ResponseGetWithdrawalHistory*/
type ResponseGetWithdrawalHistory struct {
	Code    int                 `json:"code" example:"200"`
	Message string              `json:"message" example:"Success"`
	Data    []WithdrawalHistory `json:"data"`
}

/*Struct ResponseGetRevenueAnalysis*/
type ResponseGetRevenueAnalysis struct {
	Code    int                            `json:"code" example:"200"`
	Message string                         `json:"message" example:"Success"`
	Data    GetRevenueAnalysisData `json:"data"`
}

/*Data of ResponseGetRevenueAnalysisData*/
type GetRevenueAnalysisData struct {
	CallNumberSuccess   int     `json:"call_number_success" `
	CallNumberFail      int     `json:"call_number_fail" `
	Rating              float64 `json:"rating" `
	TotalRevenue        int64   `json:"total_revenue" `
	AvailableWithdrawal int64   `json:"available_withdrawal" `
}
