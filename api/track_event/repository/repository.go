package repository

import (
	"sandexcare_backend/api/track_event/model"
	"sandexcare_backend/db"
)

type trackActionInterface interface {
	CreateTrackActionListener(data model.TrackActionListener) error
	CreateTrackCall(data model.TrackCall) error
}

func PublishInterfaceTrackAction() trackActionInterface {
	return &trackActionResource{}
}

type trackActionResource struct {
}

func (r *trackActionResource) CreateTrackActionListener(data model.TrackActionListener) (err error) {
	_, err = db.Collection(model.CollectionTrackActionListener).InsertOne(db.GetContext(), data)
	return err
}

func (r *trackActionResource) CreateTrackCall(data model.TrackCall) (err error) {
	_, err = db.Collection(model.CollectionTrackCall).InsertOne(db.GetContext(), data)
	return err
}