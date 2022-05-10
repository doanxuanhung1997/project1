package repository

import (
	"houze_ops_backend/api/user/model"
	"houze_ops_backend/db"
)

type userInterface interface {
	CreateUser(data model.SysUser) error
}

func PublishInterfaceUser() userInterface {
	return &userResource{}
}

type userResource struct {
}

func (r *userResource) CreateUser(data model.SysUser) (err error) {
	_, err = db.Collection(model.CollectionSysUser).InsertOne(db.GetContext(), data)
	return err
}
