package agora

import (
	"fmt"
	rtcTokenBuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/RtcTokenBuilder"
	"sandexcare_backend/helpers/config"
	"time"
)

// Use RtcTokenBuilder to Generate an RTC token.
func GenerateRtcToken(uidStr string, channel string) (rtcToken string) {

	//appID := "<Your App ID>"
	//appCertificate := "<Your App Certificate>"
	//// Number of seconds after which the token expires.
	//// For demonstration purposes the expiry time is set to 40 seconds. This shows you the automatic token renew actions of the client.
	//expireTimeInSeconds := uint32(40)
	//// Get current timestamp.
	//currentTimestamp := uint32(time.Now().UTC().Unix())
	//// Timestamp when the token expires.
	//expireTimestamp := currentTimestamp + expireTimeInSeconds
	env := config.GetEnvValue()

	appID := env.Agora.AppId
	appCertificate := env.Agora.AppCertificate
	expireTimeInSeconds := uint32(1800)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	rtcToken, err := rtcTokenBuilder.BuildTokenWithUserAccount(appID, appCertificate, channel, uidStr, rtcTokenBuilder.RoleAttendee, expireTimestamp)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Token with userAccount: %s\n", rtcToken)
	}
	return rtcToken
}
