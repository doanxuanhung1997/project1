package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionConsultingField = "consulting_field"
	CollectionTimeSlot        = "time_slot"
)

/*ConsultingField Model*/
type ConsultingField struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

/*TimeSlot Model*/
type TimeSlot struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TimeSlot    string             `json:"time_slot" bson:"time_slot"`
	Surcharge   float64            `json:"surcharge" bson:"surcharge"`
	BookingTime []string           `json:"booking_time" bson:"booking_time"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}