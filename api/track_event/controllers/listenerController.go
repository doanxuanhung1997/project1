package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"sandexcare_backend/api/track_event/model"
	"sandexcare_backend/api/track_event/repository"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/message"
	"sandexcare_backend/helpers/middlewares"
	wsControllers "sandexcare_backend/server_websocket/controllers"
	"time"
)

// @Summary
// @Schemes
// @Description API Tracking Action Listener
// @Param input body model.InputTrackActionListener true "Input tracking action of listener"
// @Tags Track Event
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseSuccess
// @Failure 400 {object} model.ResponseError
// @Router  /track-event/listener/action [post]
func TrackingActionListener() gin.HandlerFunc {
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
		if tokenInfo.Role != constant.RoleListener {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorPermissionDenied,
			})
			return
		}
		/*Parameter c.GetRawData*/
		rawBody, _ := c.GetRawData()
		input := model.InputTrackActionListener{}
		err := json.Unmarshal(rawBody, &input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorConvertInput,
			})
			return
		}

		if input.Action != constant.CheckIn && input.Action != constant.CheckOut &&
			input.Action != constant.OnCall && input.Action != constant.OffCall {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": message.MessageErrorActionInvalid,
			})
			return
		}

		hub := wsControllers.GetHub()
		switch input.Action {
		case constant.OnCall:
			wsControllers.HandleListenerOnCall(hub, tokenInfo.Id)
			break
		case constant.CheckOut:
			// TODO Check if the listener has turned on ready call.
			// if action is check_out => auto first off_call
			CreateTrackAction(tokenInfo.Id, constant.OffCall)
			wsControllers.HandleListenerOffCall(hub, tokenInfo.Id)
		case constant.OffCall:
			wsControllers.HandleListenerOffCall(hub, tokenInfo.Id)
			break
		default:
			break
		}
		payload := wsControllers.SocketEventStruct{}
		payload.EventName = constant.WSEventReadyCall
		payload.EventPayload = wsControllers.GetAllListenersReadyCall(hub)

		wsControllers.BroadcastSocketEventToAllClient(hub, payload)
		CreateTrackAction(tokenInfo.Id, input.Action)

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": message.MessageSuccess,
		})
	}
}

func CreateTrackAction(listenerId string, action string) {
	// init model
	dataTrackAction := model.TrackActionListener{}
	dataTrackAction.Id = primitive.NewObjectID()
	dataTrackAction.ListenerId = listenerId
	dataTrackAction.Action = action
	dataTrackAction.CreatedAt = time.Now().UTC()
	dataTrackAction.UpdatedAt = time.Now().UTC()
	_ = repository.PublishInterfaceTrackAction().CreateTrackActionListener(dataTrackAction)
}

func CreateTrackCall(listenerId string, callId string, action string) {
	// init model
	dataTrackCall := model.TrackCall{}
	dataTrackCall.Id = primitive.NewObjectID()
	dataTrackCall.CallId = callId
	dataTrackCall.ListenerId = listenerId
	dataTrackCall.Action = action
	dataTrackCall.CreatedAt = time.Now().UTC()
	dataTrackCall.UpdatedAt = time.Now().UTC()
	_ = repository.PublishInterfaceTrackAction().CreateTrackCall(dataTrackCall)
}
