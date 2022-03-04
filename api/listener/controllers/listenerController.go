package controllers

import (
	"encoding/json"
	conversationRepository "sandexcare_backend/api/conversation/repository"
	"sandexcare_backend/api/listener/model"
	"sandexcare_backend/api/listener/repository"
	paymentRepository "sandexcare_backend/api/payment/repository"
	trackEventController "sandexcare_backend/api/track_event/controllers"
	"sandexcare_backend/helpers/common"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/mail"
	"sandexcare_backend/helpers/message"
	"sandexcare_backend/helpers/middlewares"
	wsControllers "sandexcare_backend/server_websocket/controllers"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	//"github.com/gorilla/websocket"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Create account listener
// @Schemes
// @Description Xử lý tạo account chuyên viên và chuyên gia!
// @Param input body model.InputCreateListener true "Add InputCreateListener"
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /listener/create [post]
func CreateListener() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputCreateListener{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate email
		if common.IsEmpty(input.Email) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorEmailEmpty,
			})
			return
		}
		if !common.CheckValidationEmail(input.Email) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorEmailInvalidFormat,
			})
			return
		}
		if !common.CheckLength(input.Email, constant.MinLengthEmail, constant.MaxLengthEmail) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorEmailLength,
			})
			return
		}

		if common.CheckEmailExist(input.Email, constant.RoleListener) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorEmailExist,
			})
			return
		}

		/*validate first name*/
		if common.IsEmpty(input.FirstName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorFirstNameEmpty,
			})
			return
		}
		if !common.CheckLength(input.FirstName, constant.MinLengthFirstName, constant.MaxLengthFirstName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorFirstNameLength,
			})
			return
		}
		if common.CheckSpecialCharacters(input.FirstName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorFirstNameSpecialChar,
			})
			return
		}

		/*validate last name*/
		if common.IsEmpty(input.LastName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorLastNameEmpty,
			})
			return
		}
		if !common.CheckLength(input.LastName, constant.MinLengthLastName, constant.MaxLengthLastName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorLastNameLength,
			})
			return
		}
		if common.CheckSpecialCharacters(input.LastName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorLastNameSpecialChar,
			})
			return
		}

		// validate birthday
		if common.IsEmpty(input.Dob) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorBirthdayEmpty,
			})
			return
		}

		if !common.CheckFormatDate(input.Dob) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorBirthdayInvalid,
			})
			return
		}

		// validate phone number
		if common.IsEmpty(input.PhoneNumber) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPhoneNumberEmpty,
			})
			return
		}

		if !common.CheckValidationPhoneNumber(input.PhoneNumber) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPhoneNumberInvalid,
			})
			return
		}

		if common.CheckPhoneNumberExist(input.PhoneNumber, constant.RoleListener) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPhoneNumberExist,
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

		if common.CheckEmployeeIdExist(input.EmployeeId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorEmployeeIdExist,
			})
			return
		}

		//validate role
		if input.Role != constant.RoleListener && input.Role != constant.RoleExperts {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorRoleInvalid,
			})
			return
		}
		// init model
		var listener model.Listener
		listener.Id = primitive.NewObjectID()
		listener.EmployeeId = input.EmployeeId
		listener.Email = input.Email
		listener.PhoneNumber = input.PhoneNumber
		listener.Name.FirstName = input.FirstName
		listener.Name.LastName = input.LastName
		listener.Dob = input.Dob
		listener.Gender = input.Gender
		listener.Address = input.Address
		listener.PersonalId = input.PersonalId
		listener.Bank.Owner = input.Owner
		listener.Bank.Account = input.Account
		listener.Bank.BankName = input.BankName
		listener.Role = input.Role
		listener.Status = constant.Active
		listener.Description = input.Description
		listener.MainTopic = input.MainTopic
		listener.Price = input.Price
		listener.Avatar = input.Avatar
		passwordGenerate := common.GenerateTokenString(12)
		listener.Password = common.HashPassword(passwordGenerate)
		listener.CreatedAt = time.Now().UTC()
		listener.UpdatedAt = time.Now().UTC()

		// handle save data
		errCreateListener := repository.PublishInterfaceListener().CreateListener(listener)
		if errCreateListener == nil {
			subject := mail.GetSubject(mail.TypeCreateListener)
			content := mail.GetHtmlContentCreateListener(input.Role, input.PhoneNumber, passwordGenerate)

			_, _ = mail.SendEmail(input.Email, subject, content)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errCreateListener.Error(),
			})
			return
		}
	}
}

