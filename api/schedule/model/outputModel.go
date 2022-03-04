package model

type DataScheduleWork struct {
	Date     string `json:"date"`
	TimeSlot string `json:"time_slot"`
	IsBook   bool   `json:"is_book"`
}

type ResponseScheduleWork struct {
	Code    int                `json:"code" example:"200"`
	Message string             `json:"message" example:"Success"`
	Data    []DataScheduleWork `json:"data"`
}

type DataWorkingDay struct {
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
	TimeSlot    string `json:"time_slot"`
	BookingTime string `json:"booking_time"`
}

type ResponseWorkingDay struct {
	Code    int              `json:"code" example:"200"`
	Message string           `json:"message" example:"Success"`
	Data    []DataWorkingDay `json:"data"`
}

type DataListenerInfo struct {
	ListenerId  string  `json:"listener_id"`
	EmployeeId  string  `json:"employee_id"`
	Name        string  `json:"name"`
	Avatar      string  `json:"avatar"`
	Star        float64 `json:"star"`
	Description string  `json:"description"`
}

type DataDetailAppointmentInDay struct {
	BookingTime string             `json:"booking_time"`
	IsFull      bool               `json:"is_full"`
}

type ResponseGetScheduleWorkAppointment struct {
	Code    int                          `json:"code" example:"200"`
	Message string                       `json:"message" example:"Success"`
	Data    []DataDetailAppointmentInDay `json:"data"`
}

type ResponseGetDetailScheduleWorkListener struct {
	Code    int      `json:"code" example:"200"`
	Message string   `json:"message" example:"Success"`
	Data    []string `json:"data" example:"[00:30,01:00]"`
}

type ResponseGetListenerForBookAppointment struct {
	Code    int                `json:"code" example:"200"`
	Message string             `json:"message" example:"Success"`
	Data    []DataListenerInfo `json:"data"`
}
