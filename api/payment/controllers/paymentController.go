package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"sandexcare_backend/agora"
	conversationModel "sandexcare_backend/api/conversation/model"
	conversationRepository "sandexcare_backend/api/conversation/repository"
	listenerRepository "sandexcare_backend/api/listener/repository"
	masterRepository "sandexcare_backend/api/master/repository"
	"sandexcare_backend/api/payment/model"
	"sandexcare_backend/api/payment/repository"
	scheduleRepository "sandexcare_backend/api/schedule/repository"
	userRepository "sandexcare_backend/api/user/repository"
	"sandexcare_backend/helpers/cache"
	"sandexcare_backend/helpers/common"
	"sandexcare_backend/helpers/config"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/message"
	"sandexcare_backend/helpers/middlewares"
	wsControllers "sandexcare_backend/server_websocket/controllers"
	"strings"
	"time"
)

// @Summary
// @Schemes
// @Description API payment order appointment
// @Param input body model.InputOrderPayment true "Input payment order appointment"
// @Tags Payment
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /payment [post]
func OrderPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorTokenInvalid,
			})
			return
		}
		if tokenInfo.Role != constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputOrderPayment{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate listener_id
		if common.IsEmpty(input.ListenerId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorListenerIdEmpty,
			})
			return
		}

		listenerInfo, errListener := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(input.ListenerId)
		if errListener != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errListener.Error(),
			})
			return
		}

		// validate time_slot
		if common.IsEmpty(input.TimeSlot) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorTimeSlotEmpty,
			})
			return
		}

		// validate booking_time
		if common.IsEmpty(input.BookingTime) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorBookingTimeEmpty,
			})
			return
		}

		// validate date
		if common.IsEmpty(input.Date) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateEmpty,
			})
			return
		}

		// convert date from string to date format
		date, err := time.Parse(constant.DateFormat, input.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateFormat,
			})
			return
		}

		callTimeString := input.Date + " " + input.BookingTime + ":00"
		// convert date from string to date format
		callDatetime, _ := time.Parse(constant.DateTimeFormat, callTimeString)

		timeNow, _ := time.Parse(constant.DateTimeFormat, time.Now().UTC().Add(constant.UTC7*time.Hour).Format(constant.DateTimeFormat))
		if callDatetime.Sub(timeNow).Minutes() < 20 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorAppointmentLate,
			})
			return
		}

		dateNow, _ := time.Parse(constant.DateFormat, time.Now().Format(constant.DateFormat))
		if date.Sub(dateNow).Hours()/24 > 6 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorOrderDateLimit,
			})
			return
		}

		if !scheduleRepository.PublishInterfaceSchedule().CheckScheduleWorkExistListener(input.ListenerId, date, input.TimeSlot) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorScheduleNotExist,
			})
			return
		}

		dataBooked := repository.PublishInterfacePayment().GetScheduleAppointmentForUser(tokenInfo.Id)
		for _, b := range dataBooked {
			listenerBooked, _ := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(b.ListenerId)
			if listenerBooked.Role == listenerInfo.Role {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorOrderNotAllowed,
				})
				return
			}
		}

		// check booking exist
		if repository.PublishInterfacePayment().CheckOrderPaymentExist(input.ListenerId, callDatetime) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorAppointmentExist,
			})
			return
		} else {
			timeSlotData, errTimeSlot := masterRepository.PublishInterfaceMaster().GetDetailTimeSlot(input.TimeSlot)
			if errTimeSlot == nil {
				env := config.GetEnvValue()
				callPrice := env.Call.AppointmentPrice
				surcharge := int64(float64(callPrice) * timeSlotData.Surcharge)
				diamondOrder := callPrice + surcharge
				var discount int64
				if !common.IsEmpty(input.CouponId) {
					couponInfo, errCoupon := userRepository.PublishInterfaceUser().GetCouponInfo(tokenInfo.Id, input.CouponId)
					if errCoupon != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"code":    http.StatusBadRequest,
							"message": message.MessageErrorCouponNotExist,
						})
						return
					}
					if couponInfo.Type != constant.CouponBookingCV {
						c.JSON(http.StatusBadRequest, gin.H{
							"code":    http.StatusBadRequest,
							"message": message.MessageErrorCouponInvalid,
						})
						return
					}
					discount = int64(float64(diamondOrder) * couponInfo.Discount)
				}

				diamondPayment := diamondOrder - discount
				userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(tokenInfo.Id)
				if userInfo.Diamond < diamondPayment {
					gc := cache.Cache()
					key := cache.CreateKeyLockPayment(input.ListenerId, input.Date, input.BookingTime)
					_, exist := gc.Get(key)
					if !exist {
						gc.Set(key, tokenInfo.Id, 6*time.Minute)
					}
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorDiamondNotEnough,
					})
					return
				}
				// init model
				var order model.OrderPayment
				order.Id = primitive.NewObjectID()
				order.UserId = tokenInfo.Id
				order.ListenerId = input.ListenerId
				order.DiamondOrder = diamondOrder
				order.DiamondPayment = diamondPayment
				order.Surcharge = surcharge
				order.CouponId = input.CouponId
				order.Date = date
				order.TimeSlot = input.TimeSlot
				order.BookingTime = input.BookingTime
				order.CallDatetime = callDatetime
				order.Type = constant.OrderBookAppointment
				order.Status = constant.Active
				order.CreatedAt = time.Now().UTC()
				order.UpdatedAt = time.Now().UTC()

				// handle create data order call now
				errCreateOrder := repository.PublishInterfacePayment().CreateOrderPayment(order)
				if errCreateOrder == nil {
					userInfo.Diamond -= diamondPayment
					userInfo.UpdatedAt = time.Now().UTC()
					_ = userRepository.PublishInterfaceUser().UpdateUser(userInfo)
					if !common.IsEmpty(order.CouponId) {
						couponInfo, _ := userRepository.PublishInterfaceUser().GetCouponInfo(tokenInfo.Id, input.CouponId)
						couponInfo.Status = constant.Used
						couponInfo.UpdatedAt = time.Now().UTC()
						_ = userRepository.PublishInterfaceUser().UpdateCouponsUser(couponInfo)
					}
					c.JSON(http.StatusOK, gin.H{
						"code":    http.StatusOK,
						"message": message.MessageSuccess,
					})
					return
				} else {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": errCreateOrder.Error(),
					})
					return
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": errTimeSlot.Error(),
				})
				return
			}
		}
	}
}

