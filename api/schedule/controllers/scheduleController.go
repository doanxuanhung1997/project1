package controllers

import (
	"encoding/json"
	"net/http"
	listenerRepository "sandexcare_backend/api/listener/repository"
	masterRepository "sandexcare_backend/api/master/repository"
	paymentModel "sandexcare_backend/api/payment/model"
	paymentRepository "sandexcare_backend/api/payment/repository"
	"sandexcare_backend/api/schedule/model"
	"sandexcare_backend/api/schedule/repository"
	"sandexcare_backend/helpers/cache"
	"sandexcare_backend/helpers/common"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/message"
	"sandexcare_backend/helpers/middlewares"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary
// @Schemes
// @Description API Create schedule work for listener
// @Param input body model.InputCreateScheduleWork true "Input create schedule for listener"
// @Tags Schedule
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /schedule-work/create [post]
func CreateSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get token info
		_, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorTokenInvalid,
			})
			return
		}

		////check role is Admin
		//if tokenInfo.Role != constant.RoleAdmin {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"code":    http.StatusBadRequest,
		//		"message": message.MessageErrorPermissionDenied,
		//	})
		//	return
		//}

		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputCreateScheduleWork{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate employee_id
		if common.IsEmpty(input.EmployeeId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorEmployeeIdEmpty,
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

		// validate date work
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

		dateNow, _ := time.Parse(constant.DateFormat, time.Now().Format(constant.DateFormat))
		if !date.After(dateNow) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateScheduleInvalid,
			})
			return
		}

		listenerInfo, errInfo := listenerRepository.PublishInterfaceListener().GetListenerByEmployeeId(input.EmployeeId)
		if errInfo != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorEmployeeIdNotExist,
			})
			return
		}

		// init model
		var schedule model.ScheduleWork
		schedule.Id = primitive.NewObjectID()
		schedule.ListenerId = listenerInfo.Id.Hex()
		schedule.Date = date
		schedule.TimeSlot = input.TimeSlot
		schedule.Status = constant.Active
		schedule.CreatedAt = time.Now().UTC()
		schedule.UpdatedAt = time.Now().UTC()

		//TODO check trùng lịch

		// handle create data schedule
		errCreateSchedule := repository.PublishInterfaceSchedule().CreateScheduleWork(schedule)
		if errCreateSchedule == nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errCreateSchedule.Error(),
			})
			return
		}
	}
}

// @Summary
// @Schemes
// @Description API Get schedule work for listener.
// @Param  start_date path string true "Start date" Format(string)
// @Param  end_date path string true "End date" Format(string)
// @Tags Schedule
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseScheduleWork
// @Failure 401 {object} model.ResponseError
// @Failure 400 {object} model.ResponseError
// @Router  /schedule-work [get]
func GetScheduleWorkForListener() gin.HandlerFunc {
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
		//check role is listener
		if tokenInfo.Role == constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		//get listener info by phone_number in token
		listener, _ := listenerRepository.PublishInterfaceListener().GetListenerByPhoneNumber(tokenInfo.PhoneNumber)
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		// validate start_date
		if common.IsEmpty(startDate) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateEmpty,
			})
			return
		}

		// convert start_date from string to date format
		dateFrom, err := time.Parse(constant.DateFormat, startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateFormat,
			})
			return
		}

		// validate end_date
		if common.IsEmpty(endDate) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateEmpty,
			})
			return
		}

		// convert end_date from string to date format
		dateTo, err := time.Parse(constant.DateFormat, endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateFormat,
			})
			return
		}

		dataSchedule := repository.PublishInterfaceSchedule().GetScheduleWordForListener(listener.Id.Hex(), dateFrom, dateTo)
		c.JSON(http.StatusOK, model.ResponseScheduleWork{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    dataSchedule,
		})
	}
}

// @Summary
// @Schemes
// @Description API Get schedule work in day for listener.
// @Param  day path string true "Day" Format(string)
// @Param  time_slot path string true "Time Slot - epx: 00:00-06:00 " Format(string)
// @Tags Schedule
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseWorkingDay
// @Failure 401 {object} model.ResponseError
// @Failure 400 {object} model.ResponseError
// @Router  /schedule-work/day [get]
func GetScheduleInDay() gin.HandlerFunc {
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
		//check role is listener
		if tokenInfo.Role == constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		//get listener info by phone_number in token
		listener, _ := listenerRepository.PublishInterfaceListener().GetListenerByPhoneNumber(tokenInfo.PhoneNumber)

		day := c.Query("day")
		timeSlot := c.Query("time_slot")

		// validate time_slot
		if common.IsEmpty(timeSlot) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorTimeSlotEmpty,
			})
			return
		}

		// validate day
		if common.IsEmpty(day) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateEmpty,
			})
			return
		}

		// convert day from string to day format
		workingDay, err := time.Parse(constant.DateFormat, day)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateFormat,
			})
			return
		}

		data := repository.PublishInterfaceSchedule().GetScheduleInDay(listener.Id.Hex(), workingDay, timeSlot)
		c.JSON(http.StatusOK, model.ResponseWorkingDay{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    data,
		})
	}
}

