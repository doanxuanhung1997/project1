package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionScheduleWork = "schedule_work"
)

/*ScheduleWork Model*/
type ScheduleWork struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ListenerId  string             `json:"listener_id" bson:"listener_id"`
	Date        time.Time          `json:"date" bson:"date"`
	TimeSlot    string             `json:"time_slot" bson:"time_slot"`
	Status      string             `json:"status" bson:"status"`
	OrderStatus string             `json:"order_status" bson:"order_status"`
	UserLock    string             `json:"user_lock" bson:"user_lock"`
	TimeLock    time.Time          `json:"time_lock" bson:"time_lock"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
