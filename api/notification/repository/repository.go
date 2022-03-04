package repository

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sandexcare_backend/api/notification/model"
	"sandexcare_backend/db"
	"sandexcare_backend/helpers/common"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/message"
	"time"
)

type notificationInterface interface {
	CreateNotification(data model.Notification) error
	GetListNotification(receiverId string, role int) (data []model.DataGetListNotification)
	ReadNotification(id string, readAll bool, idToken string, roleToken int) error
}

func PublishInterfaceNotification() notificationInterface {
	return &notificationResource{}
}

type notificationResource struct {
}

// Create data notification
func (r *notificationResource) CreateNotification(data model.Notification) (err error) {
	_, err = db.Collection(model.CollectionNotification).InsertOne(db.GetContext(), data)
	return err
}

// Get list notification
func (r *notificationResource) GetListNotification(receiverId string, role int) (data []model.DataGetListNotification) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	cur, err := db.Collection(model.CollectionNotification).Find(db.GetContext(), bson.M{
		"receiver_id":   receiverId,
		"receiver_role": role,
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData model.Notification
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			var item model.DataGetListNotification
			item.Id = curData.Id.Hex()
			item.Content = curData.Content
			item.ReadFlag = curData.ReadFlag
			item.Type = curData.Type
			item.CreatedAt = curData.CreatedAt.Format(constant.DateTimeFormat)
			data = append(data, item)
		}
	}
	return data
}

// Read notification
func (r *notificationResource) ReadNotification(id string, readAll bool, idToken string, roleToken int) error {
	if common.IsEmpty(id) && !readAll {
		return errors.New(message.MessageErrorReadNotifyFail)
	}
	var condition = bson.M{}
	if readAll {
		condition = bson.M{
			"receiver_id":   idToken,
			"receiver_role": roleToken,
		}
	}else {
		notifyObjectId, errObjectId := primitive.ObjectIDFromHex(id)
		if errObjectId != nil {
			return errObjectId
		}
		condition = bson.M{
			"_id": notifyObjectId,
		}
	}
	_, err := db.Collection(model.CollectionNotification).UpdateMany(db.GetContext(),
		condition,
		bson.M{"$set": bson.M{
			"read_flag":    true,
			"updated_at": time.Now().UTC(),
		}})
	return err
}