// @Summary
// @Schemes
// @Description API thanh toán đặt lịch chuyên gia
// @Param input body model.InputOrderPayment true "Input payment order appointment"
// @Tags Payment
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /payment/booking-expert [post]
func PaymentBookingExpert() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorTokenInvalid,
			})
			return
		}
		if tokenInfo.Role != constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputOrderPayment{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate listener_id
		if common.IsEmpty(input.ListenerId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorListenerIdEmpty,
			})
			return
		}

		listenerInfo, errListener := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(input.ListenerId)
		if errListener != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errListener.Error(),
			})
			return
		}

		// validate time_slot
		if common.IsEmpty(input.TimeSlot) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorTimeSlotEmpty,
			})
			return
		}

		// validate booking_time
		if common.IsEmpty(input.BookingTime) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorBookingTimeEmpty,
			})
			return
		}

		// validate date
		if common.IsEmpty(input.Date) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateEmpty,
			})
			return
		}

		// convert date from string to date format
		date, err := time.Parse(constant.DateFormat, input.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateFormat,
			})
			return
		}

		callTimeString := input.Date + " " + input.BookingTime + ":00"
		// convert date from string to date format
		callDatetime, _ := time.Parse(constant.DateTimeFormat, callTimeString)

		timeNow, _ := time.Parse(constant.DateTimeFormat, time.Now().UTC().Add(constant.UTC7*time.Hour).Format(constant.DateTimeFormat))
		if callDatetime.Sub(timeNow).Minutes() < 20 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorAppointmentLate,
			})
			return
		}

		dateNow, _ := time.Parse(constant.DateFormat, time.Now().Format(constant.DateFormat))
		if date.Sub(dateNow).Hours()/24 > 6 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorOrderDateLimit,
			})
			return
		}

		if !scheduleRepository.PublishInterfaceSchedule().CheckScheduleWorkExistListener(input.ListenerId, date, input.TimeSlot) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorScheduleNotExist,
			})
			return
		}

		dataBooked := repository.PublishInterfacePayment().GetScheduleAppointmentForUser(tokenInfo.Id)
		for _, b := range dataBooked {
			listenerBooked, _ := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(b.ListenerId)
			if listenerBooked.Role == listenerInfo.Role {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorOrderNotAllowed,
				})
				return
			}
		}

		// check booking exist
		if repository.PublishInterfacePayment().CheckOrderPaymentExist(input.ListenerId, callDatetime) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorAppointmentExist,
			})
			return
		} else {
			diamondOrder := listenerInfo.Price
			var discount int64
			if !common.IsEmpty(input.CouponId) {
				couponInfo, errCoupon := userRepository.PublishInterfaceUser().GetCouponInfo(tokenInfo.Id, input.CouponId)
				if errCoupon != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorCouponNotExist,
					})
					return
				}
				if couponInfo.Type != constant.CouponNameBookingCG {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorCouponInvalid,
					})
					return
				}
				discount = int64(float64(diamondOrder) * couponInfo.Discount)
			}

			diamondPayment := diamondOrder - discount
			userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(tokenInfo.Id)
			if userInfo.Diamond < diamondPayment {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorDiamondNotEnough,
				})
				return
			}
			// init model
			var order model.OrderPayment
			order.Id = primitive.NewObjectID()
			order.UserId = tokenInfo.Id
			order.ListenerId = input.ListenerId
			order.DiamondOrder = diamondOrder
			order.DiamondPayment = diamondPayment
			order.CouponId = input.CouponId
			order.Date = date
			order.TimeSlot = input.TimeSlot
			order.BookingTime = input.BookingTime
			order.CallDatetime = callDatetime
			order.Type = constant.OrderBookAppointment
			order.Status = constant.Active
			order.CreatedAt = time.Now().UTC()
			order.UpdatedAt = time.Now().UTC()

			// handle create data order call now
			errCreateOrder := repository.PublishInterfacePayment().CreateOrderPayment(order)
			if errCreateOrder == nil {
				userInfo.Diamond -= diamondPayment
				userInfo.UpdatedAt = time.Now().UTC()
				_ = userRepository.PublishInterfaceUser().UpdateUser(userInfo)
				if !common.IsEmpty(order.CouponId) {
					couponInfo, _ := userRepository.PublishInterfaceUser().GetCouponInfo(tokenInfo.Id, input.CouponId)
					couponInfo.Status = constant.Used
					couponInfo.UpdatedAt = time.Now().UTC()
					_ = userRepository.PublishInterfaceUser().UpdateCouponsUser(couponInfo)
				}

				// gửi sms thông báo đến chuyên gia
				content := message.SmsBooking
				dayOfWeek := common.GetDayOfWeek(date)
				dateSMS := date.Format("02/01/2006")
				replacer := strings.NewReplacer("{hour}", input.BookingTime, "{day_of_week}", dayOfWeek, "{date}", dateSMS)
				content = replacer.Replace(content)
				common.SendSMS(listenerInfo.PhoneNumber, content)

				dataScheduleWork, _ := scheduleRepository.PublishInterfaceSchedule().GetDetailScheduleWorkListener(input.ListenerId, date, input.TimeSlot)
				dataScheduleWork.OrderStatus = "yes"
				dataScheduleWork.UpdatedAt = time.Now().UTC()
				_ = scheduleRepository.PublishInterfaceSchedule().UpdateScheduleWork(dataScheduleWork)

				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": message.MessageSuccess,
				})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": errCreateOrder.Error(),
				})
				return
			}
		}
	}
}

