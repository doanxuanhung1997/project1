package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sandexcare_backend/helpers/config"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func SendNotification(content string) {
	data := url.Values{}
	data.Set("text", content)
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot2082980755:AAHkEB4RsO2x-6YBVjBbOZtQXMe4_AJdAMg/sendMessage?chat_id=-1001749804629", strings.NewReader(data.Encode()))
	// Header - API get user information
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", os.Getenv("LINE_TOKEN"))
	if err != nil {
		fmt.Print(err.Error())
	}
	//
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
}

func GetIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ""
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			v := ip.String()
			return v
		}
	}
	return ""
}

func SendNotificationStarted() {
	env := config.GetEnvValue()

	log.Error(env.Server)
	if env.Server.Type == "build_dev" || env.Server.Type == "build_pro" {
		SendNotification("🆘🆘🆘 Server started at: " + time.Now().String() + "\n IP: " + GetIP())
	}
}

func SendHealthCheck(cpu, ram string) {
	env := config.GetEnvValue()

	log.Error(env.Server)
	if env.Server.Type == "build_dev" || env.Server.Type == "build_pro" {
		SendNotification("Server : " + GetIP() +
			"\n 🎄 HEAP: " + ram +
			"\n 🎄 STACK: " + cpu)
	}
}

func SendAPITime(t time.Duration, url, status string) {
	SendNotification("GỌI API: " + url + " THỜI GIAN API PHẢN HỒI: " + fmt.Sprint(t) + " STATUS:" + fmt.Sprint(status))
}

func SendNotifyFirebase(title, content, icon, to string) error {
	message := map[string]interface{}{
		"data": map[string]interface{}{
			"notification": map[string]interface{}{
				"title": title,
				"body":  content,
				"icon":  icon,
			},
		},
		"to": to,
	}
	start := time.Now()

	data := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(data)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(message)

	req, err := http.NewRequest(http.MethodPost, "https://fcm.googleapis.com/fcm/send", data)
	if err != nil {
		return errors.New("can not make a request by abnormal reason")
	}
	// Header - API get user information
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "key=AAAA8ZqDyjw:APA91bFsAnG1DodfAYhabyssmL-Pm5i3xoV4t0qPCiTKkl73qdTh4eP4ZpIY8nZW0uPsFpQGnw2za_062jMZtRJ_AzadMV7UwyzWYTIzamj6r7msLG-BvXYFx-kW4npkpBU10aqTq2RA")

	//
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	end := time.Since(start)
	if err != nil {
		log.Error(err)
		return errors.New("can not send request to API ")
	}
	defer func() {
		resp.Body.Close()
	}()
	SendAPITime(end, "FIREBASE NOTIFY", resp.Status)
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		// callAPI(token, wsToken, aType, url, body)
		return errors.New("forbidden API")
	default:
		return errors.New("abnormal error wwith api authen")
	}
}
