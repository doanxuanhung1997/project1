package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"houze_ops_backend/api/gen_blog/model"
	genBlogRepo "houze_ops_backend/api/gen_blog/repository"
	"houze_ops_backend/helpers/message"
	"houze_ops_backend/helpers/middlewares"
	"net/http"
	"strconv"
)

func GetAllCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		category := genBlogRepo.PublishInterface().GetAllCategory()
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
			"data":    category,
		})
	}
}

func GetAllBlogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, errAuth := middlewares.AuthenticateToken(c)
		if errAuth != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		approvalStatus := c.Query("approval_status")
		status, _ := strconv.Atoi(approvalStatus)
		author := c.Query("author")
		category := c.Query("category")
		categoryId, _ := strconv.Atoi(category)
		name := c.Query("name")

		blogs := genBlogRepo.PublishInterface().GetAllBlogs(status, author, categoryId, name)

		var resultBlogs []*model.ResultBlog
		for _, b := range blogs {
			resultBlog := model.ConvertToResultBlog(b)
			resultBlog.Category = genBlogRepo.PublishInterface().GetBlogCategory(resultBlog.Id)
			resultBlogs = append(resultBlogs, resultBlog)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
			"data":    resultBlogs,
		})
	}
}

func GetDetailBlog() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, errAuth := middlewares.AuthenticateToken(c)
		if errAuth != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		id := c.Query("id")
		blogId, _ := strconv.Atoi(id)
		blog, err := genBlogRepo.PublishInterface().GetDetailBlog(blogId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorNotFound,
			})
			return
		}
		resultBlog := model.ConvertToResultBlog(blog)
		resultBlog.Category = genBlogRepo.PublishInterface().GetBlogCategory(resultBlog.Id)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
			"data":    resultBlog,
		})
		return
	}
}

func CreateBlog() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, errAuth := middlewares.AuthenticateToken(c)
		if errAuth != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		// get param from body
		rawBody, _ := c.GetRawData()
		input := model.InputCreateBlog{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		genBlog := model.ConvertToGenBlog(input, tokenInfo.Id)
		genBlogCategory := model.ConvertToGenBlogCategory(input.Category, tokenInfo.Id)

		id, err := genBlogRepo.PublishInterface().CreateBlog(genBlog, genBlogCategory)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageFail,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
			"data":    id,
		})
		return
	}
}

func UpdateBlog() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, errAuth := middlewares.AuthenticateToken(c)
		if errAuth != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		// get param from body
		rawBody, _ := c.GetRawData()
		input := model.InputUpdateBlog{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		blog, err := genBlogRepo.PublishInterface().GetBlogById(input.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorNotFound,
			})
			return
		}
		blog = model.ConvertToGenBlogUpdate(blog, input, tokenInfo.Id)
		blogCategory := model.ConvertToGenBlogCategoryUpdate(blog.Id, input.Category, tokenInfo.Id)

		err = genBlogRepo.PublishInterface().UpdateBlog(&blog, blogCategory)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageFail,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
		})
		return
	}
}

func DeleteBlog() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, errAuth := middlewares.AuthenticateToken(c)
		if errAuth != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		// get param from body
		rawBody, _ := c.GetRawData()
		input := model.InputDeleteBlog{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}
		_, err = genBlogRepo.PublishInterface().GetBlogById(input.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorNotFound,
			})
			return
		}
		err = genBlogRepo.PublishInterface().DeleteBlog(input.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageFail,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
		})
		return
	}
}