// @Summary
// @Schemes
// @Description API Get schedule work in day for listener.
// @Param  day path string true "Day" Format(string)
// @Param  listener_id path string true "Listener Id" Format(string)
// @Tags Schedule
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetDetailScheduleWorkListener
// @Failure 401 {object} model.ResponseError
// @Failure 400 {object} model.ResponseError
// @Router  /schedule-work/detail-date/listener [get]
func GetDetailScheduleWorkListener() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get token info
		_, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorTokenInvalid,
			})
			return
		}

		day := c.Query("day")
		listenerId := c.Query("listener_id")

		// validate date
		if common.IsEmpty(day) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateEmpty,
			})
			return
		}
		// convert date from string to date format
		date, err := time.Parse(constant.DateFormat, day)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateFormat,
			})
			return
		}

		// validate listener_id
		if common.IsEmpty(listenerId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorListenerIdEmpty,
			})
			return
		}

		_, errListener := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(listenerId)
		if errListener != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errListener.Error(),
			})
			return
		}
		timeSlotData := masterRepository.PublishInterfaceMaster().GetAllTimeSlot()
		var bookingAble []string
		for _, v := range timeSlotData {
			if repository.PublishInterfaceSchedule().CheckScheduleWorkExistListener(listenerId, date, v.TimeSlot) {
				for _, b := range v.BookingTime {
					callTimeString := day + " " + b
					// convert date from string to date format
					callDatetime, _ := time.Parse(constant.DateTimeFormat, callTimeString)
					if !paymentRepository.PublishInterfacePayment().CheckOrderPaymentExist(listenerId, callDatetime) {
						bookingAble = append(bookingAble, b)
					}
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
			"data":    bookingAble,
		})
		return
	}
}

// @Summary
// @Schemes
// @Description API Get schedule work in day for listener.
// @Param  day path string true "Day" Format(string)
// @Param  time_slot path string true "Time slot" Format(string)
// @Tags Schedule
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetScheduleWorkAppointment
// @Failure 401 {object} model.ResponseError
// @Failure 400 {object} model.ResponseError
// @Router  /schedule-work/appointment [get]
func GetScheduleWorkAppointment() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get token info
		_, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": message.MessageErrorTokenInvalid,
			})
			return
		}

		timeSlot := c.Query("time_slot")
		day := c.Query("day")

		// validate time_slot
		if common.IsEmpty(timeSlot) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorTimeSlotEmpty,
			})
			return
		}
		// validate day
		if common.IsEmpty(day) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateEmpty,
			})
			return
		}
		// convert date from string to date format
		date, err := time.Parse(constant.DateFormat, day)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateFormat,
			})
			return
		}
		var responseData []model.DataDetailAppointmentInDay
		timeSlotData, errTimeSlot := masterRepository.PublishInterfaceMaster().GetDetailTimeSlot(timeSlot)
		if errTimeSlot == nil {
			listenersInfo := repository.PublishInterfaceSchedule().GetListenersWorkTimeSlot(timeSlot, date)
			if len(listenersInfo) > 0 {
				for _, b := range timeSlotData.BookingTime {
					responseItem := model.DataDetailAppointmentInDay{}
					responseItem.BookingTime = b
					orderData := paymentRepository.PublishInterfacePayment().GetAppointmentsBooked(date, b)
					if len(orderData) == len(listenersInfo) {
						responseItem.IsFull = true
					}
					responseData = append(responseData, responseItem)
				}
			}
		}

		c.JSON(http.StatusOK, model.ResponseGetScheduleWorkAppointment{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    responseData,
		})
		return
	}
}

func Contains(s []paymentModel.OrderPayment, str string) bool {
	for _, v := range s {
		if v.ListenerId == str {
			return true
		}
	}
	return false
}

