package repository

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sandexcare_backend/api/payment/model"
	"sandexcare_backend/db"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/message"
	"time"
)

type paymentInterface interface {
	CreateOrderPayment(data model.OrderPayment) error
	GetOrderPaymentById(id string) (data model.OrderPayment, err error)
	GetUpcomingAppointment() (data []model.OrderPayment)
	CheckOrderPaymentExist(listenerId string, callDatetime time.Time) bool
	GetScheduleAppointmentForUser(userId string) (data []model.OrderPayment)
	UpdateOrderPayment(data model.OrderPayment) error
	GetAppointmentsBooked(date time.Time, bookingTime string) (data []model.OrderPayment)
	GetUpcomingAppointmentOfListener(listenerId string) (data model.OrderPayment, err error)
}

func PublishInterfacePayment() paymentInterface {
	return &paymentResource{}
}

type paymentResource struct {
}

func (r *paymentResource) CreateOrderPayment(data model.OrderPayment) (err error) {
	_, err = db.Collection(model.CollectionOrderPayment).InsertOne(db.GetContext(), data)
	return err
}

func (r *paymentResource) UpdateOrderPayment(data model.OrderPayment) (err error) {
	_, err = db.Collection(model.CollectionOrderPayment).UpdateOne(db.GetContext(), bson.M{"_id": data.Id}, bson.M{"$set": data})
	return err
}

func (r *paymentResource) GetOrderPaymentById(id string) (data model.OrderPayment, err error) {
	objectId, errObjectId := primitive.ObjectIDFromHex(id)
	if errObjectId != nil {
		return data, errors.New(message.MessageErrorOrderIdInvalid)
	}
	err = db.Collection(model.CollectionOrderPayment).FindOne(db.GetContext(), bson.M{"_id": objectId}).Decode(&data)
	return data, err
}

func (r *paymentResource) GetUpcomingAppointment() (data []model.OrderPayment) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"call_datetime", 1}})
	cur, err := db.Collection(model.CollectionOrderPayment).Find(db.GetContext(), bson.M{
		"status": constant.Active,
		"type":   constant.OrderBookAppointment,
		"call_datetime": bson.M{
			"$gt": time.Now().UTC().Add(constant.UTC7 * time.Hour),
		},
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData model.OrderPayment
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			data = append(data, curData)
		}
	}
	return data
}

func (r *paymentResource) GetUpcomingAppointmentOfListener(listenerId string) (data model.OrderPayment, err error) {
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"call_datetime", 1}})
	err = db.Collection(model.CollectionOrderPayment).FindOne(db.GetContext(), bson.M{
		"listener_id": listenerId,
		"call_datetime": bson.M{
			"$gt": time.Now().UTC().Add(constant.UTC7 * time.Hour),
		},
		"status": constant.Active,
	}, findOptions).Decode(&data)
	return data, err
}

func (r *paymentResource) CheckOrderPaymentExist(listenerId string, callDatetime time.Time) (result bool) {
	data := model.OrderPayment{}
	err := db.Collection(model.CollectionOrderPayment).FindOne(db.GetContext(), bson.M{
		"listener_id":   listenerId,
		"call_datetime": callDatetime,
		"status":        constant.Active,
	}).Decode(&data)
	if err == nil {
		result = true
	}
	return result
}

func (r *paymentResource) GetScheduleAppointmentForUser(userId string) (data []model.OrderPayment) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"call_datetime", 1}})
	cur, err := db.Collection(model.CollectionOrderPayment).Find(db.GetContext(), bson.M{
		"status":  constant.Active,
		"type":    constant.OrderBookAppointment,
		"user_id": userId,
		"call_datetime": bson.M{
			"$gt": time.Now().UTC().Add(constant.UTC7 * time.Hour).Add(-30 * time.Minute),
		},
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData model.OrderPayment
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			data = append(data, curData)
		}
	}
	return data
}

func (r *paymentResource) GetAppointmentsBooked(date time.Time, bookingTime string) (data []model.OrderPayment) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	cur, err := db.Collection(model.CollectionOrderPayment).Find(db.GetContext(), bson.M{
		"status":       constant.Active,
		"type":         constant.OrderBookAppointment,
		"date":         date,
		"booking_time": bookingTime,
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData model.OrderPayment
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			data = append(data, curData)
		}
	}
	return data
}
