package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionTrackCall = "track_call"
	CollectionTrackActionListener = "track_action_listener"
)

/*TrackCall Model*/
type TrackCall struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CallId     string             `json:"call_id" bson:"call_id"`
	ListenerId string             `json:"listener_id" bson:"listener_id"`
	Action     string             `json:"action" bson:"action"` // JOIN_CALL, MISS_CALL
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

/*TrackActionListener Model*/
type TrackActionListener struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ListenerId string             `json:"listener_id" bson:"listener_id"`
	Action     string             `json:"action" bson:"action"` //CHECK_IN, CHECK_OUT, ON_CALL, OFF_CALL
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}