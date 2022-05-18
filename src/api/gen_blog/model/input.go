package model

// InputCreateBlog params data create blog
type InputCreateBlog struct {
	Title            string `validate:"required" json:"title"`
	MetaDescription  string `json:"meta_description"`
	Category         []int  `json:"category" validate:"required"`
	Author           string `json:"author" validate:"required"`
	Content          string `json:"content" validate:"required"`
	Tag              string `json:"tag" validate:"required"`
	ApprovalStatus   int    `json:"approval_status" validate:"numeric,min=1,max=4"`
	Url              string `json:"url" validate:"required,url"`
	PublishTime      string `json:"publish_time" validate:"required"`
	Image            string `json:"image" validate:"base64"`
	ImageDescription string `json:"image_description"`
	LinkUrl          string `json:"link_url"`
	LinkName         string `json:"link_name"`
	LinkIsOpenNewTab bool   `json:"link_is_open_new_tab"`
}

type InputUpdateBlog struct {
	Id               int    `json:"id"`
	Title            string `json:"title"`
	MetaDescription  string `json:"meta_description"`
	Category         []int  `json:"category"`
	Author           string `json:"author"`
	Content          string `json:"content"`
	ApprovalStatus   int    `json:"approval_status"`
	Tag              string `json:"tag"`
	Url              string `json:"url"`
	PublishTime      string `json:"publish_time"`
	Image            string `json:"image"`
	ImageDescription string `json:"image_description"`
	LinkUrl          string `json:"link_url"`
	LinkName         string `json:"link_name"`
	LinkIsOpenNewTab bool   `json:"link_is_open_new_tab"`
}

type InputDeleteBlog struct {
	Id int `json:"id"`
}
