package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sandexcare_backend/api/master/model"
	"sandexcare_backend/db"
)

type masterInterface interface {
	GetAllConsultingField() (data []model.ConsultingField)
	GetAllTimeSlot() (data []model.TimeSlot)
	CreateConsultingField(data model.ConsultingField) (err error)
	CreateTimeSlot(data model.TimeSlot) (err error)
	GetDetailTimeSlot(timeSlot string) (data model.TimeSlot, err error)
}

func PublishInterfaceMaster() masterInterface {
	return &masterResource{}
}

type masterResource struct {
}

func (r *masterResource) GetAllConsultingField() (data []model.ConsultingField) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"name", 1}})
	cur, err := db.Collection(model.CollectionConsultingField).Find(db.GetContext(), bson.M{}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData model.ConsultingField
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			data = append(data, curData)
		}
	}
	return data
}

func (r *masterResource) GetAllTimeSlot() (data []model.TimeSlot) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"time_slot", 1}})
	cur, err := db.Collection(model.CollectionTimeSlot).Find(db.GetContext(), bson.M{}, findOptions)
	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	for cur.Next(db.GetContext()) {
		var curData model.TimeSlot
		err = cur.Decode(&curData)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		} else {
			data = append(data, curData)
		}
	}
	return data
}

func (r *masterResource) CreateConsultingField(data model.ConsultingField) (err error) {
	_, err = db.Collection(model.CollectionConsultingField).InsertOne(db.GetContext(), data)
	return err
}

func (r *masterResource) CreateTimeSlot(data model.TimeSlot) (err error) {
	_, err = db.Collection(model.CollectionTimeSlot).InsertOne(db.GetContext(), data)
	return err
}

func (r *masterResource) GetDetailTimeSlot(timeSlot string) (data model.TimeSlot, err error) {
	err = db.Collection(model.CollectionTimeSlot).FindOne(db.GetContext(), bson.M{"time_slot": timeSlot}).Decode(&data)
	return data, err
}
