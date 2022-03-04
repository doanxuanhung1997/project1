package model

type InputCreateScheduleWork struct {
	EmployeeId string `json:"employee_id" bson:"employee_id"`
	Date       string `json:"date" bson:"date"`
	TimeSlot   string `json:"time_slot" bson:"time_slot"`
}