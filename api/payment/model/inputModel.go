package model

type InputOrderPayment struct {
	ListenerId  string `json:"listener_id" bson:"listener_id"`
	Date        string `json:"date" bson:"date"`
	BookingTime string `json:"booking_time" bson:"booking_time"`
	TimeSlot    string `json:"time_slot" bson:"time_slot"`
	CouponId    string `json:"coupon_id" bson:"coupon_id"`
}

type InputCallPayment struct {
	Date        string `json:"date" bson:"date"`
	BookingTime string `json:"booking_time" bson:"booking_time"`
	TimeSlot    string `json:"time_slot" bson:"time_slot"`
	CouponId    string `json:"coupon_id" bson:"coupon_id"`
}

type InputUpdateOrderPayment struct {
	OrderId     string `json:"order_id"`
	ListenerId  string `json:"listener_id"`
	Date        string `json:"date"`
	BookingTime string `json:"booking_time"`
	TimeSlot    string `json:"time_slot"`
	Surcharge   int64  `json:"surcharge"`
}

type InputPaymentRefund struct {
	OrderId string `json:"order_id"`
	Type    int    `json:"type"`
}

type InputUnlock struct {
	ListenerId  string `json:"listener_id"`
	Date        string `json:"date"`
	BookingTime string `json:"booking_time"`
}