// @Summary
// @Schemes
// @Description API call payment
// @Param input body model.InputCallPayment true "Input call payment"
// @Tags Payment
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /payment/call [post]
func CallPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorTokenInvalid,
			})
			return
		}
		if tokenInfo.Role != constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputCallPayment{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate time_slot
		if common.IsEmpty(input.TimeSlot) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorTimeSlotEmpty,
			})
			return
		}

		// validate booking_time
		if common.IsEmpty(input.BookingTime) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorBookingTimeEmpty,
			})
			return
		}

		// validate date
		if common.IsEmpty(input.Date) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateEmpty,
			})
			return
		}

		// convert date from string to date format
		date, err := time.Parse(constant.DateFormat, input.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateFormat,
			})
			return
		}

		callTimeString := input.Date + " " + input.BookingTime + ":00"
		// convert date from string to date format
		callDatetime, _ := time.Parse(constant.DateTimeFormat, callTimeString)

		timeSlotData, errTimeSlot := masterRepository.PublishInterfaceMaster().GetDetailTimeSlot(input.TimeSlot)
		if errTimeSlot == nil {
			env := config.GetEnvValue()
			callPrice := env.Call.AppointmentPrice
			surcharge := int64(float64(callPrice) * timeSlotData.Surcharge)
			diamondOrder := callPrice + surcharge
			var discount int64
			if !common.IsEmpty(input.CouponId) {
				couponInfo, errCoupon := userRepository.PublishInterfaceUser().GetCouponInfo(tokenInfo.Id, input.CouponId)
				if errCoupon != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorCouponNotExist,
					})
					return
				}
				if couponInfo.Type != constant.CouponCallNow {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorCouponInvalid,
					})
					return
				}
				discount = int64(float64(diamondOrder) * couponInfo.Discount)
			}

			diamondPayment := diamondOrder - discount
			userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(tokenInfo.Id)
			if userInfo.Diamond < diamondPayment {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorDiamondNotEnough,
				})
				return
			}
			// init model
			var order model.OrderPayment
			order.Id = primitive.NewObjectID()
			order.UserId = tokenInfo.Id
			order.DiamondOrder = diamondOrder
			order.DiamondPayment = diamondPayment
			order.Surcharge = surcharge
			order.Date = date
			order.TimeSlot = input.TimeSlot
			order.BookingTime = input.BookingTime
			order.CallDatetime = callDatetime
			order.CouponId = input.CouponId
			order.Type = constant.OrderBookCallNow
			order.Status = constant.Active
			order.CreatedAt = time.Now().UTC()
			order.UpdatedAt = time.Now().UTC()

			hub := wsControllers.GetHub()
			listenerOnCall := wsControllers.GetAllListenersReadyCall(hub)
			var listenerId string
			for _, v := range listenerOnCall {
				if v.ActiveCallNow {
					listenerId = v.ListenerId
					break
				}
			}
			if common.IsEmpty(listenerId) {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorListenerUnconnect,
				})
				return
			}

			// handle create data order call now
			errCreateOrder := repository.PublishInterfacePayment().CreateOrderPayment(order)
			if errCreateOrder == nil {
				userInfo.Diamond -= diamondPayment
				userInfo.UpdatedAt = time.Now().UTC()
				_ = userRepository.PublishInterfaceUser().UpdateUser(userInfo)
				if !common.IsEmpty(order.CouponId) {
					couponInfo, _ := userRepository.PublishInterfaceUser().GetCouponInfo(tokenInfo.Id, input.CouponId)
					couponInfo.Status = constant.Used
					couponInfo.UpdatedAt = time.Now().UTC()
					_ = userRepository.PublishInterfaceUser().UpdateCouponsUser(couponInfo)
				}

				// init model conversation ==> start call
				var callHistory conversationModel.CallHistory
				callHistory.Id = primitive.NewObjectID()
				callHistory.UserId = tokenInfo.Id
				callHistory.ListenerId = listenerId
				callHistory.Channel = common.GenerateNumber(6)
				callHistory.OrderId = order.Id.Hex()
				callHistory.Status = constant.Started
				callHistory.CreatedAt = time.Now().UTC()
				callHistory.UpdatedAt = time.Now().UTC()

				errCallHistory := conversationRepository.PublishInterfaceConversation().CreateCallHistory(callHistory)
				if errCallHistory == nil {
					uidListener := common.GenerateNumber(6)
					uidUser := common.GenerateNumber(6)

					userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(tokenInfo.Id)
					payload := wsControllers.SocketEventStruct{}
					payload.EventName = constant.WSEventStartCall
					payload.EventPayload = map[string]interface{}{
						"order_id":     order.Id.Hex(),
						"call_id":      callHistory.Id.Hex(),
						"username":     userInfo.Name,
						"date":         order.Date.Format(constant.DateFormat),
						"booking_time": order.BookingTime,
						"token":        agora.GenerateRtcToken(uidListener, callHistory.Channel),
						"channel":      callHistory.Channel,
						"uid":          uidListener,
						"ringing_time": env.Call.RingingTimeCallNow,
					}
					wsControllers.EmitToSpecificClient(hub, payload, callHistory.ListenerId)

					dataResponse := model.DataPaymentCallNow{}
					dataResponse.OrderId = order.Id.Hex()
					dataResponse.CallId = callHistory.Id.Hex()
					dataResponse.Uid = uidUser
					dataResponse.Channel = callHistory.Channel
					dataResponse.Token = agora.GenerateRtcToken(uidUser, callHistory.Channel)
					dataResponse.RingingTime = env.Call.RingingTimeCallNow

					response := model.ResponsePaymentCallNow{
						Code:    http.StatusOK,
						Message: message.MessageSuccess,
						Data:    dataResponse,
					}

					c.JSON(http.StatusOK, response)
					return
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": errCreateOrder.Error(),
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errTimeSlot.Error(),
			})
			return
		}
	}
}

