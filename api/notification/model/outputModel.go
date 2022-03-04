package model

type DataGetListNotification struct {
	Id        string `json:"id,omitempty" bson:"_id,omitempty"`
	ReadFlag  bool   `json:"read_flag" bson:"read_flag"`
	Content   string `json:"content" bson:"content"`
	Type      int    `json:"type" bson:"type"`
	CreatedAt string `json:"created_at" bson:"created_at"`
}

type ResponseGetListNotification struct {
	Code    int                       `json:"code" example:"200"`
	Message string                    `json:"message" example:"success"`
	Data    []DataGetListNotification `json:"data"`
}
