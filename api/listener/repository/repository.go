package repository

import (
	"errors"
	"fmt"
	"log"
	"math"
	conversationModel "sandexcare_backend/api/conversation/model"
	"sandexcare_backend/api/listener/model"
	"sandexcare_backend/db"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/message"

	"github.com/jinzhu/now"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type listenerInterface interface {
	CreateListener(data model.Listener) error
	UpdateListener(data model.Listener) error
	GetListenerByPhoneNumber(phoneNumber string) (data model.Listener, err error)
	GetListenerByEmployeeId(employeeId string) (data model.Listener, err error)
	GetListenerByListenerId(listenerId string) (data model.Listener, err error)
	Login(phoneNumber string, password string) (data model.Listener, err error)
	CreateListenerResetPassword(data model.ListenerResetPassword) error
	UpdateListenerResetPassword(data model.ListenerResetPassword) error
	GetListenerResetPasswordByPhoneNumberAndCode(phoneNumber string, code string) (data model.ListenerResetPassword, err error)
	GetDataResetPassword(phoneNumber string, code string, status string) (data model.ListenerResetPassword, err error)
	ClearDataResetPassword(phoneNumber string) error
	CreateRequestWithdrawal(data model.WithdrawalHistory) (err error)
	GetWithdrawalHistory(listenerId string) (data []model.WithdrawalHistory)
	GetTotalRemainingWithdrawalDiamond(listenerId string) (data int64)
	GetRemainingWithdrawalDiamondDetail(listenerId string) (data []conversationModel.CallHistory)
	CountCallNumberListenerByStatus(listenerId string, status string) (data int)
	AggregateStarRatingListener(listenerId string, thisMonth bool) (data float64)
	GetTotalRevenueDiamondListener(listenerId string) (data int64)
}

func PublishInterfaceListener() listenerInterface {
	return &listenerResource{}
}

type listenerResource struct {
}

func (r *listenerResource) CreateListener(data model.Listener) (err error) {
	_, err = db.Collection(model.CollectionListeners).InsertOne(db.GetContext(), data)
	return err
}

func (r *listenerResource) GetListenerByPhoneNumber(phoneNumber string) (data model.Listener, err error) {
	err = db.Collection(model.CollectionListeners).FindOne(db.GetContext(), bson.M{"phone_number": phoneNumber}).Decode(&data)
	return data, err
}

func (r *listenerResource) GetListenerByEmployeeId(employeeId string) (data model.Listener, err error) {
	err = db.Collection(model.CollectionListeners).FindOne(db.GetContext(), bson.M{"employee_id": employeeId}).Decode(&data)
	return data, err
}

func (r *listenerResource) GetListenerByListenerId(listenerId string) (data model.Listener, err error) {
	listenerObjectId, errObjectId := primitive.ObjectIDFromHex(listenerId)
	if errObjectId != nil {
		return data, errors.New(message.MessageErrorListenerIdInvalid)
	}
	err = db.Collection(model.CollectionListeners).FindOne(db.GetContext(), bson.M{"_id": listenerObjectId}).Decode(&data)
	if err != nil {
		return data, errors.New(message.MessageErrorListenerIdNotExist)
	}
	return data, err
}

func (r *listenerResource) UpdateListener(data model.Listener) (err error) {
	_, err = db.Collection(model.CollectionListeners).UpdateOne(db.GetContext(), bson.M{"_id": data.Id}, bson.M{"$set": data})
	return err
}

func (r *listenerResource) Login(phoneNumber string, password string) (data model.Listener, err error) {
	err = db.Collection(model.CollectionListeners).FindOne(db.GetContext(), bson.M{"phone_number": phoneNumber, "password": password}).Decode(&data)
	return data, err
}

func (r *listenerResource) CreateListenerResetPassword(data model.ListenerResetPassword) (err error) {
	_, err = db.Collection(model.CollectionListenersResetPassword).InsertOne(db.GetContext(), data)
	return err
}

func (r *listenerResource) UpdateListenerResetPassword(data model.ListenerResetPassword) (err error) {
	_, err = db.Collection(model.CollectionListenersResetPassword).UpdateOne(db.GetContext(), bson.M{"_id": data.Id}, bson.M{"$set": data})
	return err
}

func (r *listenerResource) GetListenerResetPasswordByPhoneNumberAndCode(phoneNumber string, code string) (data model.ListenerResetPassword, err error) {
	err = db.Collection(model.CollectionListenersResetPassword).FindOne(db.GetContext(), bson.M{"phone_number": phoneNumber, "reset_code": code}).Decode(&data)
	return data, err
}

func (r *listenerResource) GetDataResetPassword(phoneNumber string, code string, status string) (data model.ListenerResetPassword, err error) {
	err = db.Collection(model.CollectionListenersResetPassword).FindOne(db.GetContext(), bson.M{"phone_number": phoneNumber, "reset_code": code, "status": status}).Decode(&data)
	return data, err
}

func (r *listenerResource) ClearDataResetPassword(phoneNumber string) (err error) {
	_, err = db.Collection(model.CollectionListenersResetPassword).DeleteMany(db.GetContext(), bson.M{"phone_number": phoneNumber})
	return err
}

func (r *listenerResource) CreateRequestWithdrawal(data model.WithdrawalHistory) (err error) {
	_, err = db.Collection(model.CollectionWithdrawalHistory).InsertOne(db.GetContext(), data)
	return err
}

// Get list withdrawal history of listener
func (r *listenerResource) GetWithdrawalHistory(listenerId string) (data []model.WithdrawalHistory) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	cur, err := db.Collection(model.CollectionWithdrawalHistory).Find(db.GetContext(), bson.M{
		"listener_id": listenerId,
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData model.WithdrawalHistory
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			data = append(data, curData)
		}
	}
	return data
}

