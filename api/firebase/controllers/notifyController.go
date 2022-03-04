package controllers

import (
	"encoding/json"
	"net/http"

	"sandexcare_backend/api/firebase/model"
	"sandexcare_backend/api/user/repository"
	"sandexcare_backend/helpers/message"
	"sandexcare_backend/helpers/middlewares"
	"sandexcare_backend/notification"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1
// @Summary
// @Schemes
// @Description Send Background Notification By Firebase Service -  Gửi tin nhắn khi ngoại tuyến và trực tuyến
// @Tags Firebase
// @Param   Message     body    model.Message     false        "Message: Tin nhắn gồm tựa đề, nội dung, hình ảnh và token của người gửi"
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /firebase/send [post]
func SendNotify() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawBody, _ := c.GetRawData()
		msg := model.Message{}
		err := json.Unmarshal(rawBody, &msg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.FirebaseInvalidFormat,
			})
			return
		}
		err = notification.SendNotifyFirebase(msg.Title, msg.Content, msg.Icon, msg.To)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.FirebaseError,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.FirebaseSuccess,
		})
	}
}

// @BasePath /api/v1
// @Summary
// @Schemes
// @Description Send Background Notification By Firebase Service BY UUID (TODO: Cần token service to service) -  Gửi tin nhắn khi ngoại tuyến và trực tuyến
// @Tags Firebase
// @Param uuid   path string false "UUID của user"
// @Param   Message     body    model.Message     false        "Message: Tin nhắn gồm tựa đề, nội dung, hình ảnh và token của người gửi"
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /firebase/send/{uuid} [post]
func SendNotifyByUUID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// rawBody, _ := c.GetRawData()
		// msg := model.Message{}
		// err := json.Unmarshal(rawBody, &msg)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"code":    http.StatusBadRequest,
		// 		"message": message.FirebaseInvalidFormat,
		// 	})
		// 	return
		// }
		// err = notification.SendNotifyFirebase(msg.Title, msg.Content, msg.Icon, msg.To)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"code":    http.StatusBadRequest,
		// 		"message": message.FirebaseError,
		// 	})
		// 	return
		// }
		// c.JSON(http.StatusOK, gin.H{
		// 	"code":    http.StatusOK,
		// 	"message": message.FirebaseSuccess,
		// })
	}
}

// @Summary
// @Schemes
// @Description Register Device Token - Đăng kí token device
// @Tags Firebase
// @Param Authorization header string true "Token lấy từ API authen có prefix Token"
// @Param   Token  body    model.InputToken     false        "Token của người dùng lấy từ firebase, khi gọi lại API này, token cũ sẽ bị ghi đè (Xoá vĩnh viễn)"
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /firebase/register [post]
func RegisterNotify() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, err := middlewares.AuthenticateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorTokenInvalid,
			})
			return
		}
		rawBody, _ := c.GetRawData()
		token := model.InputToken{}
		err = json.Unmarshal(rawBody, &token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.FirebaseInvalidFormat,
			})
			return
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.FirebaseRegisterInvalid,
			})
			return
		}
		err = repository.PublishInterfaceUser().RegisterFirebase(tokenInfo.Id, token.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.FirebaseRegisterInvalid,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.FirebaseRegisterSuccess,
		})
	}
}
