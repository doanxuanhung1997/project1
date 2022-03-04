package model

/*Input create master data consulting_field or time_slot*/
type InputCreateMasterData struct {
	Table       string   `json:"table"`
	Value       string   `json:"value"`
	BookingTime []string `json:"booking_time"`
}

// Input send event web socket
type InputSendEventWebSocket struct {
	FromUserId string `json:"from_user_id"`
	Event      string `json:"event"`
	Message    string `json:"message"`
}