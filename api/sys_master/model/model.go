package model

type SysProductType struct {
	tableName struct{} `pg:"public.sys_product_type"`

	Id   int    `json:"id" pg:"type:serial,pk"`
	Name string `json:"name" pg:",notnull"`
}
