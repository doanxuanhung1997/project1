package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math"
	conversationRepository "sandexcare_backend/api/conversation/repository"
	paymentRepository "sandexcare_backend/api/payment/repository"
	"sandexcare_backend/helpers/constant"
	"sandexcare_backend/helpers/middlewares"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func unRegisterAndCloseConnection(c *Client) {
	c.Hub.unregister <- c
	c.WebSocketConnection.Close()
}

func setSocketPayloadReadConfig(c *Client) {
	c.WebSocketConnection.SetReadLimit(maxMessageSize)
	c.WebSocketConnection.SetReadDeadline(time.Now().Add(pongWait))
	c.WebSocketConnection.SetPongHandler(func(string) error { c.WebSocketConnection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
}

func handleSocketPayloadEvents(client *Client, socketEventPayload SocketEventStruct) {
	var socketEventResponse SocketEventStruct
	switch socketEventPayload.EventName {
	case "join":
		log.Printf(client.Username + " Join Event triggered")
		//BroadcastSocketEventToAllClient(client.Hub, SocketEventStruct{
		//	EventName: socketEventPayload.EventName,
		//	//EventPayload: getAllReadyCallUsers(client.Hub),
		//	EventPayload: JoinDisconnectPayload{
		//		UserID: client.UserID,
		//		Users:  getAllConnectedUsers(client.Hub),
		//	},
		//})
		BroadcastSocketEventToAllClient(client.Hub, SocketEventStruct{
			EventName:    constant.WSEventReadyCall,
			EventPayload: GetAllListenersReadyCall(client.Hub),
		})

	case "disconnect":
		log.Printf(client.Username + " Disconnect Event triggered")
		//BroadcastSocketEventToAllClient(client.Hub, SocketEventStruct{
		//	EventName: socketEventPayload.EventName,
		//	//EventPayload: getAllReadyCallUsers(client.Hub),
		//	EventPayload: JoinDisconnectPayload{
		//		UserID: client.UserID,
		//		Users:  getAllConnectedUsers(client.Hub),
		//	},
		//})
		BroadcastSocketEventToAllClient(client.Hub, SocketEventStruct{
			EventName:    constant.WSEventReadyCall,
			EventPayload: GetAllListenersReadyCall(client.Hub),
		})

	case "destroy_call":
		log.Printf("Destroy call event triggered")
		receiverId := socketEventPayload.EventPayload.(map[string]interface{})["receiverId"].(string)
		socketEventResponse.EventName = "destroy_call"
		socketEventResponse.EventPayload = map[string]interface{}{
			"call_id": socketEventPayload.EventPayload.(map[string]interface{})["callId"].(string),
		}
		EmitToSpecificClient(client.Hub, socketEventResponse, receiverId)

	case "payment":
		log.Printf("Payment event triggered")
		receiverId := socketEventPayload.EventPayload.(map[string]interface{})["receiverId"].(string)
		socketEventResponse.EventName = "payment"
		socketEventResponse.EventPayload = map[string]interface{}{
			"method": socketEventPayload.EventPayload.(map[string]interface{})["method"].(string),
			"status": socketEventPayload.EventPayload.(map[string]interface{})["status"].(string),
			"extra_data": socketEventPayload.EventPayload.(map[string]interface{})["extraData"].(string),
			"total": socketEventPayload.EventPayload.(map[string]interface{})["total"].(float64),
		}
		EmitToSpecificClient(client.Hub, socketEventResponse, receiverId)
	}

}

func getUsernameByUserID(hub *Hub, userID string) string {
	var username string
	for client := range hub.clients {
		if client.UserID == userID {
			username = client.Username
		}
	}
	return username
}

func getAllConnectedUsers(hub *Hub) []UserStruct {
	var users []UserStruct
	for singleClient := range hub.clients {
		users = append(users, UserStruct{
			Username:  singleClient.Username,
			UserID:    singleClient.UserID,
			Role:      singleClient.Role,
			ReadyCall: singleClient.CallReady,
		})
	}
	return users
}

func GetAllListenersReadyCall(hub *Hub) []ListenersReadyCall {
	var listeners []ListenersReadyCall
	for singleClient := range hub.clients {
		if singleClient.CallReady && singleClient.Role == constant.RoleListener {
			listeners = append(listeners, ListenersReadyCall{
				PhoneNumber: singleClient.Username,
				ListenerId:  singleClient.UserID,
				Role:        singleClient.Role,
			})
		}
	}
	return SortListenersReadyCall(listeners)

}

// Sắp xếp theo: CV nào có cuộc gọi gần nhất thì sếp cuối cùng. Ngược lại chuyên viên chưa có cuộc gọi nào thì đầu tiên.
func SortListenersReadyCall(listeners []ListenersReadyCall) (data []ListenersReadyCall) {
	timeNow, _ := time.Parse(constant.DateTimeFormat, time.Now().Format(constant.DateTimeFormat))
	for l, _ := range listeners {
		orderInfo, errOrder := paymentRepository.PublishInterfacePayment().GetUpcomingAppointmentOfListener(listeners[l].ListenerId)
		if errOrder != nil {
			listeners[l].ActiveCallNow = true
		} else {
			minute := math.Round(orderInfo.CallDatetime.Sub(timeNow).Minutes())
			if minute > 45 {
				listeners[l].ActiveCallNow = true
			}
		}

		lastCall, errCall := conversationRepository.PublishInterfaceConversation().GetLastCallOfListener(listeners[l].ListenerId)
		if errCall != nil {
			listeners[l].LastCall = timeNow.AddDate(-100, 0, 0)
		} else {
			listeners[l].LastCall = lastCall.EndCall
		}
	}

	var n = len(listeners)
	for i := 0; i < n; i++ {
		var minIdx = i
		for j := i; j < n; j++ {
			if listeners[j].LastCall.Unix() < listeners[minIdx].LastCall.Unix() {
				minIdx = j
			}
		}
		listeners[i], listeners[minIdx] = listeners[minIdx], listeners[i]
	}
	return listeners
}

func (c *Client) readPump() {
	var socketEventPayload SocketEventStruct

	defer unRegisterAndCloseConnection(c)

	setSocketPayloadReadConfig(c)

	for {
		_, payload, err := c.WebSocketConnection.ReadMessage()

		decoder := json.NewDecoder(bytes.NewReader(payload))
		decoderErr := decoder.Decode(&socketEventPayload)

		if decoderErr != nil {
			log.Printf("error: %v", decoderErr)
			break
		}

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ===: %v", err)
			}
			break
		}

		handleSocketPayloadEvents(c, socketEventPayload)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.WebSocketConnection.Close()
	}()
	for {
		select {
		case payload, ok := <-c.Send:
			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(payload)
			finalPayload := reqBodyBytes.Bytes()

			c.WebSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.WebSocketConnection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.WebSocketConnection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(finalPayload)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				json.NewEncoder(reqBodyBytes).Encode(<-c.Send)
				w.Write(reqBodyBytes.Bytes())
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.WebSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.WebSocketConnection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// CreateNewSocketUser creates a new socket user
func CreateNewSocketUser(hub *Hub, connection *websocket.Conn, tokenInfo middlewares.TokenInfo) {
	client := &Client{
		Hub:                 hub,
		WebSocketConnection: connection,
		Send:                make(chan SocketEventStruct),
		Username:            tokenInfo.PhoneNumber,
		UserID:              tokenInfo.Id,
		Role:                tokenInfo.Role,
	}

	go client.writePump()
	go client.readPump()

	client.Hub.register <- client
}

// HandleUserRegisterEvent will handle the Join event for New socket users
func HandleUserRegisterEvent(hub *Hub, client *Client) {
	hub.clients[client] = true
	handleSocketPayloadEvents(client, SocketEventStruct{
		EventName:    "join",
		EventPayload: client.UserID,
	})
}

// HandleUserDisconnectEvent will handle the Disconnect event for socket users
func HandleUserDisconnectEvent(hub *Hub, client *Client) {
	_, ok := hub.clients[client]
	if ok {
		delete(hub.clients, client)
		close(client.Send)

		handleSocketPayloadEvents(client, SocketEventStruct{
			EventName:    "disconnect",
			EventPayload: client.UserID,
		})
	}
}

// EmitToSpecificClient will emit the socket event to specific socket user
func EmitToSpecificClient(hub *Hub, payload SocketEventStruct, userID string) {
	for client := range hub.clients {
		if client.UserID == userID {
			select {
			case client.Send <- payload:
			default:
				close(client.Send)
				delete(hub.clients, client)
			}
		}
	}
}

// BroadcastSocketEventToAllClient will emit the socket events to all socket users
func BroadcastSocketEventToAllClient(hub *Hub, payload SocketEventStruct) {
	for client := range hub.clients {
		if client.Role == constant.RoleUser {
			select {
			case client.Send <- payload:
			default:
				close(client.Send)
				delete(hub.clients, client)
			}
		}
	}
}

func HandleListenerOnCall(hub *Hub, listenerId string) {
	for client := range hub.clients {
		if client.UserID == listenerId {
			client.CallReady = true
		}
	}
}

func HandleListenerOffCall(hub *Hub, listenerId string) {
	for client := range hub.clients {
		if client.UserID == listenerId {
			client.CallReady = false
		}
	}
}
