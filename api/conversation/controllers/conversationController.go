package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"net/http"
	"sandexcare_backend/agora"
	"sandexcare_backend/api/conversation/model"
	"sandexcare_backend/api/conversation/repository"
	listenerRepository "sandexcare_backend/api/listener/repository"
	masterModel "sandexcare_backend/api/master/model"
	paymentRepository "sandexcare_backend/api/payment/repository"
	trackEventControllers "sandexcare_backend/api/track_event/controllers"
	userRepository "sandexcare_backend/api/user/repository"
	"sandexcare_backend/helpers/common"
	"sandexcare_backend/helpers/config"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/message"
	"sandexcare_backend/helpers/middlewares"
	wsControllers "sandexcare_backend/server_websocket/controllers"
	"time"
)

// @Summary
// @Schemes
// @Description API user start call of appointment booked
// @Param input body model.InputStartCall true "Input start call"
// @Tags Conversation
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseStartConversation
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /conversation/start-call [post]
func StartConversation() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		input := model.InputStartCall{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// get order data
		orderInfo, errOrder := paymentRepository.PublishInterfacePayment().GetOrderPaymentById(input.OrderId)
		if errOrder != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errOrder.Error(),
			})
			return
		}

		// order valid
		if orderInfo.UserId != tokenInfo.Id || orderInfo.Status != constant.Active {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}

		endCall := orderInfo.CallDatetime.Add(time.Hour*time.Duration(-7) + time.Minute*time.Duration(30))
		timeNow, _ := time.Parse(constant.DateTimeFormat, time.Now().UTC().Format(constant.DateTimeFormat))

		if timeNow.Before(orderInfo.CallDatetime.Add(time.Hour * time.Duration(-7))) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorStartConversation,
			})
			return
		} else if endCall.Before(timeNow) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorStartConversation,
			})
			return
		}

		hub := wsControllers.GetHub()
		listenerOnCall := wsControllers.GetAllListenersReadyCall(hub)
		onFlag := false
		for _, v := range listenerOnCall {
			if v.ListenerId == orderInfo.ListenerId {
				onFlag = true
				break
			}
		}

		// init model
		var callHistory model.CallHistory
		callHistory.Id = primitive.NewObjectID()
		callHistory.UserId = tokenInfo.Id
		callHistory.Channel = common.GenerateNumber(6)
		callHistory.ListenerId = orderInfo.ListenerId
		callHistory.OrderId = input.OrderId
		callHistory.CreatedAt = time.Now().UTC()
		callHistory.UpdatedAt = time.Now().UTC()

		if !onFlag {
			// Block chuyên viên
			// TODO Tạm thời comment code block chuyên viên
			//listenerInfo, _ := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(orderInfo.ListenerId)
			//listenerInfo.Status = constant.Inactive
			//_ = listenerRepository.PublishInterfaceListener().UpdateListener(listenerInfo)
			//
			//subject := mail.GetSubject(mail.TypeNotifyBlockListener)
			//content := mail.GetHtmlContentBlockListener()
			//_, _ = mail.SendEmail(tokenInfo.Email, subject, content)

			callHistory.Status = constant.Unconnected
			_ = repository.PublishInterfaceConversation().CreateCallHistory(callHistory)

			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": message.MessageErrorListenerUnconnect,
			})
			return
		}

		callHistory.Status = constant.Started
		// handle create data call history
		errCallHistory := repository.PublishInterfaceConversation().CreateCallHistory(callHistory)
		if errCallHistory == nil {
			env := config.GetEnvValue()
			userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(tokenInfo.Id)
			uidListener := common.GenerateNumber(6)
			uidUser := common.GenerateNumber(6)
			payload := wsControllers.SocketEventStruct{}
			payload.EventName = constant.WSEventStartCall
			payload.EventPayload = map[string]interface{}{
				"order_id":     input.OrderId,
				"call_id":      callHistory.Id.Hex(),
				"username":     userInfo.Name,
				"date":         orderInfo.Date.Format(constant.DateFormat),
				"booking_time": orderInfo.BookingTime,
				"token":        agora.GenerateRtcToken(uidListener, callHistory.Channel),
				"channel":      callHistory.Channel,
				"uid":          uidListener,
				"ringing_time": env.Call.RingingTimeAppointment,
			}
			wsControllers.EmitToSpecificClient(hub, payload, orderInfo.ListenerId)

			responseData := model.DataStartConversation{
				OrderId:     input.OrderId,
				CallId:      callHistory.Id.Hex(),
				Token:       agora.GenerateRtcToken(uidUser, callHistory.Channel),
				Channel:     callHistory.Channel,
				Uid:         uidUser,
				RingingTime: env.Call.RingingTimeAppointment,
			}
			c.JSON(http.StatusOK, model.ResponseStartConversation{
				Code:    http.StatusOK,
				Message: message.MessageSuccess,
				Data:    responseData,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errCallHistory.Error(),
			})
			return
		}
	}
}

