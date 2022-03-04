package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	listenerModel "sandexcare_backend/api/listener/model"
	userModel "sandexcare_backend/api/user/model"
	"sandexcare_backend/db"
	"sandexcare_backend/helpers/config"
	"sandexcare_backend/helpers/constant"
	"strings"
	"time"
	"unicode/utf8"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var letterRunesNumber = []rune("0123456789")
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UTC().UnixNano()))

/*General Token Register Code Client*/
func GenerateTokenString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetDayOfWeek(date time.Time) string {
	day := "Thu Hai"
	weekday := int(date.Weekday())
	switch weekday {
	case 0:
		day = "Chu Nhat"
		break
	case 1:
		day = "Thu Hai"
		break
	case 2:
		day = "Thu Ba"
		break
	case 3:
		day = "Thu Tu"
		break
	case 4:
		day = "Thu Nam"
		break
	case 5:
		day = "Thu Sau"
		break
	case 6:
		day = "Thu Bay"
		break
	}
	return day
}

/*General Number*/
func GenerateNumber(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunesNumber[seededRand.Intn(len(letterRunesNumber))]
	}
	return string(b)
}

/*Send SMS*/
func SendSMS(phoneNumber string, content string) bool {
	env := config.GetEnvValue()
	postBody, _ := json.Marshal(map[string]string{
		"u":     env.OTP.Username,
		"pwd":   env.OTP.Pwd,
		"from":  env.OTP.From,
		"phone": phoneNumber,
		"sms":   content,
	})
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post(env.OTP.Endpoint, "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Printf(err.Error())
		return false
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(err.Error())
		return false
	} else {
		responseSendSMS := ResponseSendSMS{}
		json.Unmarshal([]byte(body), &responseSendSMS)
		if responseSendSMS.Error != 0 {
			log.Printf(responseSendSMS.Log)
		}
	}
	return true
}

//Response send sms
type ResponseSendSMS struct {
	Carrier   string `json:"id"`
	Error     int    `json:"error"`
	ErrorCode string `json:"error_code"`
	Msgid     string `json:"msgid"`
	Log       string `json:"log"`
}

//hash password from user
func HashPassword(pass string) string {
	env := config.GetEnvValue()
	mixPass := pass + env.Secret.Salt
	bytePass := []byte(mixPass)

	hash := sha256.New()
	hash.Write(bytePass)
	passwordHash := hex.EncodeToString(hash.Sum(nil))
	//return  passwordHash
	return env.Secret.Buffer + passwordHash
}

/*Check email exist return true if email existed in listener table or user table*/
func CheckEmailExist(value string, role int) bool {
	if role == constant.RoleUser {
		var user userModel.User
		err := db.Collection(userModel.CollectionUsers).FindOne(db.GetContext(), bson.M{"email": value}).Decode(&user)
		if err == nil {
			return true
		}
	} else {
		var listener listenerModel.Listener
		err := db.Collection(listenerModel.CollectionListeners).FindOne(db.GetContext(), bson.M{"email": value}).Decode(&listener)
		if err == nil {
			return true
		}
	}
	return false
}

/*Check phone number exist return true if phone number existed in listener table or user table*/
func CheckPhoneNumberExist(value string, role int) bool {
	if role == constant.RoleUser {
		var user userModel.User
		err := db.Collection(userModel.CollectionUsers).FindOne(db.GetContext(), bson.M{"phone_number": value}).Decode(&user)
		if err == nil {
			return true
		}
	} else {
		var listener listenerModel.Listener
		err := db.Collection(listenerModel.CollectionListeners).FindOne(db.GetContext(), bson.M{"phone_number": value}).Decode(&listener)
		if err == nil {
			return true
		}
	}
	return false
}

/*Check employee_id existed in listener table*/
func CheckEmployeeIdExist(value string) bool {
	var listener listenerModel.Listener
	err := db.Collection(listenerModel.CollectionListeners).FindOne(db.GetContext(), bson.M{"employee_id": value}).Decode(&listener)
	if err == nil {
		return true
	}

	return false
}

/*IsEmpty return True if string is empty*/
func IsEmpty(value string) bool {
	if utf8.RuneCountInString(value) == 0 || value == "" {
		return true
	}
	return false
}

/*CheckValidationEmail return true when string is email and match regex*/
func CheckValidationEmail(email string) bool {
	re := regexp.MustCompile(`(?:[a-z0-9!#$%&'*+/=?^_'{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_'{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`)
	if !re.MatchString(email) {
		return false
	}
	return true
}

/*CheckValidationPhoneNumber return true when string is phone and match regex*/
func CheckValidationPhoneNumber(value string) bool {
	re := regexp.MustCompile(`(0[3|5|7|8|9][0-9]{8})\b`)
	if re.MatchString(value) {
		return true
	}
	return false
}

/*Check string is format date*/
func CheckFormatDate(value string) bool {
	_, err := time.Parse(constant.DateFormat, value)
	if err != nil {
		return false
	}
	return true
}

/*CheckLength return true if string length input passed min,max input*/
func CheckLength(value string, minLength int, maxLength int) bool {
	if utf8.RuneCountInString(value) < minLength || utf8.RuneCountInString(value) > maxLength {
		return false
	}
	return true
}

/*Check string is format number*/
func CheckIsNumber(value string) bool {
	re := regexp.MustCompile(`^[0-9]*$`)
	if !re.MatchString(value) {
		return false
	}
	return true
}

/*Check contain Special in string*/
func CheckSpecialCharacters(value string) bool {
	re := regexp.MustCompile(`[!@#$%^&*(),._?'"` + `:{+}|<>/-]`)
	if re.MatchString(value) {
		return true
	}
	return false
}

/*Get fullname of listener*/
func GetFullNameOfListener(listener listenerModel.Listener) (name string) {
	name = listener.Name.FirstName + " " + listener.Name.LastName
	return name
}

/*Check Validation Password*/
func IsValidPasswordFormat(value string) bool {
	result := true
	regexps := []string{".{8,}", "[a-z]", "[A-Z]", "[0-9]", `[!@#$%&*()\-_+=\[\]{}|;:<>?/.,]`}
	for _, r := range regexps {
		t, _ := regexp.MatchString(r, value)
		if !t {
			result = false
			break
		}
	}
	return result
}

func InitCouponForUser(userId string, discount float64, numberDateExpire int, typeCoupon string) {
	couponName := constant.CouponNameCallNow
	switch typeCoupon {
	case constant.CouponBookingCV:
		couponName = constant.CouponNameBookingCV
		break
	case constant.CouponBookingCG:
		couponName = constant.CouponNameBookingCV
		break
	default:
		break
	}
	replacer := strings.NewReplacer("{discount}", fmt.Sprintf("%0.f", discount*100))
	couponName = replacer.Replace(couponName)
	couponUser := userModel.CouponsUser{}
	couponUser.UserId = userId
	couponUser.Type = typeCoupon
	couponUser.Status = constant.Active
	couponUser.Discount = discount
	couponUser.Name = couponName
	couponUser.ExpiresAt = time.Now().UTC().AddDate(0, 0, numberDateExpire)
	couponUser.CreatedAt = time.Now().UTC()
	couponUser.UpdatedAt = time.Now().UTC()
	_, _ = db.Collection(userModel.CollectionCouponsUser).InsertOne(db.GetContext(), couponUser)
}
