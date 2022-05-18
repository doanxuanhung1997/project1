package model

import (
	"houze_ops_backend/helpers/common"
	"houze_ops_backend/helpers/constants"
	"time"
)

//GenBlog Model
type GenBlog struct {
	tableName struct{} `pg:"public.gen_blog"`

	Id               int       `json:"id" pg:"type:serial,pk"`
	Title            string    `json:"title" pg:",notnull"`
	Author           string    `json:"author"`
	MetaDescription  string    `json:"meta_description"`
	Tag              string    `json:"tag"`
	Content          string    `json:"content"`
	Url              string    `json:"url"`
	ApprovalStatus   int       `json:"approval_status"`
	PublishTime      time.Time `json:"publish_time" pg:"type:timestamp without time zone,default:now()"`
	Image            string    `json:"image"`
	ImageDescription string    `json:"image_description"`
	LinkUrl          string    `json:"link_url"`
	LinkName         string    `json:"link_name"`
	LinkIsOpenNewTab bool      `json:"link_is_open_new_tab"`
	CreateUser       int       `json:"create_user"`
	UpdateUser       int       `json:"update_user"`
	CreateTime       time.Time `json:"create_time" pg:"type:timestamp without time zone,default:now()"`
	UpdateTime       time.Time `json:"update_time" pg:"type:timestamp without time zone,default:now()"`
}

//GenCategory Model
type GenCategory struct {
	tableName struct{} `pg:"public.gen_category"`

	Id         int       `json:"id" pg:"type:serial,pk"`
	Name       string    `json:"name" pg:",notnull"`
	CreateUser int       `json:"create_user"`
	UpdateUser int       `json:"update_user"`
	CreateTime time.Time `json:"create_time" pg:"type:timestamp without time zone,default:now()"`
	UpdateTime time.Time `json:"update_time" pg:"type:timestamp without time zone,default:now()"`
}

//GenBlogCategory Model
type GenBlogCategory struct {
	tableName struct{} `pg:"public.gen_blog_category"`

	BlogId     int       `json:"blog_id" pg:",pk"`
	CategoryId int       `json:"CategoryId" pg:",pk"`
	CreateUser int       `json:"create_user"`
	UpdateUser int       `json:"update_user"`
	CreateTime time.Time `json:"create_time" pg:"type:timestamp without time zone,default:now()"`
	UpdateTime time.Time `json:"update_time" pg:"type:timestamp without time zone,default:now()"`
}

func ConvertToGenBlog(input InputCreateBlog, userId int) *GenBlog {
	publishTime, _ := time.Parse(constants.DateTimeFormat, input.PublishTime)
	return &GenBlog{
		Title:            input.Title,
		Author:           input.Author,
		MetaDescription:  input.MetaDescription,
		Tag:              input.Tag,
		Content:          input.Content,
		Url:              input.Url,
		ApprovalStatus:   input.ApprovalStatus,
		PublishTime:      publishTime,
		Image:            input.Image,
		ImageDescription: input.ImageDescription,
		LinkName:         input.LinkName,
		LinkUrl:          input.LinkUrl,
		LinkIsOpenNewTab: input.LinkIsOpenNewTab,
		CreateUser:       userId,
		UpdateUser:       userId,
		CreateTime:       common.GetDateTimeNow(),
		UpdateTime:       common.GetDateTimeNow(),
	}
}

func ConvertToGenBlogCategory(arrCategory []int, userId int) []*GenBlogCategory {
	var arr []*GenBlogCategory
	for _, c := range arrCategory {
		arr = append(arr, &GenBlogCategory{
			CategoryId: c,
			CreateUser: userId,
			UpdateUser: userId,
			CreateTime: common.GetDateTimeNow(),
			UpdateTime: common.GetDateTimeNow(),
		})
	}
	return arr
}

func ConvertToGenBlogUpdate(blog GenBlog, input InputUpdateBlog, userId int) GenBlog {
	publicTime, _ := time.Parse(constants.DateTimeFormat, input.PublishTime)
	blog.Title = input.Title
	blog.Author = input.Author
	blog.ApprovalStatus = input.ApprovalStatus
	blog.MetaDescription = input.MetaDescription
	blog.Tag = input.Tag
	blog.Url = input.Url
	blog.PublishTime = publicTime
	blog.Image = input.Image
	blog.ImageDescription = input.ImageDescription
	blog.LinkName = input.LinkName
	blog.LinkUrl = input.LinkUrl
	blog.LinkIsOpenNewTab = input.LinkIsOpenNewTab
	blog.UpdateTime = common.GetDateTimeNow()
	blog.UpdateUser = userId
	return blog
}

func ConvertToGenBlogCategoryUpdate(blogId int, arrCategory []int, userId int) []*GenBlogCategory {
	var arr []*GenBlogCategory
	for _, c := range arrCategory {
		arr = append(arr, &GenBlogCategory{
			BlogId:     blogId,
			CategoryId: c,
			CreateUser: userId,
			UpdateUser: userId,
			CreateTime: common.GetDateTimeNow(),
			UpdateTime: common.GetDateTimeNow(),
		})
	}
	return arr
}
