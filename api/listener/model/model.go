package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionListeners              = "listeners"
	CollectionListenersResetPassword = "listeners_reset_password"
	CollectionWithdrawalHistory      = "withdrawal_history"
)

/*Listener Model*/
type Listener struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EmployeeId     string             `json:"employee_id" bson:"employee_id"`
	PhoneNumber    string             `json:"phone_number" bson:"phone_number"`
	Email          string             `json:"email" bson:"email"`
	Gender         string             `json:"gender" bson:"gender"`
	Password       string             `json:"password" bson:"password"`
	Status         string             `json:"status" bson:"status"`
	Name           Name               `json:"name" bson:"name"`
	Dob            string             `json:"dob" bson:"dob"`
	Address        string             `json:"address" bson:"address"`
	PersonalId     string             `json:"personal_id" bson:"personal_id"`
	Role           int                `json:"role" bson:"role"`
	Avatar         string             `json:"avatar" bson:"avatar"`
	MainTopic      []string           `json:"main_topic" bson:"main_topic"`
	Description    string             `json:"description" bson:"description"`
	Price          int64              `json:"price" bson:"price"`
	Bank           Bank               `json:"bank" bson:"bank"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedFlag    bool               `json:"deleted_flag" bson:"deleted_flag"`
	CountLoginFail int                `json:"count_login_fail" bson:"count_login_fail"`
	RefreshToken   string             `json:"refresh_token" bson:"refresh_token"`
}

type Name struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
}

type Bank struct {
	Account  string `json:"account" bson:"account" example:"me@thienhang.com"`
	BankName string `json:"bank_name" bson:"bank_name"`
	Owner    string `json:"owner" bson:"owner"`
}

/*ListenerResetPassword Model*/
type ListenerResetPassword struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	ResetCode   string             `json:"reset_code" bson:"reset_code"`
	ExpiresAt   time.Time          `json:"expires_at" bson:"expires_at"`
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

/*WithdrawalHistory Model*/
type WithdrawalHistory struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ListenerId  string             `json:"listener_id" bson:"listener_id"`
	AmountMoney int64              `json:"amount_money" bson:"amount_money"`
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