// @Summary
// @Schemes
// @Description API update info appointment
// @Param input body model.InputUpdateOrderPayment true "Input update appointment"
// @Tags Payment
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /payment/update [post]
func UpdateOrderPayment() gin.HandlerFunc {
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
		input := model.InputUpdateOrderPayment{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate order_id
		if common.IsEmpty(input.OrderId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorOrderIdEmpty,
			})
			return
		}

		// validate listener_id
		if common.IsEmpty(input.ListenerId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorListenerIdEmpty,
			})
			return
		}

		// validate time_slot
		if common.IsEmpty(input.TimeSlot) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorTimeSlotEmpty,
			})
			return
		}

		// validate booking_time
		if common.IsEmpty(input.BookingTime) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorBookingTimeEmpty,
			})
			return
		}

		// validate date
		if common.IsEmpty(input.Date) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateEmpty,
			})
			return
		}

		// convert date from string to date format
		date, err := time.Parse(constant.DateFormat, input.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateFormat,
			})
			return
		}

		callTimeString := input.Date + " " + input.BookingTime
		// convert date from string to date format
		callDatetime, _ := time.Parse(constant.DateTimeFormat, callTimeString)

		timeNow, _ := time.Parse(constant.DateTimeFormat, time.Now().UTC().Add(constant.UTC7*time.Hour).Format(constant.DateTimeFormat))
		if callDatetime.Sub(timeNow).Minutes() < 20 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorAppointmentLate,
			})
			return
		}

		// check booking exist
		if repository.PublishInterfacePayment().CheckOrderPaymentExist(input.ListenerId, callDatetime) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorAppointmentExist,
			})
			return
		} else {
			userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(tokenInfo.Id)
			if userInfo.Diamond < input.Surcharge {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorDiamondNotEnough,
				})
				return
			}

			orderInfo, errOrder := repository.PublishInterfacePayment().GetOrderPaymentById(input.OrderId)
			if errOrder != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": errOrder.Error(),
				})
				return
			} else {
				if tokenInfo.Id != orderInfo.UserId || orderInfo.Status != constant.Active {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorPermissionDenied,
					})
					return
				}
				orderInfo.Date = date
				orderInfo.TimeSlot = input.TimeSlot
				orderInfo.BookingTime = input.BookingTime
				orderInfo.ListenerId = input.ListenerId
				orderInfo.CallDatetime = callDatetime
				orderInfo.DiamondPayment += input.Surcharge
				orderInfo.Surcharge = input.Surcharge
				orderInfo.UpdateFlag = true
				orderInfo.UpdatedAt = time.Now().UTC()

				_ = repository.PublishInterfacePayment().UpdateOrderPayment(orderInfo)

				userInfo.Diamond -= input.Surcharge
				_ = userRepository.PublishInterfaceUser().UpdateUser(userInfo)

				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": message.MessageSuccess,
				})
				return
			}
		}
	}
}

