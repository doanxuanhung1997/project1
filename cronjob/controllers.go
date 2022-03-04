package cronjob

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	listenerRepository "sandexcare_backend/api/listener/repository"
	notificationModel "sandexcare_backend/api/notification/model"
	notificationRepository "sandexcare_backend/api/notification/repository"
	paymentRepository "sandexcare_backend/api/payment/repository"
	userRepository "sandexcare_backend/api/user/repository"
	"sandexcare_backend/helpers/common"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/message"
	wsControllers "sandexcare_backend/server_websocket/controllers"
	"strconv"
	"strings"
	"time"
)

func UpcomingAppointmentReminders() {
	dataAppointment := paymentRepository.PublishInterfacePayment().GetUpcomingAppointment()
	timeNow, _ := time.Parse(constant.DateTimeFormat, time.Now().UTC().Add(constant.UTC7*time.Hour).Format(constant.DateTimeFormat))
	for _, a := range dataAppointment {
		userInfo, _ := userRepository.PublishInterfaceUser().GetUserById(a.UserId)
		listenerInfo, _ := listenerRepository.PublishInterfaceListener().GetListenerByListenerId(a.ListenerId)
		flagSendWs := false
		minute := math.Round(a.CallDatetime.Sub(timeNow).Minutes())
		minuteString := strconv.FormatFloat(minute, 'f', 0, 64)
		switch minute {
		case 30:
			//Send SMS cho user, chuyên gia
			contentUser := message.SmsRemind
			dayOfWeek := common.GetDayOfWeek(a.Date)
			dateSMS := a.Date.Format("02/01/2006")
			replacer := strings.NewReplacer("{hour}", a.BookingTime, "{day_of_week}", dayOfWeek, "{date}", dateSMS)
			contentUser = replacer.Replace(contentUser)
			common.SendSMS(userInfo.PhoneNumber, contentUser)
			if listenerInfo.Role == constant.RoleExperts {
				contentExpert := message.SmsBooking
				contentExpert = replacer.Replace(contentExpert)
				common.SendSMS(listenerInfo.PhoneNumber, contentUser)
			}
			flagSendWs = true
			break
		case 15:
			flagSendWs = true
			break
		case 5:
			flagSendWs = true
			break
		default:
			break
		}
		if flagSendWs {
			msgNotifyUser := message.NotificationAppointmentUser
			replacerNotifyUser := strings.NewReplacer("{number}", minuteString, "{listener_code}", listenerInfo.EmployeeId)
			msgNotifyUser = replacerNotifyUser.Replace(msgNotifyUser)
			notifyUser := notificationModel.Notification{
				Id:           primitive.NewObjectID(),
				ReceiverId:   a.UserId,
				ReceiverRole: constant.RoleUser,
				Content:      msgNotifyUser,
				Type:         constant.NotifyAppointment,
				CreatedAt:    time.Now().UTC(),
				UpdatedAt:    time.Now().UTC(),
			}
			_ = notificationRepository.PublishInterfaceNotification().CreateNotification(notifyUser)
			ResponseWS(constant.WSEventNotify, msgNotifyUser, a.UserId)

			msgNotifyListener := message.NotificationAppointmentListener
			replacerNotifyListener := strings.NewReplacer("{number}", minuteString, "{user_name}", userInfo.Name)
			msgNotifyListener = replacerNotifyListener.Replace(msgNotifyListener)
			notifyListener := notificationModel.Notification{
				Id:           primitive.NewObjectID(),
				ReceiverId:   a.ListenerId,
				ReceiverRole: constant.RoleListener,
				Content:      msgNotifyListener,
				Type:         constant.NotifyAppointment,
				CreatedAt:    time.Now().UTC(),
				UpdatedAt:    time.Now().UTC(),
			}
			_ = notificationRepository.PublishInterfaceNotification().CreateNotification(notifyListener)
			ResponseWS(constant.WSEventNotify, msgNotifyListener, a.ListenerId)
		}
	}
}

func ResponseWS(event string, message string, fromUserId string) {
	hub := wsControllers.GetHub()
	payload := wsControllers.SocketEventStruct{}
	payload.EventName = event
	payload.EventPayload = map[string]interface{}{
		"message":     message,
	}
	wsControllers.EmitToSpecificClient(hub, payload, fromUserId)
	//
	//postBody, _ := json.Marshal(map[string]string{
	//	"event":        event,
	//	"message":      message,
	//	"from_user_id": fromUserId,
	//})
	//responseBody := bytes.NewBuffer(postBody)
	////Leverage Go's HTTP Post function to make request
	//resp, err := http.Post("http://127.0.0.1:2997/api/v1/master/send-event-ws", "application/json", responseBody)
	////Handle Error
	//if err != nil {
	//	log.Fatalf("An Error Occured %v", err)
	//}
	//defer resp.Body.Close()
	////Read the response body
	//_, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//sb := string(body)
	//log.Printf(sb)
}
