package model

/*Input Create Listener Data*/
type InputCreateListener struct {
	FirstName   string   `json:"first_name" example:"me@thienhang.com"`
	LastName    string   `json:"last_name" example:"me@thienhang.com"`
	PhoneNumber string   `json:"phone_number" example:"0924202404"`
	Email       string   `json:"email" example:"me@thienhang.com"`
	Account     string   `json:"account" example:"me@thienhang.com"`
	BankName    string   `json:"bank_name"`
	Owner       string   `json:"owner"`
	EmployeeId  string   `json:"employee_id"`
	Gender      string   `json:"gender"`
	Dob         string   `json:"dob" bson:"dob"`
	Address     string   `json:"address"`
	PersonalId  string   `json:"personal_id"`
	Role        int      `json:"role"`
	Avatar      string   `json:"avatar"`
	MainTopic   []string `json:"main_topic"`
	Description string   `json:"description"`
	Price       int64    `json:"price"`
}

/*Input Login*/
type InputLogin struct {
	PhoneNumber string `json:"phone_number" `
	Password    string `json:"password" `
}

/*Input Handle Miss Call*/
type InputHandleMissCall struct {
	CallId string `json:"call_id"`
}

/*Input Handle Request Withdrawal*/
type InputRequestWithdrawal struct {
	Money int64 `json:"money" bson:"money"`
}

/*Input Handle Forgot Password*/
type InputForgotPassword struct {
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}

/*Input Handle Verify Reset Password*/
type InputVerifyResetPassword struct {
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	ResetCode   string `json:"reset_code" bson:"reset_code"`
}


/*Input Handle Reset New Password*/
type InputResetNewPassword struct {
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	ResetCode   string `json:"reset_code" bson:"reset_code"`
	NewPassword string `json:"new_password" bson:"new_password"`
}

/*Input Handle Change Password*/
type InputChangePassword struct {
	Password    string `json:"password" bson:"password"`
	NewPassword string `json:"new_password" bson:"new_password"`
}
