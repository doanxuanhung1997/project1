package cronjob

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func InitCron()  {
	// init cron jobs here
	fmt.Println("init cron jobs")
	c := cron.New()
	c.AddFunc("@every 60s", UpcomingAppointmentReminders)
	c.Start()
}