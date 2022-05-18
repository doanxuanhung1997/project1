package model

import (
	"houze_ops_backend/helpers/common"
	"time"
)

//SysUser Model
type SysUser struct {
	tableName struct{} `pg:"public.sys_user"`

	Id         int       `json:"id" pg:"type:serial,pk"`
	FirstName  string    `json:"first_name" pg:",notnull"`
	LastName   string    `json:"last_name" pg:",notnull"`
	Email      string    `json:"email" pg:",notnull,unique"`
	Password   string    `json:"password" pg:",notnull"`
	Avatar     string    `json:"avatar"`
	Status     int       `json:"status"`
	LastActive time.Time `json:"last_active" pg:"type:timestamp without time zone,default:now()"`
	CreateUser int       `json:"create_user"`
	UpdateUser int       `json:"update_user"`
	CreateTime time.Time `json:"create_time" pg:"type:timestamp without time zone,default:now()"`
	UpdateTime time.Time `json:"update_time" pg:"type:timestamp without time zone,default:now()"`
	DeleteFlag bool      `json:"delete_flag"`
}

func ConvertToSysUser(input InputCreateUser) SysUser {
	return SysUser{
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Email:      input.Email,
		Password:   common.HashPassword(input.Password),
		Avatar:     input.Avatar,
		Status:     1,
		CreateTime: common.GetDateTimeNow(),
		UpdateTime: common.GetDateTimeNow(),
	}
}
