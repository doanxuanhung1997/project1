package repository

import (
	"houze_ops_backend/api/sys_master/model"
	"houze_ops_backend/db"
)

type masterInterface interface {
	GetProvince() (province []model.SysProvince)
	GetDistrict(parentCode string) (district []model.SysDistrict)
	GetWards(parentCode string) (wards []model.SysWards)
}

func PublishInterface() masterInterface {
	return &resource{}
}

type resource struct {
}

func (r *resource) GetProvince() (province []model.SysProvince) {
	_ = db.GetConnectionDB().Model(&province).Select()
	return province
}

func (r *resource) GetDistrict(parentCode string) (district []model.SysDistrict) {
	_ = db.GetConnectionDB().Model(&district).Where("parent_code = ?", parentCode).Select()
	return district
}

func (r *resource) GetWards(parentCode string) (wards []model.SysWards) {
	_ = db.GetConnectionDB().Model(&wards).Where("parent_code = ?", parentCode).Select()
	return wards
}
