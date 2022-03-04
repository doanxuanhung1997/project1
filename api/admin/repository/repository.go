package repository

import (
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sandexcare_backend/api/admin/model"
	listenerModel "sandexcare_backend/api/listener/model"
	listenerRepository "sandexcare_backend/api/listener/repository"
	userModel "sandexcare_backend/api/user/model"
	"sandexcare_backend/db"
	"sandexcare_backend/helpers/constant"
)

type adminInterface interface {
	GetAllWithdrawalHistory() (data []model.GetAllWithdrawalHistoryData)
	UpdateWithdrawalHistory(data listenerModel.WithdrawalHistory) error
	GetWithdrawalHistoryById(id string) (data listenerModel.WithdrawalHistory, err error)
	GetAllUsers() (data []model.GetAllUsersData)
}

func PublishInterfaceAdmin() adminInterface {
	return &adminResource{}
}

type adminResource struct {
}

// Get list all withdrawal history
func (r *adminResource) GetAllWithdrawalHistory() (data []model.GetAllWithdrawalHistoryData) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	cur, err := db.Collection(listenerModel.CollectionWithdrawalHistory).Find(db.GetContext(), bson.M{}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData listenerModel.WithdrawalHistory
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			itemData := model.GetAllWithdrawalHistoryData{}
			itemData.Id = curData.Id.Hex()
			itemData.ListenerId = curData.ListenerId
			listenerInfo, _ := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(curData.ListenerId)
			itemData.ListenerName = listenerInfo.Name.FirstName + listenerInfo.Name.LastName
			itemData.Status = curData.Status
			itemData.AmountMoney = curData.AmountMoney
			itemData.CreatedAt = curData.CreatedAt.Format(constant.DateTimeFormat)
			itemData.UpdatedAt = curData.UpdatedAt.Format(constant.DateTimeFormat)
			data = append(data, itemData)
		}
	}
	return data
}

// Confirm withdrawal request of listener
func (r *adminResource) UpdateWithdrawalHistory(data listenerModel.WithdrawalHistory) (err error) {
	_, err = db.Collection(listenerModel.CollectionWithdrawalHistory).UpdateOne(db.GetContext(), bson.M{"_id": data.Id}, bson.M{"$set": data})
	return err
}

func (r *adminResource) GetWithdrawalHistoryById(id string) (data listenerModel.WithdrawalHistory, err error) {
	objectId, errObjectId := primitive.ObjectIDFromHex(id)
	if errObjectId != nil {
		return data, errObjectId
	}
	err = db.Collection(listenerModel.CollectionWithdrawalHistory).FindOne(db.GetContext(), bson.M{"_id": objectId}).Decode(&data)
	return data, err
}

// Get list all users
func (r *adminResource) GetAllUsers() (data []model.GetAllUsersData) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	cur, err := db.Collection(userModel.CollectionUsers).Find(db.GetContext(), bson.M{
		"is_member": true,
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData userModel.User
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			var item model.GetAllUsersData
			item.Id = curData.Id.Hex()
			item.PhoneNumber = curData.PhoneNumber
			item.Name = curData.Name
			item.Diamond = curData.Diamond
			item.Status = curData.Status
			data = append(data, item)
		}
	}
	return data
}