// @Summary
// @Schemes
// @Description API bắt đầu cuộc gọi với chuyên gia
// @Param input body model.InputStartCall true "Input start call"
// @Tags Conversation
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseStartConversation
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /conversation/start-call/expert [post]
func StartCallExpert() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		input := model.InputStartCall{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// get order data
		orderInfo, errOrder := paymentRepository.PublishInterfacePayment().GetOrderPaymentById(input.OrderId)
		if errOrder != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errOrder.Error(),
			})
			return
		}

		// order valid
		if orderInfo.UserId != tokenInfo.Id || orderInfo.Status != constant.Active {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}

		endCall := orderInfo.CallDatetime.Add(time.Hour*time.Duration(-7) + time.Minute*time.Duration(120))
		timeNow, _ := time.Parse(constant.DateTimeFormat, time.Now().UTC().Format(constant.DateTimeFormat))

		if timeNow.Before(orderInfo.CallDatetime.Add(time.Hour * time.Duration(-7))) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorStartConversation,
			})
			return
		} else if endCall.Before(timeNow) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorStartConversation,
			})
			return
		}

		// init model
		var callHistory model.CallHistory
		callHistory.Id = primitive.NewObjectID()
		callHistory.UserId = tokenInfo.Id
		callHistory.Status = constant.Started
		callHistory.Channel = common.GenerateNumber(6)
		callHistory.ListenerId = orderInfo.ListenerId
		callHistory.OrderId = input.OrderId
		callHistory.CreatedAt = time.Now().UTC()
		callHistory.UpdatedAt = time.Now().UTC()

		// handle create data call history
		errCallHistory := repository.PublishInterfaceConversation().CreateCallHistory(callHistory)
		if errCallHistory == nil {
			hub := wsControllers.GetHub()
			env := config.GetEnvValue()
			userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(tokenInfo.Id)
			uidExpert := common.GenerateNumber(6)
			uidUser := common.GenerateNumber(6)
			payload := wsControllers.SocketEventStruct{}
			payload.EventName = constant.WSEventStartCall
			payload.EventPayload = map[string]interface{}{
				"order_id":     input.OrderId,
				"call_id":      callHistory.Id.Hex(),
				"username":     userInfo.Name,
				"date":         orderInfo.Date.Format(constant.DateFormat),
				"booking_time": orderInfo.BookingTime,
				"token":        agora.GenerateRtcToken(uidExpert, callHistory.Channel),
				"channel":      callHistory.Channel,
				"uid":          uidExpert,
				"ringing_time": env.Call.RingingTimeAppointment,
			}
			wsControllers.EmitToSpecificClient(hub, payload, orderInfo.ListenerId)

			responseData := model.DataStartConversation{
				OrderId:     input.OrderId,
				CallId:      callHistory.Id.Hex(),
				Token:       agora.GenerateRtcToken(uidUser, callHistory.Channel),
				Channel:     callHistory.Channel,
				Uid:         uidUser,
				RingingTime: env.Call.RingingTimeAppointment,
			}
			c.JSON(http.StatusOK, model.ResponseStartConversation{
				Code:    http.StatusOK,
				Message: message.MessageSuccess,
				Data:    responseData,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errCallHistory.Error(),
			})
			return
		}
	}
}

