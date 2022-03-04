package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionNotification = "notification"
)

/*Notification Model*/
type Notification struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ReceiverId   string             `json:"receiver_id" bson:"receiver_id"`
	ReceiverRole int                `json:"receiver_role" bson:"receiver_role"`
	ReadFlag     bool               `json:"read_flag" bson:"read_flag"`
	Type         int                `json:"type" bson:"type"` //1: Appointment, 2: Call
	Content      string             `json:"content" bson:"content"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}