// @Summary API login
// @Schemes
// @Description Xử lý chuyên viên/chuyên gia đăng nhập!
// @Param input body model.InputLogin true "Input handle login"
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseLogin
// @Failure 400 {object} model.ResponseError
// @Router  /listener/login [post]
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		inputLogin := model.InputLogin{}
		err := json.Unmarshal(rawBody, &inputLogin)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		phoneNumber := inputLogin.PhoneNumber
		password := inputLogin.Password

		//hash password
		hashedPass := common.HashPassword(password)

		//get data listener
		listener, errListener := repository.PublishInterfaceListener().Login(phoneNumber, hashedPass)
		if errListener == nil {
			if listener.Status != constant.Active {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorAccountInacctive,
				})
				return
			}
			var token, _ = middlewares.GenerateJWT(listener.Id.Hex(), listener.PhoneNumber, listener.Email, listener.Role)
			data := model.LoginData{}
			data.Id = listener.Id.Hex()
			data.PhoneNumber = listener.PhoneNumber
			data.Token = token
			data.Email = listener.Email

			listener.CountLoginFail = 0
			listener.UpdatedAt = time.Now().UTC()
			_ = repository.PublishInterfaceListener().UpdateListener(listener)
			ReconnectCallInProgress(listener.Id.Hex())
			c.JSON(http.StatusOK, model.ResponseLogin{
				Code:    http.StatusOK,
				Message: message.MessageSuccess,
				Data:    data,
			})
		} else {
			dataListener, errListener := repository.PublishInterfaceListener().GetListenerByPhoneNumber(inputLogin.PhoneNumber)
			if errListener == nil {
				if dataListener.CountLoginFail < 3 {
					dataListener.CountLoginFail = dataListener.CountLoginFail + 1
					_ = repository.PublishInterfaceListener().UpdateListener(dataListener)
				} else {
					c.JSON(http.StatusMethodNotAllowed, gin.H{
						"code":    http.StatusMethodNotAllowed,
						"message": message.MessageNotifyResetPassword,
					})
					return
				}
			}
			/*Login fail*/
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorLoginFail,
			})
		}
	}
}

// @Summary API login V2
// @Schemes
// @Description Xử lý chuyên viên/chuyên gia đăng nhập!
// @Param input body model.InputLogin true "Input handle login"
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseLogin
// @Failure 400 {object} model.ResponseError
// @Router  /listener/login2 [post]
func Login2() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		inputLogin := model.InputLogin{}
		err := json.Unmarshal(rawBody, &inputLogin)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		phoneNumber := inputLogin.PhoneNumber
		password := inputLogin.Password

		//hash password
		hashedPass := common.HashPassword(password)

		//get data listener
		listener, errListener := repository.PublishInterfaceListener().Login(phoneNumber, hashedPass)
		if errListener == nil {
			if listener.Status != constant.Active {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorAccountInacctive,
				})
				return
			}

			token, err := middlewares.GenerateJWT(listener.Id.Hex(), listener.PhoneNumber, listener.Email, listener.Role)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"code":    http.StatusBadGateway,
					"message": message.InvalidToken,
				})
				return
			}

			rtoken, err := middlewares.GenerateRefreshToken(listener.Id.Hex())
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"code":    http.StatusBadGateway,
					"message": message.InvalidToken,
				})
				return
			}

			data := model.LoginData{}
			data.Id = listener.Id.Hex()
			data.PhoneNumber = listener.PhoneNumber
			data.Token = token
			data.RefreshToken = rtoken
			data.Email = listener.Email

			listener.CountLoginFail = 0
			listener.RefreshToken = strings.Split(rtoken, ".")[2]
			listener.UpdatedAt = time.Now().UTC()
			_ = repository.PublishInterfaceListener().UpdateListener(listener)
			ReconnectCallInProgress(listener.Id.Hex())
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
				"data":    data,
			})
		} else {
			dataListener, errListener := repository.PublishInterfaceListener().GetListenerByPhoneNumber(inputLogin.PhoneNumber)
			if errListener == nil {
				if dataListener.CountLoginFail < 3 {
					dataListener.CountLoginFail = dataListener.CountLoginFail + 1
					_ = repository.PublishInterfaceListener().UpdateListener(dataListener)
				} else {
					c.JSON(http.StatusMethodNotAllowed, gin.H{
						"code":    http.StatusMethodNotAllowed,
						"message": message.MessageNotifyResetPassword,
					})
					return
				}
			}
			/*Login fail*/
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorLoginFail,
			})
		}
	}
}

