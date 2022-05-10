package repository

import (
	"fmt"
	"houze_ops_backend/api/user/model"
	"houze_ops_backend/db"
)

type userInterface interface {
	CreateUser() (data []model.SysProductType)
}

func PublishInterfaceUser() userInterface {
	return &userResource{}
}

type userResource struct {
}

func (r *userResource) CreateUser() (data []model.SysProductType) {
	err := db.GetConnectionDB().Model(&data).Select()
	if err != nil {
		fmt.Println(err)
	}
	return data
}
