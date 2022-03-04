package model

type ResponseVerifyOTPCode struct {
	Id           string `json:"id" `
	PhoneNumber  string `json:"phone_number" `
	IsMember     bool   `json:"is_member" `
	Token        string `json:"token" `
	RefreshToken string `json:"refresh_token" `
}

type DataAppointmentSchedule struct {
	Id            string `json:"id" `
	ListenerId    string `json:"listener_id" `
	ListenerName  string `json:"listener_name" `
	EmployeeId    string `json:"employee_id" `
	ListenerRole  int    `json:"listener_role" `
	ListenerImage string `json:"listener_image" `
	Date          string `json:"date" `
	TimeSlot      string `json:"time_slot" `
	BookingTime   string `json:"booking_time" `
}

type ResponseAppointmentSchedule struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Data    []DataAppointmentSchedule `json:"data"`
}

type ResponseDetailAppointmentSchedule struct {
	Code    int                           `json:"code" example:"200"`
	Message string                        `json:"message" example:"Success"`
	Data    DataDetailAppointmentSchedule `json:"data"`
}

type DataDetailAppointmentSchedule struct {
	ListenerId      string `json:"listener_id" `
	ListenerName    string `json:"listener_name" `
	EmployeeId      string `json:"employee_id" `
	ListenerRole    int    `json:"listener_role" `
	ListenerImage   string `json:"listener_image" `
	Date            string `json:"date" `
	TimeSlot        string `json:"time_slot" `
	BookingTime     string `json:"booking_time" `
	DiamondOrder    int64  `json:"diamond_order" `
	DiamondPayment  int64  `json:"diamond_payment"`
	DiamondDiscount int64  `json:"diamond_discount" `
}

type ResponseSendOTPCode struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Success"`
	Data    string `json:"data"  example:"457564"`
}

type DataGetListenersBookmark struct {
	ListenerId   string  `json:"listener_id"`
	EmployeeId   string  `json:"employee_id"`
	ListenerRole int     `json:"listener_role"`
	PhoneNumber  string  `json:"phone_number"`
	Name         string  `json:"name"`
	Avatar       string  `json:"avatar"`
	Description  string  `json:"description"`
	Star         float64 `json:"star"`
}

type ResponseGetListenersBookmark struct {
	Code    int                        `json:"code" example:"200"`
	Message string                     `json:"message" example:"Success"`
	Data    []DataGetListenersBookmark `json:"data"`
}

type DataGetCouponsUser struct {
	CouponName string  `json:"coupon_name" bson:"coupon_name"`
	Discount   float64 `json:"discount" bson:"discount"`
	CouponId   string  `json:"coupon_id" bson:"coupon_id"`
	Status     string  `json:"status" bson:"status" example:"ACTIVE, INACTIVE, USED"`
	Type       string  `json:"type" bson:"type" example:"cp_call_now, cp_booking_cg, cp_booking_cv"`
	ExpiresAt  string  `json:"expires_at" bson:"expires_at"`
}

type ResponseGetCouponsUser struct {
	Code    int                  `json:"code" example:"200"`
	Message string               `json:"message" example:"Success"`
	Data    []DataGetCouponsUser `json:"data"`
}