// @Summary
// @Schemes
// @Description API user end call
// @Param input body model.InputEndCall true "Input end call"
// @Tags Conversation
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /conversation/end-call [post]
func EndCall() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		input := model.InputEndCall{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// handle get data call history by id
		callData, errCall := repository.PublishInterfaceConversation().GetCallHistoryById(input.CallId)
		if errCall == nil {
			if callData.UserId == tokenInfo.Id {
				if callData.Status == constant.Unconnected || callData.Status == constant.Completed {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorPermissionDenied,
					})
					return
				}
				hub := wsControllers.GetHub()
				payload := wsControllers.SocketEventStruct{}
				payload.EventName = constant.WSEventEndCall
				payload.EventPayload = map[string]interface{}{
					"call_id": input.CallId,
				}
				wsControllers.EmitToSpecificClient(hub, payload, callData.ListenerId)
				callData.Status = constant.Completed
				callData.UpdatedAt = time.Now().UTC()
				_ = repository.PublishInterfaceConversation().UpdateCallHistory(callData)

				paymentData, _ := paymentRepository.PublishInterfacePayment().GetOrderPaymentById(callData.OrderId)
				paymentData.Status = constant.Completed
				paymentData.UpdatedAt = time.Now().UTC()
				_ = paymentRepository.PublishInterfacePayment().UpdateOrderPayment(paymentData)
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": message.MessageSuccess,
				})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorPermissionDenied,
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorCallIdInvalid,
			})
			return
		}
	}
}

// @Summary
// @Schemes
// @Description API listener join call
// @Param input body model.InputJoinCall true "Input listener join call"
// @Tags Conversation
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /conversation/join-call [post]
func JoinConversation() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		input := model.InputJoinCall{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// check token is listener
		if tokenInfo.Role == constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}

		// handle get data call history by id
		callData, errCall := repository.PublishInterfaceConversation().GetCallHistoryById(input.CallId)
		if errCall == nil {
			if callData.ListenerId == tokenInfo.Id {
				if callData.Status == constant.Unconnected || callData.Status == constant.Completed {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorPermissionDenied,
					})
					return
				}
				listenerInfo, _ := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(tokenInfo.Id)
				callData.Status = constant.Talking
				callData.UpdatedAt = time.Now().UTC()
				_ = repository.PublishInterfaceConversation().UpdateCallHistory(callData)
				trackEventControllers.CreateTrackCall(listenerInfo.Id.Hex(), callData.Id.Hex(), constant.JoinCall)
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": message.MessageSuccess,
				})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorPermissionDenied,
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorCallIdInvalid,
			})
			return
		}
	}
}

// @Summary
// @Schemes
// @Description API user switch listener
// @Param input body model.InputSwitchListener true "Input user switch listener"
// @Tags Conversation
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /conversation/switch-listener [post]
func SwitchListener() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		input := model.InputSwitchListener{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}
		orderInfo, errOrder := paymentRepository.PublishInterfacePayment().GetOrderPaymentById(input.OrderId)
		if errOrder == nil {
			// check token is user
			if tokenInfo.Id != orderInfo.UserId {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorPermissionDenied,
				})
				return
			}

			hub := wsControllers.GetHub()
			listenerOnCall := wsControllers.GetAllListenersReadyCall(hub)
			var listenerId string
			for _, v := range listenerOnCall {
				if v.ActiveCallNow {
					isCalled := false
					callHistory := repository.PublishInterfaceConversation().GetCallHistoryOfOrder(input.OrderId)
					for _, c := range callHistory {
						if v.ListenerId == c.ListenerId {
							isCalled = true
							break
						}
					}
					if !isCalled {
						listenerId = v.ListenerId
						break
					}
				}
			}
			if common.IsEmpty(listenerId) {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorListenerUnconnect,
				})
				return
			}

			// init model conversation ==> start call
			var callHistory model.CallHistory
			callHistory.Id = primitive.NewObjectID()
			callHistory.UserId = tokenInfo.Id
			callHistory.ListenerId = listenerId
			callHistory.Channel = common.GenerateNumber(6)
			callHistory.OrderId = input.OrderId
			callHistory.Status = constant.Started
			callHistory.CreatedAt = time.Now().UTC()
			callHistory.UpdatedAt = time.Now().UTC()

			errCallHistory := repository.PublishInterfaceConversation().CreateCallHistory(callHistory)
			if errCallHistory == nil {
				env := config.GetEnvValue()
				userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(tokenInfo.Id)
				uidListener := common.GenerateNumber(6)
				payload := wsControllers.SocketEventStruct{}
				payload.EventName = constant.WSEventStartCall
				payload.EventPayload = map[string]interface{}{
					"call_id":      callHistory.Id.Hex(),
					"username":     userInfo.Name,
					"date":         orderInfo.Date.Format(constant.DateFormat),
					"booking_time": orderInfo.BookingTime,
					"token":        agora.GenerateRtcToken(uidListener, callHistory.Channel),
					"channel":      callHistory.Channel,
					"uid":          uidListener,
					"ringing_time": env.Call.RingingTimeCallNow,
				}
				wsControllers.EmitToSpecificClient(hub, payload, listenerId)
			}
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errOrder.Error(),
			})
			return
		}
	}
}

