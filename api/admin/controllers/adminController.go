package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"sandexcare_backend/api/admin/model"
	"sandexcare_backend/api/admin/repository"
	userRepository "sandexcare_backend/api/user/repository"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/message"
	"sandexcare_backend/helpers/middlewares"
	"time"
)

// @Summary API get all withdrawal history
// @Schemes
// @Description Admin xem tất cả yêu cầu rút tiền của cả hệ thống.
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetAllWithdrawalHistory
// @Failure 400 {object} model.ResponseError
// @Router  /admin/withdrawal-history [get]
func GetAllWithdrawalHistory() gin.HandlerFunc {
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
		// check role of token
		if tokenInfo.Role != constant.RoleAdmin {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}

		// get data from database
		data := repository.PublishInterfaceAdmin().GetAllWithdrawalHistory()
		c.JSON(http.StatusOK, model.ResponseGetAllWithdrawalHistory{
			Code: http.StatusOK,
			Message: message.MessageSuccess,
			Data: data,
		})
		return
	}
}

// @Summary API Confirm Withdrawal Request
// @Schemes
// @Description Admin xác nhận yêu cầu rút tiền của chuyên viên/chuyên gia!
// @Param input body model.InputConfirmWithdrawal true "Input Confirm Withdrawal"
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /admin/confirm-withdrawal [post]
func ConfirmWithdrawalRequest() gin.HandlerFunc {
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

		// check role of token
		if tokenInfo.Role != constant.RoleAdmin {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}

		rawBody, _ := c.GetRawData()
		input := model.InputConfirmWithdrawal{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate status
		if input.Status != constant.Completed && input.Status != constant.Canceled {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorStatusInvalid,
			})
			return
		}

		// get data from database
		dataWithdrawal, errWithdrawal := repository.PublishInterfaceAdmin().GetWithdrawalHistoryById(input.Id)
		if errWithdrawal == nil {
			dataWithdrawal.Status = input.Status
			dataWithdrawal.UpdatedAt = time.Now().UTC()
			_ = repository.PublishInterfaceAdmin().UpdateWithdrawalHistory(dataWithdrawal)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errWithdrawal.Error(),
			})
			return
		}
	}
}

// @Summary API get all user
// @Schemes
// @Description Admin xem tất cả người dùng trên hệ thống.
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetAllUsers
// @Failure 400 {object} model.ResponseError
// @Router  /admin/users [get]
func GetAllUser() gin.HandlerFunc {
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
		// check role of token
		if tokenInfo.Role != constant.RoleAdmin {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}

		// get data from database
		data := repository.PublishInterfaceAdmin().GetAllUsers()
		c.JSON(http.StatusOK, model.ResponseGetAllUsers{
			Code: http.StatusOK,
			Message: message.MessageSuccess,
			Data: data,
		})
		return
	}
}

// @Summary API Submit Diamonds User
// @Schemes
// @Description Admin cấp kim cương cho tài khoản user!
// @Param input body model.InputSubmitDiamond true "Input Submit Diamond"
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /admin/users/submit-diamond [post]
func SubmitDiamondsUser() gin.HandlerFunc {
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
		// check role of token
		if tokenInfo.Role != constant.RoleAdmin {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}

		rawBody, _ := c.GetRawData()
		input := model.InputSubmitDiamond{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		if input.Diamond <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDiamondInvalid,
			})
			return
		}
		// get user from database
		user, errUser := userRepository.PublishInterfaceUser().GetUserById(input.Id)
		if errUser == nil {
			user.Diamond = input.Diamond
			user.UpdatedAt = time.Now().UTC()
			_ = userRepository.PublishInterfaceUser().UpdateUser(user)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errUser.Error(),
			})
			return
		}
	}
}
