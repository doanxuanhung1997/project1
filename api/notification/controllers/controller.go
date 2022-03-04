package controllers

import (
	"encoding/json"
	"net/http"
	listenerRepository "sandexcare_backend/api/listener/repository"
	"sandexcare_backend/api/notification/model"
	"sandexcare_backend/api/notification/repository"
	userRepository "sandexcare_backend/api/user/repository"
	"sandexcare_backend/helpers/common"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/message"
	"sandexcare_backend/helpers/middlewares"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Create Notification
// @Schemes
// @Description API tạo thống báo.
// @Param input body model.InputCreateNotification true "Input Create Notification"
// @Tags Notification
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /notification [post]
func CreateNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorTokenInvalid,
			})
			return
		}
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputCreateNotification{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate content
		if common.IsEmpty(input.Content) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorContentEmpty,
			})
			return
		}

		// validate receiver_id
		if common.IsEmpty(input.ReceiverId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorListenerIdEmpty,
			})
			return
		}

		if input.ReceiverRole == constant.RoleUser {
			_, errUser := userRepository.PublishInterfaceUser().GetUserById(input.ReceiverId)
			if errUser != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorReceiverNotFound,
				})
				return
			}
		} else {
			_, errListener := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(input.ReceiverId)
			if errListener != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorReceiverNotFound,
				})
				return
			}
		}

		// init model
		var notify model.Notification
		notify.Id = primitive.NewObjectID()
		notify.ReceiverId = input.ReceiverId
		notify.ReceiverRole = input.ReceiverRole
		notify.Content = input.Content
		notify.Type = input.Type
		notify.CreatedAt = time.Now().UTC()
		notify.UpdatedAt = time.Now().UTC()

		// handle create data notification
		errCreateNotify := repository.PublishInterfaceNotification().CreateNotification(notify)
		if errCreateNotify == nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errCreateNotify.Error(),
			})
			return
		}
	}
}

// @Summary
// @Schemes
// @Description API Get list notify of user token.
// @Tags Notification
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetListNotification
// @Failure 401 {object} model.ResponseError
// @Failure 400 {object} model.ResponseError
// @Router  /notification [get]
func GetListNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get token info
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorTokenInvalid,
			})
			return
		}
		//get list notify by id and role in token
		data := repository.PublishInterfaceNotification().GetListNotification(tokenInfo.Id, tokenInfo.Role)

		c.JSON(http.StatusOK, model.ResponseGetListNotification{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    data,
		})
		return
	}
}

// @Summary
// @Schemes
// @Description API Read notify
// @Param input body model.InputReadNotification true "Input read notify"
// @Tags Notification
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /notification/read [post]
func ReadNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get token info
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorTokenInvalid,
			})
			return
		}
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputReadNotification{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		//read notification
		errRead := repository.PublishInterfaceNotification().ReadNotification(input.Id, input.ReadAll, tokenInfo.Id, tokenInfo.Role)
		if errRead == nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errRead.Error(),
			})
			return
		}
	}
}
