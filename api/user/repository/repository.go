package repository

import (
	"log"
	listenerRepository "sandexcare_backend/api/listener/repository"
	"sandexcare_backend/api/user/model"
	"sandexcare_backend/db"
	"sandexcare_backend/helpers/common"
	"sandexcare_backend/helpers/constant"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userInterface interface {
	CreateUser(data model.User) error
	UpdateUser(data model.User) error
	GetUserByPhoneNumber(phoneNumber string) (data model.User, err error)
	GetUserByPhoneNumberAndCode(phoneNumber string, code string) (data model.User, err error)
	GetUserById(id string) (data model.User, err error)
	CreateListenersBookmark(data model.ListenersBookmark) error
	GetAllListenersBookmark(userId string) (data []model.DataGetListenersBookmark)
	GetListenersBookmarkByUserIdAndListenerId(userId string, listenerId string) (data model.ListenersBookmark, err error)
	DeleteListenerBookmark(userId string, listenerId string) (err error)
	GetAllCouponsUser(userId string) (data []model.DataGetCouponsUser)
	GetCouponInfo(userId string, couponId string) (data model.CouponsUser, err error)
	UpdateCouponsUser(data model.CouponsUser) (err error)
	CreateCouponUser(data model.CouponsUser) (err error)
	RegisterFirebase(userUUID, token string) (err error)
}

func PublishInterfaceUser() userInterface {
	return &userResource{}
}

type userResource struct {
}

func (r *userResource) CreateUser(data model.User) (err error) {
	_, err = db.Collection(model.CollectionUsers).InsertOne(db.GetContext(), data)
	return err
}

func (r *userResource) UpdateUser(data model.User) (err error) {
	_, err = db.Collection(model.CollectionUsers).UpdateOne(db.GetContext(), bson.M{"_id": data.Id}, bson.M{"$set": data})
	return err
}

func (r *userResource) GetUserByPhoneNumber(phoneNumber string) (data model.User, err error) {
	err = db.Collection(model.CollectionUsers).FindOne(db.GetContext(), bson.M{"phone_number": phoneNumber}).Decode(&data)
	return data, err
}

func (r *userResource) GetUserByPhoneNumberAndCode(phoneNumber string, code string) (data model.User, err error) {
	err = db.Collection(model.CollectionUsers).FindOne(db.GetContext(), bson.M{"phone_number": phoneNumber, "code": code}).Decode(&data)
	return data, err
}

func (r *userResource) GetUserById(id string) (data model.User, err error) {
	objectId, errObjectId := primitive.ObjectIDFromHex(id)
	if errObjectId != nil {
		return data, errObjectId
	}
	err = db.Collection(model.CollectionUsers).FindOne(db.GetContext(), bson.M{"_id": objectId}).Decode(&data)
	return data, err
}

func (r *userResource) CreateListenersBookmark(data model.ListenersBookmark) (err error) {
	_, err = db.Collection(model.CollectionListenersBookmark).InsertOne(db.GetContext(), data)
	return err
}

func (r *userResource) GetAllListenersBookmark(userId string) (data []model.DataGetListenersBookmark) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	cur, err := db.Collection(model.CollectionListenersBookmark).Find(db.GetContext(), bson.M{
		"user_id": userId,
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData model.ListenersBookmark
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			var item model.DataGetListenersBookmark
			listenerInfo, _ := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(curData.ListenerId)
			item.PhoneNumber = listenerInfo.PhoneNumber
			item.ListenerId = listenerInfo.Id.Hex()
			item.ListenerRole = listenerInfo.Role
			item.EmployeeId = listenerInfo.EmployeeId
			item.Description = listenerInfo.Description
			item.Avatar = listenerInfo.Avatar
			item.Star = listenerRepository.PublishInterfaceListener().AggregateStarRatingListener(curData.ListenerId, false)
			item.Name = common.GetFullNameOfListener(listenerInfo)
			data = append(data, item)
		}
	}
	return data
}

func (r *userResource) GetListenersBookmarkByUserIdAndListenerId(userId string, listenerId string) (data model.ListenersBookmark, err error) {
	err = db.Collection(model.CollectionListenersBookmark).FindOne(db.GetContext(), bson.M{"user_id": userId, "listener_id": listenerId}).Decode(&data)
	return data, err
}

func (r *userResource) DeleteListenerBookmark(userId string, listenerId string) (err error) {
	_, err = db.Collection(model.CollectionListenersBookmark).DeleteMany(db.GetContext(), bson.M{"user_id": userId, "listener_id": listenerId})
	return err
}

func (r *userResource) GetAllCouponsUser(userId string) (data []model.DataGetCouponsUser) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	cur, err := db.Collection(model.CollectionCouponsUser).Find(db.GetContext(), bson.M{
		"user_id": userId,
		"status":  constant.Active,
		"expires_at": bson.M{
			"$gt": time.Now().UTC(),
		},
	}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData model.CouponsUser
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			var item model.DataGetCouponsUser
			item.CouponName = curData.Name
			item.Discount = curData.Discount
			item.CouponId = curData.Id.Hex()
			item.Status = curData.Status
			item.Type = curData.Type
			item.ExpiresAt = curData.ExpiresAt.Format(constant.DateTimeFormat)
			data = append(data, item)
		}
	}
	return data
}

func (r *userResource) GetCouponInfo(userId string, couponId string) (data model.CouponsUser, err error) {
	objectID, _ := primitive.ObjectIDFromHex(couponId)
	err = db.Collection(model.CollectionCouponsUser).FindOne(db.GetContext(), bson.M{
		"user_id": userId,
		"_id":     objectID,
	}).Decode(&data)
	return data, err
}

func (r *userResource) UpdateCouponsUser(data model.CouponsUser) (err error) {
	_, err = db.Collection(model.CollectionCouponsUser).UpdateOne(db.GetContext(), bson.M{"_id": data.Id}, bson.M{"$set": data})
	return err
}

func (r *userResource) RegisterFirebase(userUUID, token string) (err error) {
	objectID, err := primitive.ObjectIDFromHex(userUUID)
	if err != nil {
		return err
	}
	_, err = db.Collection(model.CollectionUsers).UpdateOne(db.GetContext(),
		bson.M{"_id": objectID}, bson.M{"$set": bson.M{"fb_token": token}})
	return err
}

func (r *userResource) CreateCouponUser(data model.CouponsUser) (err error) {
	_, err = db.Collection(model.CollectionCouponsUser).InsertOne(db.GetContext(), data)
	return err
}
