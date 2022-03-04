package model

/*Input Send OTP Code*/
type InputSendOTPCode struct {
	PhoneNumber string `json:"phone_number" bson:"phone_number" example:"0335266678"`
}

/*Input flow OTP code*/
type InputVerifyOTPCode struct {
	PhoneNumber string `json:"phone_number" bson:"phone_number" example:"0335266678"`
	Code        string `json:"code" bson:"code" example:"462384"`
}

/*Input complete user info*/
type InputCompleteInfo struct {
	Name             string   `json:"name" bson:"name" example:"HungDX"`
	ReferralCode     string   `json:"referral_code" bson:"referral_code" example:"0335229337"`
	ConsultingFields []string `json:"consulting_fields" bson:"consulting_fields"`
}

/*Input update user info*/
type InputUpdateInfo struct {
	Avatar           string   `json:"avatar" bson:"avatar"`
	Name             string   `json:"name" bson:"name"`
	ConsultingFields []string `json:"consulting_fields" bson:"consulting_fields"`
}

/*Input bookmark listener*/
type InputBookmarkListener struct {
	ListenerId string `json:"listener_id" bson:"listener_id"`
}