// @Summary
// @Schemes
// @Description API Get listener to book appointment.
// @Param  day path string true "Day" Format(string)
// @Param  time_slot path string true "Time Slot" Format(string)
// @Param  booking_time path string true "Booking time of user" Format(string)
// @Tags Schedule
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetListenerForBookAppointment
// @Failure 401 {object} model.ResponseError
// @Failure 400 {object} model.ResponseError
// @Router  /schedule-work/appointment/listener [get]
func GetListenerForBookAppointment() gin.HandlerFunc {
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

		day := c.Query("day")
		timeSlot := c.Query("time_slot")
		bookingTime := c.Query("booking_time")

		// validate date
		if common.IsEmpty(day) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateEmpty,
			})
			return
		}
		// convert date from string to date format
		date, err := time.Parse(constant.DateFormat, day)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDateFormat,
			})
			return
		}

		// validate time_slot
		if common.IsEmpty(timeSlot) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorTimeSlotEmpty,
			})
			return
		}

		// validate booking_time
		if common.IsEmpty(bookingTime) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorBookingTimeEmpty,
			})
			return
		}

		var responseData []model.DataListenerInfo
		listenersInfo := repository.PublishInterfaceSchedule().GetListenersWorkTimeSlot(timeSlot, date)
		if len(listenersInfo) > 0 {
			orderData := paymentRepository.PublishInterfacePayment().GetAppointmentsBooked(date, bookingTime)
			for i, _ := range listenersInfo {
				if !Contains(orderData, listenersInfo[i].ListenerId) {
					gc := cache.Cache()
					key := cache.CreateKeyLockPayment(listenersInfo[i].ListenerId, day, bookingTime)
					dataLock, flag := gc.Get(key)
					if flag && dataLock != tokenInfo.Id {
						continue
					}
					responseData = append(responseData, listenersInfo[i])
				}
			}
		}

		c.JSON(http.StatusOK, model.ResponseGetListenerForBookAppointment{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    responseData,
		})
		return
	}
}

//// @Summary
//// @Schemes
//// @Description API get chi tiết ca làm việc của chuyên gia để đặt lịch hẹn
//// @Param  day path string true "Day" Format(string)
//// @Param  listener_id path string true "Chuyên viên Id" Format(string)
//// @Tags Schedule
//// @Accept json
//// @Produce json
//// @Success 200 {object} model.ResponseGetScheduleWorkAppointment
//// @Failure 401 {object} model.ResponseError
//// @Failure 400 {object} model.ResponseError
//// @Router  /schedule-work/appointment/expert [get]
//func GetDetailScheduleWorkExpertForBookAppointment() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		//get token info
//		tokenInfo, errToken := middlewares.AuthenticateToken(c)
//		if errToken != nil {
//			c.JSON(http.StatusUnauthorized, gin.H{
//				"code":    http.StatusUnauthorized,
//				"message": message.MessageErrorTokenInvalid,
//			})
//			return
//		}
//
//		day := c.Query("day")
//		listenerId := c.Query("listener_id")
//
//		// validate date
//		if common.IsEmpty(day) {
//			c.JSON(http.StatusBadRequest, gin.H{
//				"code":    http.StatusBadRequest,
//				"message": message.MessageErrorDateEmpty,
//			})
//			return
//		}
//		// convert date from string to date format
//		date, err := time.Parse(constant.DateFormat, day)
//		if err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{
//				"code":    http.StatusBadRequest,
//				"message": message.MessageErrorDateFormat,
//			})
//			return
//		}
//
//		// validate listener_id
//		if common.IsEmpty(listenerId) {
//			c.JSON(http.StatusBadRequest, gin.H{
//				"code":    http.StatusBadRequest,
//				"message": message.MessageErrorExpertIdEmpty,
//			})
//			return
//		}
//
//		bookingTimeExpert := [...]string{"09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"}
//		var responseData []model.DataDetailAppointmentInDay
//		for _, b := range bookingTimeExpert {
//			responseItem := model.DataDetailAppointmentInDay{}
//			responseItem.BookingTime = b + ":00"
//			temp, _ := strconv.Atoi(b)
//			timeSlot := b + ":00" + "-" + strconv.Itoa(temp+1) + ":00"
//
//			if repository.PublishInterfaceSchedule().CheckScheduleWorkExistListener(listenerId, date, timeSlot) {
//				callTimeString := day + " " + responseItem.BookingTime + ":00"
//				// convert date from string to date format
//				callDatetime, _ := time.Parse(constant.DateTimeFormat, callTimeString)
//				if paymentRepository.PublishInterfacePayment().CheckOrderPaymentExist(listenerId, callDatetime) {
//					responseItem.IsFull = true
//				} else {
//					scheduleWorkDetail, _ := repository.PublishInterfaceSchedule().GetDetailScheduleWorkListener(listenerId, date, timeSlot)
//					if scheduleWorkDetail.OrderStatus == "lock" && !(scheduleWorkDetail.UserLock == tokenInfo.Id) {
//						responseItem.IsFull = true
//					}
//				}
//			} else {
//				responseItem.IsFull = true
//			}
//			responseData = append(responseData, responseItem)
//		}
//
//		c.JSON(http.StatusOK, model.ResponseGetScheduleWorkAppointment{
//			Code:    http.StatusOK,
//			Message: message.MessageSuccess,
//			Data:    responseData,
//		})
//		return
//	}
//}
