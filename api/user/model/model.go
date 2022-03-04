package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUsers             = "users"
	CollectionCoupons           = "coupons"
	CollectionCouponsUser       = "coupons_user"
	CollectionListenersBookmark = "listeners_bookmark"
)

/*User Model*/
type User struct {
	Id               primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	PhoneNumber      string               `json:"phone_number" bson:"phone_number"`
	Name             string               `json:"name" bson:"name"`
	Dob              string               `json:"dob" bson:"dob"`
	Avatar           string               `json:"avatar" bson:"avatar"`
	Diamond          int64                `json:"diamond" bson:"diamond"`
	Status           string               `json:"status" bson:"status"`
	Code             string               `json:"code" bson:"code"`
	ExpiresAt        time.Time            `json:"expires_at" bson:"expires_at"`
	CountSendCode    int                  `json:"count_send_code" bson:"count_send_code"`
	IsMember         bool                 `json:"is_member" bson:"is_member"`
	ConsultingFields []string             `json:"consulting_fields" bson:"consulting_fields"`
	FirebaseToken    string               `json:"fb_token" bson:"fb_token"`
	BookmarkPost     []primitive.ObjectID `json:"bookmark_post" bson:"bookmark_post"`
	RefreshToken     string               `json:"refresh_token" bson:"refresh_token"`
	CreatedAt        time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at" bson:"updated_at"`
	DeletedFlag      bool                 `json:"deleted_flag" bson:"deleted_flag"`
	CountReferral    int                  `json:"count_referral" bson:"count_referral"`
}

/*ListenersBookmark Model*/
type ListenersBookmark struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId       string             `json:"user_id" bson:"user_id"`
	ListenerId   string             `json:"listener_id" bson:"listener_id"`
	ListenerRole int                `json:"listener_role" bson:"listener_role"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

/*CouponsUser Model*/
type CouponsUser struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId    string             `json:"user_id" bson:"user_id"`
	Name      string             `json:"name" bson:"name"`
	Status    string             `json:"status" bson:"status"` //ACTIVE, INACTIVE, USED
	Discount  float64            `json:"discount" bson:"discount"`
	Type      string             `json:"type" bson:"type"` // cp_call_now, cp_booking_cg, cp_booking_cv
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
