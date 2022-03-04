package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"sandexcare_backend/api/listener/model"
	"sandexcare_backend/api/listener/repository"
	"sandexcare_backend/helpers/common"
	"sandexcare_backend/helpers/config"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/mail"
	"sandexcare_backend/helpers/message"
	"sandexcare_backend/helpers/middlewares"
	"strings"
	"time"
)

// @Summary API forgot password
// @Schemes
// @Description Xử lý quên mật khẩu!
// @Param input body model.InputForgotPassword true "Input Forgot Password"
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /listener/forgot-password [post]
func ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputForgotPassword{}
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

		// check phone number exist in listener table
		_, err = repository.PublishInterfaceListener().GetListenerByPhoneNumber(input.PhoneNumber)
		if err == nil {
			// Clear data reset password old
			_ = repository.PublishInterfaceListener().ClearDataResetPassword(input.PhoneNumber)

			code := common.GenerateNumber(constant.MaxLenOTPCode)
			contentSmsOTP := message.SmsOTP
			replacer := strings.NewReplacer("{otp_code}", code)
			contentSmsOTP = replacer.Replace(contentSmsOTP)

			dataForgotPassword := model.ListenerResetPassword{}
			dataForgotPassword.Id = primitive.NewObjectID()
			dataForgotPassword.PhoneNumber = input.PhoneNumber
			dataForgotPassword.ResetCode = code
			dataForgotPassword.Status = constant.Processing
			dataForgotPassword.ExpiresAt = time.Now().UTC().Add(time.Minute * time.Duration(constant.ExpiresOTP))
			dataForgotPassword.CreatedAt = time.Now().UTC()
			dataForgotPassword.UpdatedAt = time.Now().UTC()
			errCreate := repository.PublishInterfaceListener().CreateListenerResetPassword(dataForgotPassword)
			if errCreate == nil && common.SendSMS(input.PhoneNumber, contentSmsOTP) {
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": message.MessageSuccess,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorForgotPassword,
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPhoneNumberNotExist,
			})
		}
	}
}

// @Summary API resend forgot password
// @Schemes
// @Description Xử lý gửi lại quên mật khẩu!
// @Param input body model.InputForgotPassword true "Input Forgot Password"
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /listener/resend-forgot-password [post]
func ResendForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputForgotPassword{}
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

		// check phone number exist in listener table
		listenerInfo, err := repository.PublishInterfaceListener().GetListenerByPhoneNumber(input.PhoneNumber)
		if err == nil {
			// Clear data reset password old
			_ = repository.PublishInterfaceListener().ClearDataResetPassword(input.PhoneNumber)
			dataForgotPassword := model.ListenerResetPassword{}
			dataForgotPassword.Id = primitive.NewObjectID()
			dataForgotPassword.PhoneNumber = input.PhoneNumber
			dataForgotPassword.ResetCode = common.GenerateNumber(constant.MaxLenOTPCode)
			dataForgotPassword.Status = constant.Processing
			dataForgotPassword.ExpiresAt = time.Now().UTC().AddDate(100, 0, 0)
			dataForgotPassword.CreatedAt = time.Now().UTC()
			dataForgotPassword.UpdatedAt = time.Now().UTC()
			subject := mail.GetSubject(mail.TypeForgotPassword)
			content := mail.GetHtmlContentResetPassword(listenerInfo.Role, dataForgotPassword.PhoneNumber, dataForgotPassword.ResetCode)

			// get email admin system from file config
			env := config.GetEnvValue()
			check, _ := mail.SendEmail(env.Mail.Admin, subject, content)
			if check {
				err := repository.PublishInterfaceListener().CreateListenerResetPassword(dataForgotPassword)
				if err == nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    http.StatusOK,
						"message": message.MessageSuccess,
					})
				} else {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorForgotPassword,
					})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorForgotPassword,
				})
			}

		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPhoneNumberNotExist,
			})
		}
	}
}

// @Summary API verify reset password
// @Schemes
// @Description Xử lý xác nhận đặt lại mật khẩu!
// @Param input body model.InputVerifyResetPassword true "Input Verify Reset Password"
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /listener/verify-reset-password [post]
func VerifyResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputVerifyResetPassword{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		// get data request reset password
		data, err := repository.PublishInterfaceListener().GetListenerResetPasswordByPhoneNumberAndCode(input.PhoneNumber, input.ResetCode)
		if err == nil {
			if data.ExpiresAt.After(time.Now().UTC()) {
				data.ExpiresAt = time.Now().UTC()
				data.Status = constant.Verified
				data.UpdatedAt = time.Now().UTC()
				errUpdate := repository.PublishInterfaceListener().UpdateListenerResetPassword(data)
				if errUpdate == nil {
					c.JSON(http.StatusOK, gin.H{
						"code":    http.StatusOK,
						"message": message.MessageSuccess,
					})
				} else {
					c.JSON(http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"message": message.MessageErrorVerifyResetPassword,
					})
				}
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

// @Summary API reset new password
// @Schemes
// @Description Xử lý đặt lại mật khẩu mới!
// @Param input body model.InputResetNewPassword true "Input Reset New Password"
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /listener/reset-new-password [post]
func ResetNewPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputResetNewPassword{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		if !common.IsValidPasswordFormat(input.NewPassword) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPasswordFormat,
			})
			return
		}
		// check data reset password exist and verified
		dataReset, errDataReset := repository.PublishInterfaceListener().GetDataResetPassword(input.PhoneNumber, input.ResetCode, constant.Verified)
		if errDataReset == nil {
			//Hash password
			hashedPass := common.HashPassword(input.NewPassword)
			dataListener, err := repository.PublishInterfaceListener().GetListenerByPhoneNumber(input.PhoneNumber)
			if err == nil {
				dataListener.Password = hashedPass
				dataListener.UpdatedAt = time.Now().UTC()
				_ = repository.PublishInterfaceListener().UpdateListener(dataListener)
				dataReset.Status = constant.Completed
				dataReset.UpdatedAt = time.Now().UTC()
				_ = repository.PublishInterfaceListener().UpdateListenerResetPassword(dataReset)
				c.JSON(http.StatusOK, gin.H{
					"code":    http.StatusOK,
					"message": message.MessageSuccess,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": message.MessageErrorResetNewPassword,
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorResetNewPassword,
			})
		}
	}
}

// @Summary API change password
// @Schemes
// @Description Xử lý thay đổi mật khẩu!
// @Param input body model.InputChangePassword true "Input Change Password"
// @Tags Listener
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /listener/change-password [post]
func ChangePassword() gin.HandlerFunc {
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
		//check role is listener
		if tokenInfo.Role == constant.RoleUser {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		//get listener info by phone_number in token
		listener, _ := repository.PublishInterfaceListener().GetListenerByPhoneNumber(tokenInfo.PhoneNumber)

		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputChangePassword{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		//check validate
		if common.IsEmpty(input.Password) || common.IsEmpty(input.NewPassword) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPasswordEmpty,
			})
			return
		}

		if !common.IsValidPasswordFormat(input.NewPassword) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPasswordFormat,
			})
			return
		}

		if input.NewPassword == input.Password {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPasswordSameOld,
			})
			return
		}

		//hash password
		oldPassword := common.HashPassword(input.Password)
		newPassword := common.HashPassword(input.NewPassword)

		//check password old
		if listener.Password == oldPassword {
			listener.Password = newPassword
			listener.UpdatedAt = time.Now().UTC()
			_ = repository.PublishInterfaceListener().UpdateListener(listener)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": message.MessageSuccess,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorChangePassword,
			})
		}
	}
}
