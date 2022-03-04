package controllers

import (
	"encoding/json"
	listenerRepository "sandexcare_backend/api/listener/repository"
	masterModel "sandexcare_backend/api/master/model"
	paymentRepository "sandexcare_backend/api/payment/repository"
	"sandexcare_backend/api/user/model"
	"sandexcare_backend/api/user/repository"
	"sandexcare_backend/helpers/common"
	"sandexcare_backend/helpers/config"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/mail"
	"sandexcare_backend/helpers/message"
	"sandexcare_backend/helpers/middlewares"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	//"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary
// @Schemes
// @Description API Send OTP Code
// @Param input body model.InputSendOTPCode true "Input Send OTP code"
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSendOTPCode
// @Failure 400 {object} model.ResponseError
// @Router  /user/send-otp-code [post]
func SendOTPCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputSendOTPCode{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
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
		dataUser, errUser := repository.PublishInterfaceUser().GetUserByPhoneNumber(input.PhoneNumber)

		var code string
		if input.PhoneNumber == "0999999990" || input.PhoneNumber == "0999999991" || input.PhoneNumber == "0999999992" {
			code = "111111"
		} else {
			code = common.GenerateNumber(constant.MaxLenOTPCode)
		}
		contentSmsOTP := message.SmsOTP
		replacer := strings.NewReplacer("{otp_code}", code)
		contentSmsOTP = replacer.Replace(contentSmsOTP)

		if errUser == nil {
			dataUser.Code = code
			dataUser.ExpiresAt = time.Now().UTC().Add(time.Minute * time.Duration(constant.ExpiresOTP))
			dataUser.UpdatedAt = time.Now().UTC()
			errUpdate := repository.PublishInterfaceUser().UpdateUser(dataUser)
			if errUpdate == nil && common.SendSMS(dataUser.PhoneNumber, contentSmsOTP) {
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": message.MessageSuccess,
					"data":    code,
				})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorSendOTP,
				})
				return
			}
		} else {
			dataUser := model.User{}
			dataUser.Id = primitive.NewObjectID()
			dataUser.PhoneNumber = input.PhoneNumber
			dataUser.Code = code
			dataUser.Status = constant.Active
			dataUser.ExpiresAt = time.Now().UTC().Add(time.Minute * time.Duration(constant.ExpiresOTP))
			dataUser.CreatedAt = time.Now().UTC()
			dataUser.UpdatedAt = time.Now().UTC()
			errCreate := repository.PublishInterfaceUser().CreateUser(dataUser)
			if errCreate == nil && common.SendSMS(input.PhoneNumber, contentSmsOTP) {
				c.JSON(http.StatusOK, model.ResponseSendOTPCode{
					Code:    http.StatusOK,
					Message: message.MessageSuccess,
					Data:    code,
				})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorSendOTP,
				})
				return
			}
		}
	}
}

// @Summary
// @Schemes
// @Description API Resend OTP Code
// @Param input body model.InputSendOTPCode true "Input resend OTP code"
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSendOTPCode
// @Failure 400 {object} model.ResponseError
// @Router  /user/resend-otp-code [post]
func ResendOTPCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputSendOTPCode{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
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

		code := common.GenerateNumber(constant.MaxLenOTPCode)
		dataUser, errUser := repository.PublishInterfaceUser().GetUserByPhoneNumber(input.PhoneNumber)
		if errUser == nil {
			dataUser.Code = code
			dataUser.ExpiresAt = time.Now().UTC().AddDate(100, 0, 0)
			dataUser.UpdatedAt = time.Now().UTC()
			_ = repository.PublishInterfaceUser().UpdateUser(dataUser)
			subject := mail.GetSubject(mail.TypeForgotPassword)
			content := mail.GetHtmlContentResetPassword(constant.RoleUser, input.PhoneNumber, code)

			// get email admin system from file config
			env := config.GetEnvValue()
			check, _ := mail.SendEmail(env.Mail.Admin, subject, content)
			if !check {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorResendOTP,
				})
				return
			}
			c.JSON(http.StatusOK, model.ResponseSendOTPCode{
				Code:    http.StatusOK,
				Message: message.MessageSuccess,
				Data:    code,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, masterModel.ResponseError{
				Code:    http.StatusBadRequest,
				Message: message.MessageErrorPhoneNumberNotExist,
			})
			return
		}
	}
}