// @Summary
// @Schemes
// @Description Get history call of user by token listener.
// @Tags Conversation
// @Accept json
// @Produce json
// @Param  user_id path string true "User Id."
// @Success 200 {object} model.ResponseGetCallHistoryUser
// @Failure 400 {object} model.ResponseError
// @Router  /conversation/call-history/user [get]
func GetCallHistoryUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorTokenInvalid,
			})
			return
		}
		userId := c.Query("user_id")
		if tokenInfo.Role == constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		// handle get data call history of user
		data := repository.PublishInterfaceConversation().GetCallHistoryOfUser(userId)

		c.JSON(http.StatusOK, model.ResponseGetCallHistoryUser{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    data,
		})
		return

	}
}

// @Summary
// @Schemes
// @Description Get history call of listener.
// @Tags Conversation
// @Accept json
// @Produce json
// @Param  status path string false "Ex: status = COMPLETED. Status of call. If there is no status param => get all status."
// @Success 200 {object} model.ResponseGetCallHistoryListener
// @Failure 400 {object} model.ResponseError
// @Router  /conversation/call-history/listener [get]
func GetCallHistoryListener() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorTokenInvalid,
			})
			return
		}
		if tokenInfo.Role == constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}

		status := c.Query("status")

		// handle get data call history of listener
		data := repository.PublishInterfaceConversation().GetCallHistoryOfListener(tokenInfo.Id, status)

		c.JSON(http.StatusOK, model.ResponseGetCallHistoryListener{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    data,
		})

	}
}

// @Summary
// @Schemes
// @Description Get detail call history of listener.
// @Tags Conversation
// @Accept json
// @Produce json
// @Param  call_id path string true "Call Id."
// @Success 200 {object} model.ResponseGetDetailCallHistoryListener
// @Failure 400 {object} model.ResponseError
// @Router  /conversation/call-history/listener/detail [get]
func GetDetailCallHistoryForListener() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorTokenInvalid,
			})
			return
		}
		if tokenInfo.Role == constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}

		callId := c.Query("call_id")
		// handle get data call history for listener
		dataQuery, err := repository.PublishInterfaceConversation().GetCallHistoryById(callId)
		if err == nil {
			var responseData model.DataGetCallHistoryListener
			responseData.Content = dataQuery.Content
			responseData.ConsultingField = dataQuery.ConsultingField
			responseData.StartCall = dataQuery.StartCall.Format(constant.DateTimeFormat)
			responseData.EndCall = dataQuery.EndCall.Format(constant.DateTimeFormat)
			responseData.CallId = dataQuery.Id.Hex()
			responseData.Status = dataQuery.Status
			responseData.UserId = dataQuery.UserId
			user, errUser := userRepository.PublishInterfaceUser().GetUserById(dataQuery.UserId)
			if errUser == nil {
				responseData.UserName = user.Name
				responseData.PhoneNumber = user.PhoneNumber
			}
			c.JSON(http.StatusOK, model.ResponseGetDetailCallHistoryListener{
				Code:    http.StatusOK,
				Message: message.MessageSuccess,
				Data:    responseData,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorCallIdInvalid,
			})
		}
	}
}