// @Summary
// @Schemes
// @Description API Verify Token
// @Param Authorization header string true "Refresh Token: Prefix Token xxxxxx"
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseVerifyOTPCode
// @Failure 400 {object} model.ResponseError
// @Failure 403 {object} model.ResponseError
// @Router  /listener/verify-token [get]
func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// responseData := make(map[string]string)
		// responseData["access_token"] = token
		// responseData["refresh_token"] = refreshtoken
		uuid, signature, err := middlewares.GetIDByRefreshToken(c)
		logrus.Info(uuid)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    err.Error(),
				"message": message.InvalidRefreshToken,
			})
			return
		}

		//get data listener
		listener, errListener := repository.PublishInterfaceListener().GetListenerByListenerId(uuid)
		if errListener == nil {
			if listener.Status != constant.Active {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorAccountInacctive,
				})
				return
			}

			if listener.RefreshToken != signature {
				c.JSON(http.StatusForbidden, gin.H{
					"code":    http.StatusForbidden,
					"message": message.InvalidRefreshToken,
				})
				return
			}
			token, err := middlewares.GenerateJWT(listener.Id.Hex(), listener.PhoneNumber, listener.Email, listener.Role)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"code":    http.StatusBadGateway,
					"message": message.InvalidToken,
				})
				return
			}

			rtoken, err := middlewares.GenerateRefreshToken(listener.Id.Hex())
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"code":    http.StatusBadGateway,
					"message": message.InvalidToken,
				})
				return
			}

			data := make(map[string]string)

			data["token"] = token
			data["refresh_token"] = rtoken

			listener.CountLoginFail = 0
			listener.RefreshToken = strings.Split(rtoken, ".")[2]
			listener.UpdatedAt = time.Now().UTC()
			_ = repository.PublishInterfaceListener().UpdateListener(listener)
			ReconnectCallInProgress(listener.Id.Hex())
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
				"data":    data,
			})
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": message.InvalidRefreshToken,
			})
		}
	}
}

// @Summary API Logout
// @Schemes
// @Description Xử lý đăng xuất!
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /listener/logout [post]
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get token info
		_, errToken := middlewares.AuthenticateToken(c)
		if errToken == nil {
			cookie := http.Cookie{
				Name:   "token",
				MaxAge: -1}
			http.SetCookie(c.Writer, &cookie)

			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errToken.Error(),
			})
			return
		}
	}
}

// @Summary API get listener info
// @Schemes
// @Description Xử lý lấy thông tin của chuyên viên/chuyên gia bằng token
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetInfoListener
// @Failure 400 {object} model.ResponseError
// @Router  /listener/info [get]
func GetInfoListener() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get token info
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken == nil {
			if tokenInfo.Role == constant.RoleUser {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorPermissionDenied,
				})
				return
			}
			ReconnectCallInProgress(tokenInfo.Id)
			listener, _ := repository.PublishInterfaceListener().GetListenerByListenerId(tokenInfo.Id)
			listener.Password = ""
			c.JSON(http.StatusOK, model.ResponseGetInfoListener{
				Code:    http.StatusOK,
				Message: message.MessageSuccess,
				Data:    listener,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errToken.Error(),
			})
			return
		}
	}
}

// @Summary API Yêu cầu rút tiền
// @Schemes
// @Description Xử lý gửi yều cầu rút tiền!
// @Param input body model.InputRequestWithdrawal true "Input Request Withdrawal"
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /listener/request-withdrawal [post]
func RequestWithdrawal() gin.HandlerFunc {
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

		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputRequestWithdrawal{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		if input.Money <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDiamondInvalid,
			})
			return
		}

		remainingDiamond := repository.PublishInterfaceListener().GetTotalRemainingWithdrawalDiamond(tokenInfo.Id)
		if remainingDiamond < input.Money {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorDiamondNotEnough,
			})
			return
		}

		numberDiamond := input.Money
		remainingWithdrawalDetail := repository.PublishInterfaceListener().GetRemainingWithdrawalDiamondDetail(tokenInfo.Id)
		for _, v := range remainingWithdrawalDetail {
			if numberDiamond == 0 {
				break
			}
			if numberDiamond >= v.WithdrawalDiamond {
				numberDiamond = numberDiamond - v.WithdrawalDiamond
				v.WithdrawalDiamond = 0
			} else {
				v.WithdrawalDiamond = v.WithdrawalDiamond - numberDiamond
				numberDiamond = 0
			}
			_ = conversationRepository.PublishInterfaceConversation().UpdateCallHistory(v)
		}

		// init model
		var request model.WithdrawalHistory
		request.Id = primitive.NewObjectID()
		request.ListenerId = tokenInfo.Id
		request.Status = constant.Processing
		request.AmountMoney = input.Money
		request.CreatedAt = time.Now().UTC()
		request.UpdatedAt = time.Now().UTC()

		// handle save data
		errRequest := repository.PublishInterfaceListener().CreateRequestWithdrawal(request)
		if errRequest == nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errRequest.Error(),
			})
			return
		}
	}
}

