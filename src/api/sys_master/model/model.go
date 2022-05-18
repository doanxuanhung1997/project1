package model

type SysProductType struct {
	tableName struct{} `pg:"public.sys_product_type"`

	Id   int    `json:"id" pg:"type:serial,pk"`
	Name string `json:"name" pg:",notnull"`
}

type SysProvince struct {
	tableName struct{} `pg:"public.sys_province"`

	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
	NameWithType string `json:"name_with_type"`
	Code         string `json:"code"`
}

type SysDistrict struct {
	tableName struct{} `pg:"public.sys_district"`

	Name         string `json:"name"`
	Type         string `json:"type"`
	Slug         string `json:"slug"`
	NameWithType string `json:"name_with_type"`
	Path         string `json:"path"`
	PathWithType string `json:"path_with_type"`
	Code         string `json:"code"`
	ParentCode   string `json:"parent_code"`
}

type SysWards struct {
	tableName struct{} `pg:"public.sys_wards"`

	Name         string `json:"name"`
	Type         string `json:"type"`
	Slug         string `json:"slug"`
	NameWithType string `json:"name_with_type"`
	Path         string `json:"path"`
	PathWithType string `json:"path_with_type"`
	Code         string `json:"code"`
	ParentCode   string `json:"parent_code"`
}