// @Summary
// @Schemes
// @Description API listener submit information conversation
// @Param input body model.InputSubmitInfoConversation true "Input listener submit information conversation"
// @Tags Conversation
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /conversation/info-call [post]
func SubmitInfoConversation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token info
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
		input := model.InputSubmitInfoConversation{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate CallId
		if common.IsEmpty(input.CallId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorCallIdInvalid,
			})
			return
		}

		// validate StartCall
		if common.IsEmpty(input.StartCall) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateEmpty,
			})
			return
		}

		// convert start call from string to datetime format
		StartTime, err := time.Parse(constant.DateTimeFormat, input.StartCall)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateFormat,
			})
			return
		}

		// validate EndCall
		if common.IsEmpty(input.EndCall) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateEmpty,
			})
			return
		}

		// convert end call from string to datetime format
		EndTime, err := time.Parse(constant.DateTimeFormat, input.EndCall)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateFormat,
			})
			return
		}

		// get data call by call_id
		dataCall, errCall := repository.PublishInterfaceConversation().GetCallHistoryById(input.CallId)
		if errCall == nil {
			//check permission
			if tokenInfo.Id != dataCall.ListenerId {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorPermissionDenied,
				})
				return
			}

			env := config.GetEnvValue()
			fee := env.Fee.ListenerCall
			if tokenInfo.Role == constant.RoleExperts {
				rate := env.Fee.ExpertsCall
				dataOrder, _ := paymentRepository.PublishInterfacePayment().GetOrderPaymentById(dataCall.OrderId)
				fee = int64(float64(dataOrder.DiamondPayment) * rate)
			}

			dataCall.StartCall = StartTime
			dataCall.EndCall = EndTime
			dataCall.Status = constant.Completed
			dataCall.Content = input.Content
			dataCall.ConsultingField = input.ConsultingField
			dataCall.ConsultingFeeListener = fee
			dataCall.WithdrawalDiamond = dataCall.ConsultingFeeListener / 2
			dataCall.UpdatedAt = time.Now().UTC()
			_ = repository.PublishInterfaceConversation().UpdateCallHistory(dataCall)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errCall.Error(),
			})
		}
	}
}

// @Summary
// @Schemes
// @Description API user submit call evaluation
// @Param input body model.InputSubmitCallEvaluation true "Input user submit call evaluation"
// @Tags Conversation
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /conversation/evaluation [post]
func SubmitCallEvaluation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token info
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
		input := model.InputSubmitCallEvaluation{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// get data call by id
		dataCall, errCall := repository.PublishInterfaceConversation().GetCallHistoryById(input.CallId)
		if errCall == nil {
			//check role is user
			if tokenInfo.Id != dataCall.UserId {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorPermissionDenied,
				})
				return
			}

			tipRate := config.GetEnvValue().Fee.ListenerTip
			listenerInfo, _ := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(dataCall.ListenerId)
			if listenerInfo.Role == constant.RoleExperts {
				tipRate = config.GetEnvValue().Fee.ExpertsTip
			}

			if input.Tip > 0 {
				userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(tokenInfo.Id)
				if userInfo.Diamond < input.Tip {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorDiamondNotEnough,
					})
					return
				}
				userInfo.Diamond = userInfo.Diamond - input.Tip
				userInfo.UpdatedAt = time.Now().UTC()
				_ = userRepository.PublishInterfaceUser().UpdateUser(userInfo)
			}

			// init model
			callEvaluation := model.CallEvaluation{}
			callEvaluation.Id = primitive.NewObjectID()
			callEvaluation.UserId = tokenInfo.Id
			callEvaluation.ListenerId = dataCall.ListenerId
			callEvaluation.CallId = input.CallId
			callEvaluation.NoteCompany = input.NoteCompany
			callEvaluation.NoteListener = input.NoteListener
			callEvaluation.Star = input.Star
			callEvaluation.Tip = input.Tip
			callEvaluation.TipRate = tipRate
			callEvaluation.CreatedAt = time.Now().UTC()
			callEvaluation.UpdatedAt = time.Now().UTC()
			_ = repository.PublishInterfaceConversation().CreateCallEvaluation(callEvaluation)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errCall.Error(),
			})
		}
	}
}