// @Summary
// @Schemes
// @Description API Verify OTP Code
// @Param input body model.InputVerifyOTPCode true "Input Verify OTP code"
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseVerifyOTPCode
// @Failure 400 {object} model.ResponseError
// @Router  /user/verify-otp-code [post]
func VerifyOTPCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputVerifyOTPCode{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
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

		// validate code
		if common.IsEmpty(input.Code) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorCodeEmpty,
			})
			return
		}

		// get data user by phone_number and code
		dataUser, errUser := repository.PublishInterfaceUser().GetUserByPhoneNumberAndCode(input.PhoneNumber, input.Code)
		if errUser == nil {
			if dataUser.ExpiresAt.After(time.Now().UTC()) {
				dataUser.ExpiresAt = time.Now().UTC()
				dataUser.UpdatedAt = time.Now().UTC()
				if dataUser.Status != constant.Active {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorAccountInacctive,
					})
					return
				}
				_ = repository.PublishInterfaceUser().UpdateUser(dataUser)
				var token, _ = middlewares.GenerateJWT(dataUser.Id.Hex(), dataUser.PhoneNumber, "", constant.RoleUser)
				responseData := model.ResponseVerifyOTPCode{}
				responseData.Id = dataUser.Id.Hex()
				responseData.PhoneNumber = dataUser.PhoneNumber
				responseData.IsMember = dataUser.IsMember
				responseData.Token = token
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": message.MessageSuccess,
					"data":    responseData,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorCodeExpired,
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorCodeIncorrect,
			})
		}
	}
}

// @Summary
// @Schemes
// @Description API Verify OTP Code V2
// @Param input body model.InputVerifyOTPCode true "Input Verify OTP code"
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseVerifyOTPCode
// @Failure 400 {object} model.ResponseError
// @Router  /user/verify-otp-code2 [post]
func VerifyOTPCode2() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputVerifyOTPCode{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
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

		// validate code
		if common.IsEmpty(input.Code) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorCodeEmpty,
			})
			return
		}

		// get data user by phone_number and code
		dataUser, errUser := repository.PublishInterfaceUser().GetUserByPhoneNumberAndCode(input.PhoneNumber, input.Code)
		if errUser == nil {
			if dataUser.ExpiresAt.After(time.Now().UTC()) {
				dataUser.ExpiresAt = time.Now().UTC()
				dataUser.UpdatedAt = time.Now().UTC()
				if dataUser.Status != constant.Active {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorAccountInacctive,
					})
					return
				}

				token, err := middlewares.GenerateJWT(dataUser.Id.Hex(), dataUser.PhoneNumber, "", constant.RoleUser)
				if err != nil {
					logrus.Error(err)
					c.JSON(http.StatusBadGateway, gin.H{
						"code":    "E_Abnormal",
						"message": message.Abnormal,
					})
				}

				//
				refreshtoken, err := middlewares.GenerateRefreshToken(dataUser.Id.Hex())
				if err != nil {
					logrus.Error(err)
					c.JSON(http.StatusBadGateway, gin.H{
						"code":    "E_Abnormal",
						"message": message.Abnormal,
					})
				}

				dataUser.RefreshToken = strings.Split(refreshtoken, ".")[2]
				err = repository.PublishInterfaceUser().UpdateUser(dataUser)
				if err != nil {
					logrus.Error(err)
					c.JSON(http.StatusBadGateway, gin.H{
						"code":    "E_Abnormal",
						"message": message.Abnormal,
					})
				}
				responseData := model.ResponseVerifyOTPCode{}
				responseData.Id = dataUser.Id.Hex()
				responseData.PhoneNumber = dataUser.PhoneNumber
				responseData.IsMember = dataUser.IsMember
				responseData.Token = token
				responseData.RefreshToken = refreshtoken
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": message.MessageSuccess,
					"data":    responseData,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorCodeExpired,
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorCodeIncorrect,
			})
		}
	}
}

