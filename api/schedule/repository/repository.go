package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	listenerRepository "sandexcare_backend/api/listener/repository"
	paymentModel "sandexcare_backend/api/payment/model"
	"sandexcare_backend/api/schedule/model"
	userRepository "sandexcare_backend/api/user/repository"
	"sandexcare_backend/db"
	"sandexcare_backend/helpers/common"
	"sandexcare_backend/helpers/constant"
	"time"
)

type scheduleInterface interface {
	CreateScheduleWork(data model.ScheduleWork) error
	UpdateScheduleWork(data model.ScheduleWork) error
	GetScheduleWordForListener(listenerId string, dateFrom time.Time, dateTo time.Time) (data []model.DataScheduleWork)
	GetScheduleInDay(listenerId string, workingDay time.Time, timeSlot string) (data []model.DataWorkingDay)
	CheckScheduleWorkExistListener(listenerId string, date time.Time, timeSlot string) (result bool)
	GetListenersWorkTimeSlot(timeSlot string, day time.Time) (data []model.DataListenerInfo)
	GetDetailScheduleWorkListener(listenerId string, date time.Time, timeSlot string) (data model.ScheduleWork, err error)
}

func PublishInterfaceSchedule() scheduleInterface {
	return &scheduleResource{}
}

type scheduleResource struct {
}

func (r *scheduleResource) CreateScheduleWork(data model.ScheduleWork) (err error) {
	_, err = db.Collection(model.CollectionScheduleWork).InsertOne(db.GetContext(), data)
	return err
}

func (r *scheduleResource) UpdateScheduleWork(data model.ScheduleWork) (err error) {
	_, err = db.Collection(model.CollectionScheduleWork).UpdateOne(db.GetContext(), bson.M{"_id": data.Id}, bson.M{"$set": data})
	return err
}

// Get schedule work for listener from schedule_work table
func (r *scheduleResource) GetScheduleWordForListener(listenerId string, dateFrom time.Time, dateTo time.Time) (data []model.DataScheduleWork) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"date", 1}})
	cur, err := db.Collection(model.CollectionScheduleWork).Find(db.GetContext(), bson.M{
		"listener_id": listenerId,
		"status":      constant.Active,
		"date": bson.M{
			"$gte": dateFrom,
			"$lte": dateTo,
		},
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData model.ScheduleWork
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			var itemData model.DataScheduleWork
			workingDay := PublishInterfaceSchedule().GetScheduleInDay(listenerId, curData.Date, curData.TimeSlot)
			if len(workingDay) > 0 {
				itemData.IsBook = true
			}
			itemData.Date = curData.Date.Format(constant.DateFormat)
			itemData.TimeSlot = curData.TimeSlot
			data = append(data, itemData)
		}
	}
	return data
}

// Get schedule work in day for listener
func (r *scheduleResource) GetScheduleInDay(listenerId string, workingDay time.Time, timeSlot string) (data []model.DataWorkingDay) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"booking_time", 1}})
	cur, err := db.Collection(paymentModel.CollectionOrderPayment).Find(db.GetContext(), bson.M{
		"listener_id": listenerId,
		"status":      constant.Active,
		"type":        constant.OrderBookAppointment,
		"date":        workingDay,
		"time_slot":   timeSlot,
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData paymentModel.OrderPayment
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			var itemData model.DataWorkingDay
			itemData.UserId = curData.UserId
			userInfo, errUser := userRepository.PublishInterfaceUser().GetUserById(curData.UserId)
			if errUser == nil {
				itemData.Name = userInfo.Name
			}
			itemData.BookingTime = curData.BookingTime
			itemData.TimeSlot = curData.TimeSlot
			data = append(data, itemData)
		}
	}
	return data
}

// Get schedule work in day for listener
func (r *scheduleResource) GetDetailScheduleWorkInTimeSlot(workingDay time.Time, timeSlot string) (data []model.DataWorkingDay) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"booking_time", 1}})
	cur, err := db.Collection(paymentModel.CollectionOrderPayment).Find(db.GetContext(), bson.M{
		"status":    constant.Active,
		"type":      constant.OrderBookAppointment,
		"date":      workingDay,
		"time_slot": timeSlot,
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData paymentModel.OrderPayment
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			var itemData model.DataWorkingDay
			itemData.UserId = curData.UserId
			userInfo, errUser := userRepository.PublishInterfaceUser().GetUserById(curData.UserId)
			if errUser == nil {
				itemData.Name = userInfo.Name
			}
			itemData.BookingTime = curData.BookingTime
			itemData.TimeSlot = curData.TimeSlot
			data = append(data, itemData)
		}
	}
	return data
}

func (r *scheduleResource) CheckScheduleWorkExistListener(listenerId string, date time.Time, timeSlot string) (result bool) {
	data := model.ScheduleWork{}
	err := db.Collection(model.CollectionScheduleWork).FindOne(db.GetContext(), bson.M{
		"listener_id": listenerId,
		"date":        date,
		"time_slot":   timeSlot,
	}).Decode(&data)
	if err == nil {
		result = true
	}
	return result
}

func (r *scheduleResource) GetListenersWorkTimeSlot(timeSlot string, day time.Time) (data []model.DataListenerInfo) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	cur, err := db.Collection(model.CollectionScheduleWork).Find(db.GetContext(), bson.M{
		"time_slot": timeSlot,
		"date":      day,
		"status":    constant.Active,
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData model.ScheduleWork
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			listenerInfo, _ := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(curData.ListenerId)
			itemData := model.DataListenerInfo{}
			itemData.ListenerId = curData.ListenerId
			itemData.Avatar = listenerInfo.Avatar
			itemData.Name = common.GetFullNameOfListener(listenerInfo)
			itemData.EmployeeId = listenerInfo.EmployeeId
			itemData.Description = listenerInfo.Description
			itemData.Star = listenerRepository.PublishInterfaceListener().AggregateStarRatingListener(curData.ListenerId, false)
			data = append(data, itemData)
		}
	}
	return data
}

func (r *scheduleResource) GetDetailScheduleWorkListener(listenerId string, date time.Time, timeSlot string) (data model.ScheduleWork, err error) {
	err = db.Collection(model.CollectionScheduleWork).FindOne(db.GetContext(), bson.M{
		"listener_id": listenerId,
		"date":        date,
		"time_slot":   timeSlot,
	}).Decode(&data)
	return data, err
}
