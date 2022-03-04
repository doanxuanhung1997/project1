package controllers

import (
	"github.com/gorilla/websocket"
	"time"
)

// UserStruct is used for sending users with socket id
type UserStruct struct {
	Username  string `json:"username"`
	UserID    string `json:"user_id"`
	ReadyCall bool   `json:"ready_call"`
	Role      int    `json:"role"`
}

// ListenersReadyCall
type ListenersReadyCall struct {
	PhoneNumber   string    `json:"phone_number"`
	ListenerId    string    `json:"listener_id"`
	Role          int       `json:"role"`
	ActiveCallNow bool      `json:"active_call_now"`
	LastCall      time.Time `json:"last_call"`
}

// SocketEventStruct struct of socket events
type SocketEventStruct struct {
	EventName    string      `json:"eventName"`
	EventPayload interface{} `json:"eventPayload"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub                 *Hub
	WebSocketConnection *websocket.Conn
	Send                chan SocketEventStruct
	Username            string
	UserID              string
	Role                int
	CallReady           bool
}

// JoinDisconnectPayload will have struct for payload of join disconnect
type JoinDisconnectPayload struct {
	Users  []UserStruct `json:"users"`
	UserID string       `json:"userID"`
}