// @Summary
// @Schemes
// @Description API Verify Token
// @Param Authorization header string true "Refresh Token: Prefix Token xxxxxx"
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseVerifyOTPCode
// @Failure 400 {object} model.ResponseError
// @Failure 403 {object} model.ResponseError
// @Router  /user/verify-token [get]
func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, signature, err := middlewares.GetIDByRefreshToken(c)
		logrus.Info(uuid)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    err.Error(),
				"message": message.InvalidRefreshToken,
			})
			return
		}
		//
		// get data user by phone_number and code
		u, err := repository.PublishInterfaceUser().GetUserById(uuid)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.InvalidRefreshToken,
			})
			return
		}

		if u.RefreshToken != signature {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": message.InvalidRefreshToken,
			})
			return
		}
		u.ExpiresAt = time.Now().UTC()
		u.UpdatedAt = time.Now().UTC()
		if u.Status != constant.Active {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorAccountInacctive,
			})
			return
		}

		token, err := middlewares.GenerateJWT(u.Id.Hex(), u.PhoneNumber, "", constant.RoleUser)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadGateway, gin.H{
				"code":    "E_Abnormal",
				"message": message.Abnormal,
			})
		}

		//
		refreshtoken, err := middlewares.GenerateRefreshToken(u.Id.Hex())
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadGateway, gin.H{
				"code":    "E_Abnormal",
				"message": message.Abnormal,
			})
		}
		//
		u.RefreshToken = strings.Split(refreshtoken, ".")[2]
		err = repository.PublishInterfaceUser().UpdateUser(u)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusBadGateway, gin.H{
				"code":    "E_Abnormal",
				"message": message.Abnormal,
			})
		}
		responseData := make(map[string]string)
		responseData["token"] = token
		responseData["refresh_token"] = refreshtoken

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
			"data":    responseData,
		})
	}
}

// @Summary
// @Schemes
// @Description API User complete info
// @Param input body model.InputCompleteInfo true "Input complete info for user"
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /user/complete-info [post]
func CompleteInfo() gin.HandlerFunc {
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
		//check role is user
		if tokenInfo.Role != constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputCompleteInfo{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate name
		if common.IsEmpty(input.Name) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorNameEmpty,
			})
			return
		}

		if !common.CheckLength(input.Name, constant.MinLengthName, constant.MaxLengthName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorNameLength,
			})
			return
		}

		if common.CheckSpecialCharacters(input.Name) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorNameCharSpecial,
			})
			return
		}
		// get data user by id
		dataUser, errUser := repository.PublishInterfaceUser().GetUserById(tokenInfo.Id)
		if errUser == nil {
			if dataUser.IsMember {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorPermissionDenied,
				})
				return
			}
			if !common.IsEmpty(input.ReferralCode) {
				referralInfo, errReferral := repository.PublishInterfaceUser().GetUserByPhoneNumber(input.ReferralCode)
				if errReferral != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorReferralNotExist,
					})
					return
				} else {
					if referralInfo.CountReferral < 5 {
						dateExpire := 7
						// tặng mã giảm giá cho người giới thiệu
						common.InitCouponForUser(referralInfo.Id.Hex(), 0.5, dateExpire, constant.CouponBookingCV)
						referralInfo.CountReferral += 1
						referralInfo.UpdatedAt = time.Now().UTC()
						_ = repository.PublishInterfaceUser().UpdateUser(referralInfo)
						// tặng mã giảm giá cho account mới
						for i := 0; i < 5; i++ {
							common.InitCouponForUser(dataUser.Id.Hex(), 0.5, dateExpire, constant.CouponBookingCV)
							dateExpire += 7
						}
					}
				}
			}

			dataUser.Name = input.Name
			dataUser.ConsultingFields = input.ConsultingFields
			dataUser.IsMember = true
			dataUser.UpdatedAt = time.Now().UTC()
			_ = repository.PublishInterfaceUser().UpdateUser(dataUser)
			common.InitCouponForUser(dataUser.Id.Hex(), 1, 30, constant.CouponCallNow)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errUser.Error(),
			})
		}
	}
}

// @Summary
// @Schemes
// @Description API get info by token
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /user/info [get]
func GetInfoUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get token info
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken == nil {
			if tokenInfo.Role != constant.RoleUser {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorPermissionDenied,
				})
				return
			}
			user, _ := repository.PublishInterfaceUser().GetUserById(tokenInfo.Id)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
				"data":    user,
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

