package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"sandexcare_backend/api/master/model"
	"sandexcare_backend/api/master/repository"
	"sandexcare_backend/helpers/config"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/message"
	wsControllers "sandexcare_backend/server_websocket/controllers"
	"time"
)

//func CreateMasterData() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		rawBody, _ := c.GetRawData()
//		input := model.InputCreateMasterData{}
//		err := json.Unmarshal(rawBody, &input)
//		if err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{
//				"code":    http.StatusBadRequest,
//				"message": message.MessageErrorConvertInput,
//			})
//			return
//		}
//		if input.Table == "consulting_field" {
//			consultingField := model.ConsultingField{}
//			consultingField.Id = primitive.NewObjectID()
//			consultingField.Name = input.Value
//			consultingField.CreatedAt = time.Now().UTC()
//			consultingField.UpdatedAt = time.Now().UTC()
//			_ = repository.PublishInterfaceMaster().CreateConsultingField(consultingField)
//		}
//		if input.Table == "time_slot" {
//			timeSlot := model.TimeSlot{}
//			timeSlot.Id = primitive.NewObjectID()
//			timeSlot.TimeSlot = input.Value
//			timeSlot.BookingTime = input.BookingTime
//			timeSlot.CreatedAt = time.Now().UTC()
//			timeSlot.UpdatedAt = time.Now().UTC()
//			_ = repository.PublishInterfaceMaster().CreateTimeSlot(timeSlot)
//		}
//		c.JSON(http.StatusOK, gin.H{
//			"code":    http.StatusOK,
//			"message": message.MessageSuccess,
//		})
//	}
//}

// @Summary
// @Schemes
// @Description API get consulting field data
// @Tags Master
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetConsultingField
// @Router  /master/consulting-field [get]
func GetConsultingField() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := repository.PublishInterfaceMaster().GetAllConsultingField()
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
			"data":    data,
		})
	}
}

// @Summary
// @Schemes
// @Description API get time slot data
// @Tags Master
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetTimeSlot
// @Router  /master/time-slot [get]
func GetTimeSlot() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := repository.PublishInterfaceMaster().GetAllTimeSlot()
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
			"data":    data,
		})
	}
}

// @Summary
// @Schemes
// @Description API get config data
// @Tags Master
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetConfigData
// @Router  /master/config [get]
func GetConfigData() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := model.ConfigData{}
		env := config.GetEnvValue()
		data.AppointmentPrice = env.Call.AppointmentPrice
		c.JSON(http.StatusOK, model.ResponseGetConfigData{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    data,
		})
	}
}

// @Summary API send message event web socket
// @Schemes
// @Description Xử lý gửi tin nhắn đến 1 người qua websocket!
// @Param input body model.InputSendEventWebSocket true "Input Send Event Web Socket"
// @Tags Master
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /master/send-event-ws [post]
func SendEventWebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawBody, _ := c.GetRawData()
		inputWS := model.InputSendEventWebSocket{}
		err := json.Unmarshal(rawBody, &inputWS)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		hub := wsControllers.GetHub()
		payload := wsControllers.SocketEventStruct{}
		payload.EventName = inputWS.Event
		payload.EventPayload = map[string]interface{}{
			"message": inputWS.Message,
		}
		wsControllers.EmitToSpecificClient(hub, payload, inputWS.FromUserId)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
		})
		return
	}
}

// @Summary
// @Schemes
// @Description API get server datetime
// @Tags Master
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetServerDatetime
// @Router  /master/datetime [get]
func GetServerDatetime() gin.HandlerFunc {
	return func(c *gin.Context) {
		datetime := time.Now().UTC()
		c.JSON(http.StatusOK, model.ResponseGetServerDatetime{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    datetime.Format(constant.DateTimeFormat),
		})
	}
}
