package repository

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sandexcare_backend/api/conversation/model"
	listenerRepository "sandexcare_backend/api/listener/repository"
	userRepository "sandexcare_backend/api/user/repository"
	"sandexcare_backend/db"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/message"
)

type conversationInterface interface {
	CreateCallHistory(data model.CallHistory) error
	UpdateCallHistory(data model.CallHistory) error
	GetCallHistoryOfUser(userId string) (data []model.DataGetCallHistoryUser)
	GetCallHistoryOfListener(listenerId string, status string) (data []model.DataGetCallHistoryListener)
	GetCallHistoryById(callId string) (data model.CallHistory, err error)
	CreateCallEvaluation(data model.CallEvaluation) (err error)
	GetCallEvaluationByCallId(callId string) (data model.CallEvaluation, err error)
	GetConversationsForUser(userId string) (data []model.CallHistory)
	GetLastCallOfListener(listenerId string) (data model.CallHistory, err error)
	GetCallInProgressOfListener(listenerId string) (data model.CallHistory, err error)
	GetCallHistoryOfOrder(orderId string) (data []model.CallHistory)
}

func PublishInterfaceConversation() conversationInterface {
	return &conversationResource{}
}

type conversationResource struct {
}

func (r *conversationResource) CreateCallHistory(data model.CallHistory) (err error) {
	_, err = db.Collection(model.CollectionCallHistory).InsertOne(db.GetContext(), data)
	return err
}

func (r *conversationResource) UpdateCallHistory(data model.CallHistory) (err error) {
	_, err = db.Collection(model.CollectionCallHistory).UpdateOne(db.GetContext(), bson.M{"_id": data.Id}, bson.M{"$set": data})
	return err
}

// Get conversation history of user
func (r *conversationResource) GetCallHistoryOfUser(userId string) (data []model.DataGetCallHistoryUser) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"user_id": userId,
				"status":  constant.Completed,
			},
		},
		{
			"$project": bson.M{
				"_id":              1,
				"listener_id":      1,
				"start_call":       1,
				"end_call":         1,
				"content":          1,
				"consulting_field": 1,
			},
		},
		{
			"$sort": bson.M{"start_call": 1},
		},
	}
	curTotal, logErrorTotal := db.Collection(model.CollectionCallHistory).Aggregate(db.GetContext(), pipeline)
	if logErrorTotal != nil {
		println(logErrorTotal.Error())
		return data
	}
	for curTotal.Next(db.GetContext()) {
		var curData model.CallHistory
		err := curTotal.Decode(&curData)
		if err != nil {
			fmt.Print("Error on Decoding the document ", err)
		} else {
			var item model.DataGetCallHistoryUser
			listener, errListener := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(curData.ListenerId)
			if errListener == nil {
				item.ListenerName = listener.Name.FirstName + listener.Name.LastName
			}
			item.CallId = curData.Id.Hex()
			item.ListenerId = curData.ListenerId
			item.StartCall = curData.StartCall.Format(constant.DateTimeFormat)
			item.EndCall = curData.EndCall.Format(constant.DateTimeFormat)
			item.Content = curData.Content
			item.ConsultingField = curData.ConsultingField
			data = append(data, item)
		}
	}
	return data
}

// Get detail conversation history of listener
func (r *conversationResource) GetCallHistoryById(callId string) (data model.CallHistory, err error) {
	callObjectId, errObjectId := primitive.ObjectIDFromHex(callId)
	if errObjectId != nil {
		return data, errors.New(message.MessageErrorCallIdInvalid)
	}
	err = db.Collection(model.CollectionCallHistory).FindOne(db.GetContext(), bson.M{"_id": callObjectId}).Decode(&data)
	return data, err
}

// Get conversation history of listener
func (r *conversationResource) GetCallHistoryOfListener(listenerId string, status string) (data []model.DataGetCallHistoryListener) {
	var match = bson.M{
		"$match": bson.M{
			"listener_id": listenerId,
		},
	}

	if status == constant.Completed {
		match = bson.M{
			"$match": bson.M{
				"listener_id": listenerId,
				"status":      constant.Completed,
			},
		}
	}
	pipeline := []bson.M{
		match,
		{
			"$sort": bson.M{"start_call": 1},
		},
	}
	curTotal, logErrorTotal := db.Collection(model.CollectionCallHistory).Aggregate(db.GetContext(), pipeline)
	if logErrorTotal != nil {
		println(logErrorTotal.Error())
		return data
	}
	for curTotal.Next(db.GetContext()) {
		var curData model.CallHistory
		err := curTotal.Decode(&curData)
		if err != nil {
			fmt.Print("Error on Decoding the document ", err)
		} else {
			var item model.DataGetCallHistoryListener
			userInfo, errUser := userRepository.PublishInterfaceUser().GetUserById(curData.UserId)
			if errUser == nil {
				item.UserName = userInfo.Name
				item.PhoneNumber = userInfo.PhoneNumber
			}
			item.CallId = curData.Id.Hex()
			item.UserId = curData.UserId
			item.StartCall = curData.StartCall.Format(constant.DateTimeFormat)
			item.EndCall = curData.EndCall.Format(constant.DateTimeFormat)
			if status == constant.Completed {
				item.Content = curData.Content
				item.ConsultingField = curData.ConsultingField
			} else {
				item.Status = curData.Status
				item.ConsultingFeeListener = curData.ConsultingFeeListener
			}
			data = append(data, item)
		}
	}
	return data
}

func (r *conversationResource) CreateCallEvaluation(data model.CallEvaluation) (err error) {
	_, err = db.Collection(model.CollectionCallEvaluation).InsertOne(db.GetContext(), data)
	return err
}

// Get detail call evaluation by call id
func (r *conversationResource) GetCallEvaluationByCallId(callId string) (data model.CallEvaluation, err error) {
	err = db.Collection(model.CollectionCallEvaluation).FindOne(db.GetContext(), bson.M{"call_id": callId}).Decode(&data)
	return data, err
}

func (r *conversationResource) GetConversationsForUser(userId string) (data []model.CallHistory) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	cur, err := db.Collection(model.CollectionCallHistory).Find(db.GetContext(), bson.M{
		"status":  constant.Completed,
		"user_id": userId,
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData model.CallHistory
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			data = append(data, curData)
		}
	}
	return data
}

func (r *conversationResource) GetLastCallOfListener(listenerId string) (data model.CallHistory, err error) {
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"updated_at", -1}})
	err = db.Collection(model.CollectionCallHistory).FindOne(db.GetContext(), bson.M{
		"listener_id":   listenerId,
		"status":        constant.Completed,
	}, findOptions).Decode(&data)
	return data, err
}

func (r *conversationResource) GetCallInProgressOfListener(listenerId string) (data model.CallHistory, err error) {
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"updated_at", -1}})
	err = db.Collection(model.CollectionCallHistory).FindOne(db.GetContext(), bson.M{
		"listener_id":   listenerId,
		"status":        constant.Talking,
	}, findOptions).Decode(&data)
	return data, err
}


func (r *conversationResource) GetCallHistoryOfOrder(orderId string) (data []model.CallHistory) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	cur, err := db.Collection(model.CollectionCallHistory).Find(db.GetContext(), bson.M{
		"order_id": orderId,
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData model.CallHistory
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			data = append(data, curData)
		}
	}
	return data
}