// @Summary
// @Schemes
// @Description API update complete info of user
// @Param input body model.InputUpdateInfo true "Input update info for user"
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /user/update-info [post]
func UpdateUserInfo() gin.HandlerFunc {
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
		//check role is user
		if tokenInfo.Role != constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputUpdateInfo{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// validate name
		if common.IsEmpty(input.Name) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorNameEmpty,
			})
			return
		}

		if !common.CheckLength(input.Name, constant.MinLengthName, constant.MaxLengthName) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorNameLength,
			})
			return
		}

		if common.CheckSpecialCharacters(input.Name) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorNameCharSpecial,
			})
			return
		}

		// get data user by id
		dataUser, errUser := repository.PublishInterfaceUser().GetUserById(tokenInfo.Id)
		if errUser == nil {
			// update avatar user
			dataUser.Avatar = input.Avatar
			dataUser.UpdatedAt = time.Now().UTC()
			_ = repository.PublishInterfaceUser().UpdateUser(dataUser)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": errUser.Error(),
			})
		}
	}
}

// @Summary
// @Schemes
// @Description API handle user bookmark listener
// @Param input body model.InputBookmarkListener true "Input user bookmark listener"
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /user/listeners-bookmark [post]
func CreateListenersBookmark() gin.HandlerFunc {
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
		//check role is user
		if tokenInfo.Role != constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputBookmarkListener{}
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
				"message": message.MessageErrorNameEmpty,
			})
			return
		}

		listenerInfo, errListener := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(input.ListenerId)
		// get data listener by id
		if errListener == nil {
			//check listener bookmark exist
			_, errExist := repository.PublishInterfaceUser().GetListenersBookmarkByUserIdAndListenerId(tokenInfo.Id, input.ListenerId)
			if errExist != nil {
				// listener bookmark has not existed
				//init model
				bookmark := model.ListenersBookmark{}
				bookmark.Id = primitive.NewObjectID()
				bookmark.ListenerId = input.ListenerId
				bookmark.ListenerRole = listenerInfo.Role
				bookmark.UserId = tokenInfo.Id
				bookmark.CreatedAt = time.Now().UTC()
				bookmark.UpdatedAt = time.Now().UTC()
				_ = repository.PublishInterfaceUser().CreateListenersBookmark(bookmark)
			}
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorEmployeeIdNotExist,
			})
		}
	}
}

// @Summary
// @Schemes
// @Description API get listeners user's bookmark
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetListenersBookmark
// @Failure 400 {object} model.ResponseError
// @Router  /user/listeners-bookmark [get]
func GetListenersBookmark() gin.HandlerFunc {
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
		//check role is user
		if tokenInfo.Role != constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		data := repository.PublishInterfaceUser().GetAllListenersBookmark(tokenInfo.Id)
		c.JSON(http.StatusOK, model.ResponseGetListenersBookmark{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    data,
		})
	}
}

// @Summary
// @Schemes
// @Description API delete listener bookmark of user
// @Param  listener_id path string true "Listener ID."
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /user/listeners-bookmark [delete]
func DeleteListenersBookmark() gin.HandlerFunc {
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
		//check role is user
		if tokenInfo.Role != constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		listenerId := c.Query("listener_id")

		if common.IsEmpty(listenerId) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorListenerIdEmpty,
			})
			return
		}

		_ = repository.PublishInterfaceUser().DeleteListenerBookmark(tokenInfo.Id, listenerId)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
		})
	}
}

// @Summary
// @Schemes
// @Description API get user's coupons
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseGetCouponsUser
// @Failure 400 {object} model.ResponseError
// @Router  /user/coupons [get]
func GetCouponsUser() gin.HandlerFunc {
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
		//check role is user
		if tokenInfo.Role != constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		couponsUser := repository.PublishInterfaceUser().GetAllCouponsUser(tokenInfo.Id)
		c.JSON(http.StatusOK, model.ResponseGetCouponsUser{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    couponsUser,
		})
	}
}

