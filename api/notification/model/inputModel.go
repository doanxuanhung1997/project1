package model

type InputCreateNotification struct {
	ReceiverId   string `json:"receiver_id" bson:"receiver_id"`
	ReceiverRole int    `json:"receiver_role" bson:"receiver_role"`
	Content      string `json:"content" bson:"content"`
	Type         int    `json:"type" bson:"type"`
}

type InputReadNotification struct {
	Id      string `json:"id" bson:"id"`
	ReadAll bool   `json:"read_all" bson:"read_all"`
}