// @Summary
// @Schemes
// @Description Get list conversation history of user.
// @Tags Conversation
// @Accept json
// @Produce json
// @Header 200 {string} Token "jhfhhsid9834e8ff39fh"
// @Success 200 {object} model.ResponseGetConversationsUser
// @Failure 400 {object} model.ResponseError
// @Router  /conversation/user [get]
func GetConversationHistoryUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusBadRequest, masterModel.ResponseError{
				Code:    http.StatusBadRequest,
				Message: errToken.Error(),
			})
			return
		}
		if tokenInfo.Role != constant.RoleUser {
			c.JSON(http.StatusBadRequest, masterModel.ResponseError{
				Code:    http.StatusBadRequest,
				Message: message.MessageErrorPermissionDenied,
			})
			return
		}
		// get data conversations history of user
		data := repository.PublishInterfaceConversation().GetConversationsForUser(tokenInfo.Id)
		var responseData []model.DataConversationsUser
		for _, v := range data {
			item := model.DataConversationsUser{}
			item.Id = v.Id.Hex()
			listenerInfo, _ := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(v.ListenerId)
			item.ListenerName = common.GetFullNameOfListener(listenerInfo)
			item.EmployeeId = listenerInfo.EmployeeId
			item.ListenerImage = listenerInfo.Avatar
			item.ListenerRole = listenerInfo.Role
			orderInfo, _ := paymentRepository.PublishInterfacePayment().GetOrderPaymentById(v.OrderId)
			item.Date = orderInfo.Date.Format(constant.DateFormat)
			item.BookingTime = orderInfo.BookingTime
			item.TimeSlot = orderInfo.TimeSlot
			item.Status = v.Status
			responseData = append(responseData, item)
		}
		c.JSON(http.StatusOK, model.ResponseGetConversationsUser{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    responseData,
		})
		return

	}
}

// @Summary
// @Schemes
// @Description Get detail info conversation history of user by id.
// @Tags Conversation
// @Accept json
// @Produce json
// @Param  call_id path string true "Call Id."
// @Header 200 {string} Token "jhfhhsid9834e8ff39fh"
// @Success 200 {object} model.ResponseGetDetailConversationUser
// @Failure 400 {object} model.ResponseError
// @Router  /conversation/user/detail [get]
func GetDetailConversationHistoryUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusBadRequest, masterModel.ResponseError{
				Code:    http.StatusBadRequest,
				Message: errToken.Error(),
			})
			return
		}
		if tokenInfo.Role != constant.RoleUser {
			c.JSON(http.StatusBadRequest, masterModel.ResponseError{
				Code:    http.StatusBadRequest,
				Message: message.MessageErrorPermissionDenied,
			})
			return
		}

		callId := c.Query("call_id")
		// validate id
		if common.IsEmpty(callId) {
			c.JSON(http.StatusBadRequest, masterModel.ResponseError{
				Code:    http.StatusBadRequest,
				Message: message.MessageErrorCallIdEmpty,
			})
			return
		}

		callData, errData := repository.PublishInterfaceConversation().GetCallHistoryById(callId)
		if errData == nil {
			responseData := model.DataDetailConversationUser{}
			responseData.Id = callData.Id.Hex()
			listenerInfo, _ := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(callData.ListenerId)
			responseData.ListenerName = common.GetFullNameOfListener(listenerInfo)
			responseData.ListenerImage = listenerInfo.Avatar
			responseData.ListenerRole = listenerInfo.Role
			responseData.ConsultingField = callData.ConsultingField
			orderInfo, _ := paymentRepository.PublishInterfacePayment().GetOrderPaymentById(callData.OrderId)
			responseData.Date = orderInfo.Date.Format(constant.DateFormat)
			responseData.BookingTime = orderInfo.BookingTime
			responseData.TimeSlot = orderInfo.TimeSlot
			responseData.DiamondOrder = orderInfo.DiamondOrder
			responseData.DiamondDiscount = orderInfo.DiamondOrder - orderInfo.DiamondPayment
			responseData.DiamondPayment = orderInfo.DiamondPayment
			evaluationData, errEvaluate := repository.PublishInterfaceConversation().GetCallEvaluationByCallId(callId)
			if errEvaluate == nil {
				responseData.Star = evaluationData.Star
				responseData.Tip = evaluationData.Tip
				responseData.NoteListener = evaluationData.NoteListener
				responseData.NoteCompany = evaluationData.NoteCompany
			}
			c.JSON(http.StatusOK, model.ResponseGetDetailConversationUser{
				Code:    http.StatusOK,
				Message: message.MessageSuccess,
				Data:    responseData,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, masterModel.ResponseError{
				Code:    http.StatusBadRequest,
				Message: errData.Error(),
			})
			return
		}
	}
}