// @Summary API get withdrawal history
// @Schemes
// @Description Xử lý lấy lịch xử yêu cầu rút tiền
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetWithdrawalHistory
// @Failure 400 {object} model.ResponseError
// @Router  /listener/withdrawal-history [get]
func GetWithdrawalHistory() gin.HandlerFunc {
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
		if tokenInfo.Role == constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}

		// get data from database
		data := repository.PublishInterfaceListener().GetWithdrawalHistory(tokenInfo.Id)
		c.JSON(http.StatusOK, model.ResponseGetWithdrawalHistory{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    data,
		})
		return
	}
}

// @Summary API Get Revenue Analysis
// @Schemes
// @Description Phân tích doanh số, tổng số cuôc gọi,...
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetWithdrawalHistory
// @Failure 400 {object} model.ResponseError
// @Router  /listener/revenue-analysis [get]
func GetRevenueAnalysis() gin.HandlerFunc {
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
		if tokenInfo.Role == constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		// init model response
		data := model.GetRevenueAnalysisData{}
		data.CallNumberSuccess = repository.PublishInterfaceListener().CountCallNumberListenerByStatus(tokenInfo.Id, constant.Completed)
		data.CallNumberFail = repository.PublishInterfaceListener().CountCallNumberListenerByStatus(tokenInfo.Id, constant.Unconnected)
		data.Rating = repository.PublishInterfaceListener().AggregateStarRatingListener(tokenInfo.Id, true)
		data.TotalRevenue = repository.PublishInterfaceListener().GetTotalRevenueDiamondListener(tokenInfo.Id)
		data.AvailableWithdrawal = repository.PublishInterfaceListener().GetTotalRemainingWithdrawalDiamond(tokenInfo.Id)
		c.JSON(http.StatusOK, model.ResponseGetRevenueAnalysis{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    data,
		})
		return
	}
}

// @Summary API handle miss call
// @Schemes
// @Description Xử lý chuyên viên/chuyên gia không bắt máy!
// @Param input body model.InputHandleMissCall true "Input Handle Miss Call"
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /listener/miss-call [post]
func HandleMissCall() gin.HandlerFunc {
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
		input := model.InputHandleMissCall{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		callData, errCall := conversationRepository.PublishInterfaceConversation().GetCallHistoryById(input.CallId)
		if errCall != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errCall.Error(),
			})
			return
		} else {
			if callData.Status != constant.Started || callData.ListenerId != tokenInfo.Id {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorPermissionDenied,
				})
				return
			}
			callData.Status = constant.Unconnected
			callData.UpdatedAt = time.Now().UTC()
			_ = conversationRepository.PublishInterfaceConversation().UpdateCallHistory(callData)

			listenerInfo, _ := repository.PublishInterfaceListener().GetListenerByListenerId(callData.ListenerId)
			listenerInfo.Status = constant.Inactive
			listenerInfo.UpdatedAt = time.Now().UTC()
			_ = repository.PublishInterfaceListener().UpdateListener(listenerInfo)
			trackEventController.CreateTrackCall(callData.ListenerId, input.CallId, constant.MissCall)

			subject := mail.GetSubject(mail.TypeNotifyBlockListener)
			content := mail.GetHtmlContentBlockListener()

			_, _ = mail.SendEmail(tokenInfo.Email, subject, content)

			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
			return
		}
	}
}

// Reconnect call in progress of listener
func ReconnectCallInProgress(listenerId string) {
	// check the call in progress
	callData, errCall := conversationRepository.PublishInterfaceConversation().GetCallInProgressOfListener(listenerId)
	if errCall == nil {
		orderInfo, _ := paymentRepository.PublishInterfacePayment().GetOrderPaymentById(callData.OrderId)
		//userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(callData.UserId)
		hub := wsControllers.GetHub()
		payload := wsControllers.SocketEventStruct{}
		payload.EventName = constant.WSEventReconnectCall
		payload.EventPayload = map[string]interface{}{
			"call_id":      callData.Id.Hex(),
			"username":     " userInfo.Name",
			"date":         orderInfo.Date.Format(constant.DateFormat),
			"booking_time": orderInfo.BookingTime,
		}
		wsControllers.EmitToSpecificClient(hub, payload, callData.ListenerId)
	}
}
