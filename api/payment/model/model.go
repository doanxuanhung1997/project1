package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionOrderPayment = "order_payment"
)

/*OrderPayment Model*/
type OrderPayment struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId         string             `json:"user_id" bson:"user_id"`
	ListenerId     string             `json:"listener_id" bson:"listener_id"`
	Date           time.Time          `json:"date" bson:"date"`
	TimeSlot       string             `json:"time_slot" bson:"time_slot"`
	BookingTime    string             `json:"booking_time" bson:"booking_time"`
	CallDatetime   time.Time          `json:"call_datetime" bson:"call_datetime"`
	Type           string             `json:"type" bson:"type"`
	Status         string             `json:"status" bson:"status"` //ACTIVE, COMPLETED, REFUNDED
	CouponId       string             `json:"coupon_id" bson:"coupon_id"`
	DiamondOrder   int64              `json:"diamond_order" bson:"diamond_order"`
	Surcharge      int64              `json:"surcharge" bson:"surcharge"`
	DiamondPayment int64              `json:"diamond_payment" bson:"diamond_payment"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	FlagRemind     int                `json:"flag_remind" bson:"flag_remind"`
	UpdateFlag     bool               `json:"update_flag" bson:"update_flag"`
}
