package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"houze_ops_backend/api/sys_user/model"
	user_repo "houze_ops_backend/api/sys_user/repository"
	"houze_ops_backend/helpers/common"
	"houze_ops_backend/helpers/message"
	"houze_ops_backend/helpers/middlewares"
	"net/http"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get param from body
		rawBody, _ := c.GetRawData()
		input := model.InputLogin{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		email := input.Email
		password := input.Password

		//hash password
		hashedPass := common.HashPassword(password)
		user, err := user_repo.PublishInterfaceUser().Login(email, hashedPass)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageFail,
			})
			return
		}
		var token, _ = middlewares.GenerateJWT(user.Id, user.Email, 1)
		data := model.LoginData{}
		data.Id = user.Id
		data.Token = token
		data.Email = user.Email

		c.JSON(http.StatusOK, model.ResponseLogin{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    data,
		})
		return
	}
}

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get param from body
		rawBody, _ := c.GetRawData()
		input := model.InputCreateUser{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		sysUser := model.ConvertToSysUser(input)
		err = user_repo.PublishInterfaceUser().CreateUser(sysUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusOK,
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