//Get total withdrawal diamond of listener
func (r *listenerResource) GetTotalRemainingWithdrawalDiamond(listenerId string) (data int64) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"listener_id": listenerId,
				"status":      constant.Completed,
				"withdrawal_diamond": bson.M{
					"$gt": 0,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":     bson.M{},
				"diamond": bson.M{"$sum": "$withdrawal_diamond"},
			},
		},
	}
	curTotal, logErrorTotal := db.Collection(conversationModel.CollectionCallHistory).Aggregate(db.GetContext(), pipeline)
	if logErrorTotal != nil {
		println(logErrorTotal.Error())
		return data
	}
	for curTotal.Next(db.GetContext()) {
		var curData model.TotalWithdrawalDiamond
		err := curTotal.Decode(&curData)
		if err != nil {
			fmt.Print("Error on Decoding the document ", err)
		}
		data = curData.Diamond
	}
	return data
}

// Get list withdrawal diamond of listener
func (r *listenerResource) GetRemainingWithdrawalDiamondDetail(listenerId string) (data []conversationModel.CallHistory) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	cur, err := db.Collection(conversationModel.CollectionCallHistory).Find(db.GetContext(), bson.M{
		"listener_id": listenerId,
		"status":      constant.Completed,
		"withdrawal_diamond": bson.M{
			"$gt": 0,
		},
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData conversationModel.CallHistory
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			data = append(data, curData)
		}
	}
	return data
}

//Count call number of listener by call status
func (r *listenerResource) CountCallNumberListenerByStatus(listenerId string, status string) (data int) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"listener_id": listenerId,
				"status":      status,
			},
		},
		{
			"$count": "call_number",
		},
	}
	curTotal, logErrorTotal := db.Collection(conversationModel.CollectionCallHistory).Aggregate(db.GetContext(), pipeline)
	if logErrorTotal != nil {
		println(logErrorTotal.Error())
		return data
	}
	for curTotal.Next(db.GetContext()) {
		var curData model.AggregateRevenue
		err := curTotal.Decode(&curData)
		if err != nil {
			fmt.Print("Error on Decoding the document ", err)
		}
		data = curData.CallNumber
	}
	return data
}

//Count star of listener
func (r *listenerResource) AggregateStarRatingListener(listenerId string, thisMonth bool) (data float64) {
	var match = bson.M{
		"$match": bson.M{
			"listener_id": listenerId,
		},
	}

	if thisMonth {
		timeFrom := now.BeginningOfMonth()
		timeTo := now.EndOfMonth()
		match = bson.M{
			"$match": bson.M{
				"listener_id": listenerId,
				"created_at": bson.M{
					"$gte": timeFrom,
					"$lte": timeTo,
				},
			},
		}
	}

	pipeline := []bson.M{
		match,
		{
			"$group": bson.M{
				"_id":          bson.M{},
				"total_star":   bson.M{"$sum": "$star"},
				"count_rating": bson.M{"$sum": 1},
			},
		},
	}
	curTotal, logErrorTotal := db.Collection(conversationModel.CollectionCallEvaluation).Aggregate(db.GetContext(), pipeline)
	if logErrorTotal != nil {
		println(logErrorTotal.Error())
		return data
	}
	for curTotal.Next(db.GetContext()) {
		var curData model.AggregateRevenue
		err := curTotal.Decode(&curData)
		if err != nil {
			fmt.Print("Error on Decoding the document ", err)
		}
		data = float64(curData.TotalStart) / float64(curData.CountRating)
	}
	return math.Ceil(data*100) / 100
}

//Get total revenue diamonds star of listener
func (r *listenerResource) GetTotalRevenueDiamondListener(listenerId string) (data int64) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"listener_id": listenerId,
				"status":      constant.Completed,
				"consulting_fee_listener": bson.M{
					"$gt": 0,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":     bson.M{},
				"diamond": bson.M{"$sum": "$consulting_fee_listener"},
			},
		},
	}
	curTotal, logErrorTotal := db.Collection(conversationModel.CollectionCallHistory).Aggregate(db.GetContext(), pipeline)
	if logErrorTotal != nil {
		println(logErrorTotal.Error())
		return data
	}
	for curTotal.Next(db.GetContext()) {
		var curData model.AggregateRevenue
		err := curTotal.Decode(&curData)
		if err != nil {
			fmt.Print("Error on Decoding the document ", err)
		}
		data = curData.Diamond
	}

	// get tip revenue of listener
	pipeline2 := []bson.M{
		{
			"$match": bson.M{
				"listener_id": listenerId,
				"tip": bson.M{
					"$gt": 0,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{},
				"diamond": bson.M{"$sum": bson.M{
					"$multiply": []string{"$tip", "$tip_rate"},
				}},
			},
		},
	}
	curTotalTip, logErrorTotalTip := db.Collection(conversationModel.CollectionCallEvaluation).Aggregate(db.GetContext(), pipeline2)
	if logErrorTotalTip != nil {
		println(logErrorTotalTip.Error())
		return data
	}
	for curTotalTip.Next(db.GetContext()) {
		var curDataTip model.AggregateRevenue
		err := curTotalTip.Decode(&curDataTip)
		if err != nil {
			fmt.Print("Error on Decoding the document ", err)
		}
		data += curDataTip.Diamond
	}

	return data
}
