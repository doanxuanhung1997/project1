package repository

import (
	"houze_ops_backend/api/sys_user/model"
	"houze_ops_backend/db"
)

type userInterface interface {
	Login(email string, password string) (data model.SysUser, err error)
	CreateUser(data model.SysUser) error
}

func PublishInterfaceUser() userInterface {
	return &userResource{}
}

type userResource struct {
}

func (r *userResource) CreateUser(data model.SysUser) error {
	_, err := db.GetConnectionDB().Model(&data).Insert()
	return err
}

func (r *userResource) Login(email string, password string) (data model.SysUser, err error) {
	err = db.GetConnectionDB().Model(&data).
		Where("email = ?", email).
		Where("password = ?", password).
		Limit(1).
		Select()
	return data, err
}