// @Summary
// @Schemes
// @Description API hoàn tiền chuyên viên. type = 2 chuyên viên không nghe máy. Tất cả trường hợp khác type = 1.
// @Param input body model.InputPaymentRefund true "Input request payment refund"
// @Tags Payment
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /payment/refund [post]
func PaymentRefund() gin.HandlerFunc {
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
		input := model.InputPaymentRefund{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate order_id
		if common.IsEmpty(input.OrderId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorOrderIdEmpty,
			})
			return
		}

		userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(tokenInfo.Id)

		orderInfo, errOrder := repository.PublishInterfacePayment().GetOrderPaymentById(input.OrderId)
		if errOrder != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errOrder.Error(),
			})
			return
		} else {
			if tokenInfo.Id != orderInfo.UserId || orderInfo.Status != constant.Active {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorPermissionDenied,
				})
				return
			}

			if input.Type == 2 {
				userInfo.Diamond += int64(float64(orderInfo.DiamondPayment) * 1.5)
				if !common.IsEmpty(orderInfo.CouponId) {
					couponInfo, _ := userRepository.PublishInterfaceUser().GetCouponInfo(tokenInfo.Id, orderInfo.CouponId)
					couponInfo.Status = constant.Active
					couponInfo.UpdatedAt = time.Now().UTC()
					_ = userRepository.PublishInterfaceUser().UpdateCouponsUser(couponInfo)
				}
			} else {
				userInfo.Diamond += orderInfo.DiamondPayment
			}
			orderInfo.Status = constant.Refunded
			orderInfo.UpdatedAt = time.Now().UTC()

			_ = repository.PublishInterfacePayment().UpdateOrderPayment(orderInfo)
			_ = userRepository.PublishInterfaceUser().UpdateUser(userInfo)

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
// @Description API hoàn tiền chuyên gia. type = 2 chuyên gia không nghe máy. Tất cả trường hợp khác type = 1.
// @Param input body model.InputPaymentRefund true "Input request payment refund for experts"
// @Tags Payment
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /payment/refund-expert [post]
func PaymentRefundExpert() gin.HandlerFunc {
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
		input := model.InputPaymentRefund{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate order_id
		if common.IsEmpty(input.OrderId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorOrderIdEmpty,
			})
			return
		}

		userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(tokenInfo.Id)

		orderInfo, errOrder := repository.PublishInterfacePayment().GetOrderPaymentById(input.OrderId)
		if errOrder != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errOrder.Error(),
			})
			return
		} else {
			if tokenInfo.Id != orderInfo.UserId || orderInfo.Status != constant.Active {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorPermissionDenied,
				})
				return
			}

			if input.Type == 2 {
				userInfo.Diamond += int64(float64(orderInfo.DiamondPayment) * 1.1)
				if !common.IsEmpty(orderInfo.CouponId) {
					couponInfo, _ := userRepository.PublishInterfaceUser().GetCouponInfo(tokenInfo.Id, orderInfo.CouponId)
					couponInfo.Status = constant.Active
					couponInfo.UpdatedAt = time.Now().UTC()
					_ = userRepository.PublishInterfaceUser().UpdateCouponsUser(couponInfo)
				}
			} else {
				userInfo.Diamond += orderInfo.DiamondPayment
			}

			orderInfo.Status = constant.Refunded
			orderInfo.UpdatedAt = time.Now().UTC()

			_ = repository.PublishInterfacePayment().UpdateOrderPayment(orderInfo)
			_ = userRepository.PublishInterfaceUser().UpdateUser(userInfo)

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
// @Description API user unlock listener when payment fail
// @Param input body model.InputUnlock true "Input unlock listener"
// @Tags Payment
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Failure 401 {object} model.ResponseError
// @Router  /payment/unlock [post]
func Unlock() gin.HandlerFunc {
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
		input := model.InputUnlock{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate listenerId
		if common.IsEmpty(input.ListenerId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorListenerIdEmpty,
			})
			return
		}
		// validate date
		if common.IsEmpty(input.Date) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateEmpty,
			})
			return
		}
		// validate booking time
		if common.IsEmpty(input.BookingTime) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorBookingTimeEmpty,
			})
			return
		}
		gc := cache.Cache()
		key := cache.CreateKeyLockPayment(input.ListenerId, input.Date, input.BookingTime)
		dataCache, exist := gc.Get(key)
		if exist && dataCache == tokenInfo.Id {
			gc.Delete(key)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
		})
		return
	}
}
