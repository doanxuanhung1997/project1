package setup

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sandexcare_backend/helpers/middlewares"
	"sandexcare_backend/server_websocket/controllers"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Resolve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SetupServerWebSocket(app *gin.Engine) {
	//init Hub store websocket
	controllers.NewHub()
	hub := controllers.GetHub()
	go hub.Run()

	app.LoadHTMLGlob("server_websocket/views/index.html")
	app.Static("server_websocket/public", "./server_websocket/public")
	app.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	app.GET("/ws/:token", func(c *gin.Context) {
		token := c.Param("token")
		tokenInfo, errToken := middlewares.GetInfoByToken(token)
		if errToken != nil {
			return
		}
		connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
		}
		controllers.CreateNewSocketUser(hub, connection, tokenInfo)
	})
}