// @Summary
// @Schemes
// @Description Listener request extend call to user.
// @Param input body model.InputExtendCall true "Input send request extend call to user"
// @Tags Conversation
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /conversation/request-extend-call [post]
func RequestExtendCall() gin.HandlerFunc {
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
		input := model.InputExtendCall{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		callData, errCall := repository.PublishInterfaceConversation().GetCallHistoryById(input.CallId)
		if errCall != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errCall.Error(),
			})
			return
		} else {
			if callData.Status != constant.Talking || callData.ListenerId != tokenInfo.Id {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorPermissionDenied,
				})
				return
			}

			userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(callData.UserId)

			env := config.GetEnvValue()
			if userInfo.Diamond < env.Call.AppointmentPrice {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorDiamondNotEnough,
				})
				return
			}

			upcomingAppointment, errOrder := paymentRepository.PublishInterfacePayment().GetUpcomingAppointmentOfListener(callData.ListenerId)
			if errOrder == nil {
				timeNow, _ := time.Parse(constant.DateTimeFormat, time.Now().UTC().Add(constant.UTC7*time.Hour).Format(constant.DateTimeFormat))
				minute := math.Round(upcomingAppointment.CallDatetime.Sub(timeNow).Minutes())
				if minute < 45 {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorAppointmentComingUp,
					})
					return
				}
			}

			hub := wsControllers.GetHub()
			payload := wsControllers.SocketEventStruct{}
			payload.EventName = constant.WSEventExtendCall
			payload.EventPayload = map[string]interface{}{
				"call_id": callData.Id.Hex(),
				"status":  constant.Processing,
			}
			wsControllers.EmitToSpecificClient(hub, payload, callData.UserId)

			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
			return
		}
	}
}

// @Summary
// @Schemes
// @Description API Accept extend call of user
// @Param input body model.InputExtendCall true "Input extend call of user"
// @Tags Conversation
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /conversation/accept-extend-call [post]
func AcceptExtendCall() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token info
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
		input := model.InputExtendCall{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		callData, errCall := repository.PublishInterfaceConversation().GetCallHistoryById(input.CallId)
		if errCall != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errCall.Error(),
			})
			return
		} else {
			if callData.Status != constant.Talking || callData.UserId != tokenInfo.Id {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorPermissionDenied,
				})
				return
			}

			userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(callData.UserId)

			env := config.GetEnvValue()
			if userInfo.Diamond < env.Call.AppointmentPrice {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorDiamondNotEnough,
				})
				return
			}

			upcomingAppointment, errOrder := paymentRepository.PublishInterfacePayment().GetUpcomingAppointmentOfListener(callData.ListenerId)
			if errOrder == nil {
				timeNow, _ := time.Parse(constant.DateTimeFormat, time.Now().Format(constant.DateTimeFormat))
				minute := math.Round(upcomingAppointment.CallDatetime.Sub(timeNow).Minutes())
				if minute < 45 {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorUnableExtend,
					})
					return
				}
			}

			userInfo.Diamond -= env.Call.AppointmentPrice
			_ = userRepository.PublishInterfaceUser().UpdateUser(userInfo)

			hub := wsControllers.GetHub()
			payload := wsControllers.SocketEventStruct{}
			payload.EventName = constant.WSEventExtendCall
			payload.EventPayload = map[string]interface{}{
				"call_id": callData.Id.Hex(),
				"status":  constant.Completed,
			}
			wsControllers.EmitToSpecificClient(hub, payload, callData.ListenerId)

			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
			return
		}
	}
}
