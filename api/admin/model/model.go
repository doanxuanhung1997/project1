package model

// Input confirm withdrawal request
type InputConfirmWithdrawal struct {
	Id     string `json:"id" `
	Status string `json:"status" `
}

// Input submit diamonds for user
type InputSubmitDiamond struct {
	Id      string `json:"id" `
	Diamond int64  `json:"diamond" `
}

// Struct response get all withdrawal request
type ResponseGetAllWithdrawalHistory struct {
	Code    int                           `json:"code" example:"200"`
	Message string                        `json:"message" example:"Success"`
	Data    []GetAllWithdrawalHistoryData `json:"data"`
}

// Get all withdrawal request data
type GetAllWithdrawalHistoryData struct {
	Id           string `json:"id" `
	ListenerId   string `json:"listener_id"`
	ListenerName string `json:"listener_name"`
	AmountMoney  int64  `json:"amount_money"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// Struct Response get all users
type ResponseGetAllUsers struct {
	Code    int               `json:"code" example:"200"`
	Message string            `json:"message" example:"Success"`
	Data    []GetAllUsersData `json:"data"`
}

// Get all users data
type GetAllUsersData struct {
	Id          string `json:"id" `
	PhoneNumber string `json:"phone_number" `
	Name        string `json:"name" `
	Diamond     int64  `json:"diamond" `
	Status      string `json:"status" `
}
