package model

import (
	"houze_ops_backend/helpers/constants"
	"strings"
)

type ResultBlog struct {
	Id               int                   `json:"id"`
	Category         []*ResultBlogCategory `json:"category"`
	Title            string                `json:"title"`
	Author           string                `json:"author"`
	MetaDescription  string                `json:"meta_description"`
	Tag              []string              `json:"tag"`
	Content          string                `json:"content"`
	Url              string                `json:"url"`
	ApprovalStatus   int                   `json:"approval_status"`
	PublishTime      string                `json:"publish_time"`
	Image            string                `json:"image"`
	ImageDescription string                `json:"image_description"`
	LinkUrl          string                `json:"link_url"`
	LinkName         string                `json:"link_name"`
	LinkIsOpenNewTab bool                  `json:"link_is_open_new_tab"`
	CreateUser       int                   `json:"create_user"`
	UpdateUser       int                   `json:"update_user"`
	CreateTime       string                `json:"create_time"`
	UpdateTime       string                `json:"update_time"`
}

type ResultBlogCategory struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func ConvertToResultBlog(blog GenBlog) *ResultBlog {
	return &ResultBlog{
		Id:               blog.Id,
		Title:            blog.Title,
		Author:           blog.Author,
		MetaDescription:  blog.MetaDescription,
		Tag:              strings.Split(blog.Tag, ","),
		Content:          blog.Content,
		Url:              blog.Url,
		ApprovalStatus:   blog.ApprovalStatus,
		PublishTime:      blog.PublishTime.Format(constants.DateTimeFormat),
		Image:            blog.Image,
		ImageDescription: blog.ImageDescription,
		LinkName:         blog.LinkName,
		LinkUrl:          blog.LinkUrl,
		LinkIsOpenNewTab: blog.LinkIsOpenNewTab,
		CreateTime:       blog.CreateTime.Format(constants.DateTimeFormat),
		UpdateTime:       blog.UpdateTime.Format(constants.DateTimeFormat),
	}
}
