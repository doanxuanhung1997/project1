package repository

import (
	"context"
	"github.com/go-pg/pg/v10"
	"houze_ops_backend/api/gen_blog/model"
	"houze_ops_backend/db"
	"houze_ops_backend/helpers/common"
)

type genBlogInterface interface {
	GetAllCategory() (category []model.GenCategory)
	CreateBlog(blog *model.GenBlog, blogCategory []*model.GenBlogCategory) (int, error)
	GetAllBlogs(approvalStatus int, author string, category int, name string) (blog []model.GenBlog)
	GetBlogCategory(blogId int) (rbc []*model.ResultBlogCategory)
	GetDetailBlog(id int) (blog model.GenBlog, err error)
	GetBlogById(id int) (blog model.GenBlog, err error)
	UpdateBlog(blog *model.GenBlog, blogCategory []*model.GenBlogCategory) (err error)
	DeleteBlog(id int) (err error)
}

func PublishInterface() genBlogInterface {
	return &resource{}
}

type resource struct {
}

func (r *resource) CreateBlog(blog *model.GenBlog, blogCategory []*model.GenBlogCategory) (id int, err error) {
	err = db.GetConnectionDB().RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		_, err = tx.Model(blog).Insert()
		if err != nil {
			return err
		}
		id = blog.Id
		for _, bc := range blogCategory {
			bc.BlogId = blog.Id
		}
		_, err = tx.Model(&blogCategory).Insert()
		return err
	})
	return id, err
}

func (r *resource) GetAllCategory() (category []model.GenCategory) {
	_ = db.GetConnectionDB().Model(&category).Order("id ASC").Select()
	return category
}

func (r *resource) GetAllBlogs(approvalStatus int, author string, category int, name string) (blog []model.GenBlog) {
	query := db.GetConnectionDB().Model(&blog).Join("JOIN gen_blog_category AS gbc ON gen_blog.id = gbc.blog_id").Distinct()
	if approvalStatus != 0 {
		query.Where("approval_status = ?", approvalStatus)
	}
	if !common.IsEmpty(author) {
		query.Where("author = ?", author)
	}
	if !common.IsEmpty(name) {
		query.Where("title LIKE ?", "%"+name+"%").
			WhereOr("content LIKE ?", "%"+name+"%")
	}
	if category != 0 {
		query.Where("category_id = ?", category)
	}
	_ = query.Order("id ASC").Select()

	return blog
}

func (r *resource) GetBlogCategory(blogId int) (rbc []*model.ResultBlogCategory) {
	_ = db.GetConnectionDB().Model(&model.GenBlogCategory{}).
		ColumnExpr("gen_category.id AS value, gen_category.name AS label").
		Join("JOIN gen_category ON gen_category.id = gen_blog_category.category_id").
		Where("blog_id = ?", blogId).
		Select(&rbc)
	return rbc

}

func (r *resource) GetDetailBlog(id int) (blog model.GenBlog, err error) {
	err = db.GetConnectionDB().Model(&blog).Where("id = ?", id).Select()
	return blog, err
}

func (r *resource) GetBlogById(id int) (blog model.GenBlog, err error) {
	err = db.GetConnectionDB().Model(&blog).Where("id = ?", id).Select()
	return blog, err
}

func (r *resource) UpdateBlog(blog *model.GenBlog, blogCategory []*model.GenBlogCategory) (err error) {
	err = db.GetConnectionDB().RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		_, err = tx.Model(blog).WherePK().Update()
		if err != nil {
			return err
		}

		_, err = tx.Model(&model.GenBlogCategory{}).Where("blog_id = ?", blog.Id).Delete()
		if err != nil {
			return err
		}

		_, err = tx.Model(&blogCategory).Insert()
		return err
	})
	return err
}

func (r *resource) DeleteBlog(id int) (err error) {
	err = db.GetConnectionDB().RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		_, err = tx.Model(&model.GenBlog{}).Where("id = ?", id).Delete()
		if err != nil {
			return err
		}

		_, err = tx.Model(&model.GenBlogCategory{}).Where("blog_id = ?", id).Delete()
		return err
	})
	return err
}
