package model

type ResponseSuccess struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message"`
}

type ResponseError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message"`
}

type ResponseGetTimeSlot struct {
	Code    int        `json:"code" example:"200"`
	Message string     `json:"message" example:"success"`
	Data    []TimeSlot `json:"data"`
}

type ResponseGetConsultingField struct {
	Code    int               `json:"code" example:"200"`
	Message string            `json:"message" example:"success"`
	Data    []ConsultingField `json:"data"`
}

type ConfigData struct {
	AppointmentPrice int64 `json:"appointment_price" example:"30000"`
}

type ResponseGetConfigData struct {
	Code    int        `json:"code" example:"200"`
	Message string     `json:"message" example:"success"`
	Data    ConfigData `json:"data"`
}

type ResponseGetServerDatetime struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"success"`
	Data    string `json:"data"`
}
