package common

import (
	"crypto/sha256"
	"encoding/hex"
	"houze_ops_backend/configs"
	"houze_ops_backend/helpers/constants"
	"math/rand"
	"regexp"
	"time"
	"unicode/utf8"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var letterRunesNumber = []rune("0123456789")
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UTC().UnixNano()))

func GenerateTokenString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateNumber(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunesNumber[seededRand.Intn(len(letterRunesNumber))]
	}
	return string(b)
}

func HashPassword(pass string) string {
	env := configs.GetEnvConfig()
	mixPass := pass + env.JWTSecretKey
	bytePass := []byte(mixPass)

	hash := sha256.New()
	hash.Write(bytePass)
	passwordHash := hex.EncodeToString(hash.Sum(nil))
	//return  passwordHash
	return env.JWTBuffer + passwordHash
}

func CheckFormatDate(value string) bool {
	_, err := time.Parse(constants.DateFormat, value)
	if err != nil {
		return false
	}
	return true
}

func CheckLength(value string, minLength int, maxLength int) bool {
	if utf8.RuneCountInString(value) < minLength || utf8.RuneCountInString(value) > maxLength {
		return false
	}
	return true
}

func CheckIsNumber(value string) bool {
	re := regexp.MustCompile(`^[0-9]*$`)
	if !re.MatchString(value) {
		return false
	}
	return true
}

func CheckSpecialCharacters(value string) bool {
	re := regexp.MustCompile(`[!@#$%^&*(),._?'"` + `:{+}|<>/-]`)
	if re.MatchString(value) {
		return true
	}
	return false
}

func IsEmpty(value string) bool {
	if utf8.RuneCountInString(value) == 0 || value == "" {
		return true
	}
	return false
}

func GetDateTimeNow() time.Time {
	return time.Now().UTC()
}
