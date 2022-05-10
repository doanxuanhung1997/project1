package model

import (
	"time"
)

const (
	CollectionSysUser = "sys_user"
)

//SysUser Model
type SysUser struct {
	PhoneNumber      string    `json:"phone_number" bson:"phone_number"`
	Name             string    `json:"name" bson:"name"`
	Dob              string    `json:"dob" bson:"dob"`
	Avatar           string    `json:"avatar" bson:"avatar"`
	Diamond          int64     `json:"diamond" bson:"diamond"`
	Status           string    `json:"status" bson:"status"`
	Code             string    `json:"code" bson:"code"`
	ExpiresAt        time.Time `json:"expires_at" bson:"expires_at"`
	CountSendCode    int       `json:"count_send_code" bson:"count_send_code"`
	IsMember         bool      `json:"is_member" bson:"is_member"`
	ConsultingFields []string  `json:"consulting_fields" bson:"consulting_fields"`
	FirebaseToken    string    `json:"fb_token" bson:"fb_token"`
	RefreshToken     string    `json:"refresh_token" bson:"refresh_token"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
	DeletedFlag      bool      `json:"deleted_flag" bson:"deleted_flag"`
	CountReferral    int       `json:"count_referral" bson:"count_referral"`
}

type SysProductType struct {
	tableName struct{} `pg:"public.sys_product_type"`

	Id   int    `json:"id" pg:"type:serial"`
	Name string `json:"name" pg:",notnull"`
}
