package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionCallEvaluation = "call_evaluation"
	CollectionCallHistory    = "call_history"
)

/*CallEvaluation Model*/
type CallEvaluation struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId       string             `json:"user_id" bson:"user_id"`
	ListenerId   string             `json:"listener_id" bson:"listener_id"`
	CallId       string             `json:"call_id" bson:"call_id"`
	NoteCompany  string             `json:"note_company" bson:"note_company"`
	NoteListener string             `json:"note_listener" bson:"note_listener"`
	Star         int                `json:"star" bson:"star"`
	Tip          int64              `json:"tip" bson:"tip"`
	TipRate      float64            `json:"tip_rate" bson:"tip_rate"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

/*CallHistory Model*/
type CallHistory struct {
	Id                    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OrderId               string             `json:"order_id" bson:"order_id"`
	UserId                string             `json:"user_id" bson:"user_id"`
	ListenerId            string             `json:"listener_id" bson:"listener_id"`
	Channel               string             `json:"channel" bson:"channel"`
	StartCall             time.Time          `json:"start_call" bson:"start_call"`
	EndCall               time.Time          `json:"end_call" bson:"end_call"`
	Status                string             `json:"status" bson:"status"` // STARTED, TALKING, COMPLETED, UNCONNECTED
	Content               string             `json:"content" bson:"content"`
	ConsultingField       []string           `json:"consulting_field" bson:"consulting_field"`
	ConsultingFeeListener int64              `json:"consulting_fee_listener" bson:"consulting_fee_listener"`
	WithdrawalDiamond     int64              `json:"withdrawal_diamond" bson:"withdrawal_diamond"`
	CreatedAt             time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at" bson:"updated_at"`
}