// @Summary
// @Schemes
// @Description Get schedule an upcoming appointment of user!
// @Tags User
// @Accept json
// @Produce json
// @Header 200 {string} Token ""
// @Success 200 {object} model.ResponseAppointmentSchedule
// @Failure 400 {object} model.ResponseError
// @Router  /user/appointment-schedule [get]
func GetAppointmentScheduleForUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token info
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusBadRequest, masterModel.ResponseError{
				Code:    http.StatusBadRequest,
				Message: errToken.Error(),
			})
			return
		}
		//check role is user
		if tokenInfo.Role != constant.RoleUser {
			c.JSON(http.StatusBadRequest, masterModel.ResponseError{
				Code:    http.StatusBadRequest,
				Message: message.MessageErrorPermissionDenied,
			})
			return
		}
		dataAppointment := paymentRepository.PublishInterfacePayment().GetScheduleAppointmentForUser(tokenInfo.Id)
		var responseData []model.DataAppointmentSchedule
		for _, v := range dataAppointment {
			item := model.DataAppointmentSchedule{}
			item.Id = v.Id.Hex()
			listenerInfo, _ := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(v.ListenerId)
			item.ListenerName = common.GetFullNameOfListener(listenerInfo)
			item.ListenerId = listenerInfo.Id.Hex()
			item.EmployeeId = listenerInfo.EmployeeId
			item.ListenerRole = listenerInfo.Role
			item.ListenerImage = listenerInfo.Avatar
			item.TimeSlot = v.TimeSlot
			item.Date = v.Date.Format(constant.DateFormat)
			item.BookingTime = v.BookingTime

			responseData = append(responseData, item)
		}
		c.JSON(http.StatusOK, model.ResponseAppointmentSchedule{
			Code:    http.StatusOK,
			Message: message.MessageSuccess,
			Data:    responseData,
		})
		return
	}
}

// @Summary
// @Schemes
// @Description Get detail schedule an upcoming appointment of user!
// @Tags User
// @Accept json
// @Produce json
// @Param  order_id path string true "Appoinment Id."
// @Header 200 {string} Token "jhfhhsid9834e8ff39fh"
// @Success 200 {object} model.ResponseDetailAppointmentSchedule
// @Failure 400 {object} model.ResponseError
// @Router  /user/appointment-schedule/detail [get]
func GetDetailAppointmentSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token info
		tokenInfo, errToken := middlewares.AuthenticateToken(c)
		if errToken != nil {
			c.JSON(http.StatusBadRequest, masterModel.ResponseError{
				Code:    http.StatusBadRequest,
				Message: errToken.Error(),
			})
			return
		}
		//check role is user
		if tokenInfo.Role != constant.RoleUser {
			c.JSON(http.StatusBadRequest, masterModel.ResponseError{
				Code:    http.StatusBadRequest,
				Message: message.MessageErrorPermissionDenied,
			})
			return
		}

		orderId := c.Query("order_id")
		// validate id
		if common.IsEmpty(orderId) {
			c.JSON(http.StatusBadRequest, masterModel.ResponseError{
				Code:    http.StatusBadRequest,
				Message: message.MessageErrorOrderIdEmpty,
			})
			return
		}

		detailAppointment, errDetail := paymentRepository.PublishInterfacePayment().GetOrderPaymentById(orderId)
		if errDetail == nil {
			responseData := model.DataDetailAppointmentSchedule{}
			listenerInfo, _ := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(detailAppointment.ListenerId)
			responseData.ListenerId = listenerInfo.Id.Hex()
			responseData.EmployeeId = listenerInfo.EmployeeId
			responseData.ListenerName = common.GetFullNameOfListener(listenerInfo)
			responseData.ListenerRole = listenerInfo.Role
			responseData.ListenerImage = listenerInfo.Avatar
			responseData.TimeSlot = detailAppointment.TimeSlot
			responseData.Date = detailAppointment.Date.Format(constant.DateFormat)
			responseData.BookingTime = detailAppointment.BookingTime
			responseData.DiamondOrder = detailAppointment.DiamondPayment
			responseData.DiamondDiscount = detailAppointment.DiamondPayment - detailAppointment.DiamondPayment
			responseData.DiamondPayment = detailAppointment.DiamondPayment
			c.JSON(http.StatusOK, model.ResponseDetailAppointmentSchedule{
				Code:    http.StatusOK,
				Message: message.MessageSuccess,
				Data:    responseData,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, masterModel.ResponseError{
				Code:    http.StatusBadRequest,
				Message: errDetail.Error(),
			})
			return
		}
	}
}
