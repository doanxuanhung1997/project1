package main

import (
	"io"
	"os"
	routerAdmin "sandexcare_backend/api/admin/router"
	routerConversation "sandexcare_backend/api/conversation/router"
	routerFirebase "sandexcare_backend/api/firebase/router"
	routerListener "sandexcare_backend/api/listener/router"
	masterUser "sandexcare_backend/api/master/router"
	routerHealth "sandexcare_backend/api/monitor/router"
	routerNotification "sandexcare_backend/api/notification/router"
	routerPayment "sandexcare_backend/api/payment/router"
	routerSchedule "sandexcare_backend/api/schedule/router"
	routerTrackEvent "sandexcare_backend/api/track_event/router"
	routerUser "sandexcare_backend/api/user/router"
	cron "sandexcare_backend/cronjob"
	"sandexcare_backend/db"
	"sandexcare_backend/docs"
	"sandexcare_backend/helpers/cache"
	"sandexcare_backend/helpers/config"
	"sandexcare_backend/notification"
	wsSetup "sandexcare_backend/server_websocket/setup"
	"time"

	_ "sandexcare_backend/docs"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	f, err := os.OpenFile("sandexcare.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		// Cannot open log file. Logging to stderr
		return
	}
	log.SetOutput(io.MultiWriter(os.Stderr, f))
	log.Warn("************************************************************")
	log.Warn("**   					SANDEXCARE.COM 						**")
	log.Warn("************************************************************")
}

// @title Tài liệu đặc tả api.sandexcare.com
// @version 1.0
// @description  Có mặt trên tất cả nền tảng kỹ thuật số hàng đầu như Google Play, iOS, website. SandexCare là ứng dụng dịch vụ hỗ trợ khách hàng giải tỏa những khó khăn trong cuộc sống qua cuộc gọi 1-1 với chuyên viên lắng nghe của SandexCare, đem lại cảm hứng và động lực trong cuộc sống - công ty hàng đầu về sức khỏe tâm lý cộng đồng tại Việt Nam.
// @description SandexCare sẽ luôn bên bạn mọi lúc khi bạn cần nhất và hãy đến với chúng tôi để cùng nhau lan toả năng lượng tích cực cho cuộc sống để những vấn đề của chính bạn trở nên đơn giản hơn.
// @description
// @description Format date: "2022-01-04".
// @description Format time_slot: "00:00-06:00".
// @description Format booking_time: "01:00".
// @termsOfService http://sandexcare.com/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	//go Mem()
	// environment
	arg := "dev"
	cache.InitCache()
	config.Loads("./" + arg + ".env")
	config.SetEnv(config.EnvData)
	env := config.GetEnvValue()
	app := gin.Default()
	app.Use(ginsession.New())
	app.Use(CORSMiddleware())
	docs.SwaggerInfo.BasePath = "/api/v1"
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	routerFirebase.InitRouter(app)
	routerHealth.InitRouter(app)
	routerUser.InitRouter(app)
	routerAdmin.InitRouter(app)
	routerListener.InitRouter(app)
	routerSchedule.InitRouter(app)
	routerPayment.InitRouter(app)
	routerConversation.InitRouter(app)
	routerNotification.InitRouter(app)
	routerTrackEvent.InitRouter(app)
	masterUser.InitRouter(app)
	if err := db.InitDb(); err != nil {
		panic(err)
	}
	wsSetup.SetupServerWebSocket(app)
	cron.InitCron()
	notification.SendNotificationStarted()
	app.Run(env.Server.Host + ":" + env.Server.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, accept, origin, Cache-Control, X-Requested-With, authorization, origin, content-type, accept, token, X-Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Writer.Header().Set("Allow", "GET,POST,PUT,DELETE,OPTIONS")
		c.Writer.Header().Set("Content-Type", "application/json")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}

func PrintMemUsage() {
	// var m runtime.MemStats
	// runtime.ReadMemStats(&m)
	// // For info on each, see: https://golang.org/pkg/runtime/#MemStats
	// fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	// fmt.Printf("\tMallocs = %v\t", bToMb(m.Mallocs))
	// fmt.Printf("\tStackInUse = %v MiB", bToMb(m.StackInuse))
	// fmt.Printf("\tHeapInUse = %v\t", bToMb(m.HeapInuse))
	// fmt.Printf("\tHeapRealeased = %v\t", bToMb(m.HeapReleased))
	// fmt.Printf("\tFrees = %v\t", bToMb(m.Frees))
	// fmt.Printf("\tNumGC = %v\n", m.NumGC)
	// notification.SendHealthCheck(fmt.Sprint(m.HeapInuse), fmt.Sprint(bToMb(m.StackInuse)))
}

func Mem() {
	for {
		// Print our memory usage at each interval
		PrintMemUsage()
		time.Sleep(30 * time.Minute)
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
