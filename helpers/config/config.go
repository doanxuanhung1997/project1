package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	EnvData = "data"
)

//define envValue struct
type envValue struct {
	Database database
	Mail     mail
	Server   server
	Secret   secret
	OTP      otp
	Fee      fee
	Call     call
	Agora    agora
}

//define env struct
type env struct {
	Data envValue
}

//define config struct
type config struct {
	Env env
}

//define database struct
type database struct {
	Name     string
	Host     string
	Port     int
	User     string
	Password string
}

//define secret struct
type secret struct {
	Salt         string
	Buffer       string
	JwtSecretKey string
}

//define server struct
type server struct {
	Host               string
	Port               string
	ExpireToken        int
	ExpireRefreshToken int
	Type               string
}

//define mail struct
type mail struct {
	Admin    string
	Email    string
	Password string
	Smtp     string
	Port     int
	Id       string
}

//define otp struct
type otp struct {
	Username string
	Pwd      string
	From     string
	Endpoint string
}

//define fee conversation
type fee struct {
	ListenerCall int64
	ListenerTip  float64
	ExpertsCall  float64
	ExpertsTip   float64
}

//define call config
type call struct {
	AppointmentPrice       int64
	RingingTimeAppointment int64
	RingingTimeCallNow     int64
}

//define call config
type agora struct {
	AppId          string
	AppCertificate string
}

type Config config
type Server server
type Env env
type EnvValue envValue
type Database database
type Mail mail
type OTP otp
type Fee fee
type Secret secret
type Call call
type Agora agora

var configValue *config

//load environment from file
func Loads(filePath string) *config {
	var fileName string
	var yamlFile []byte
	var err error

	if fileName, err = filepath.Abs(filePath); err != nil {
		panic(err)
	}

	if yamlFile, err = ioutil.ReadFile(fileName); err != nil {
		panic(err)
	}
	configValue = &config{}
	if err = yaml.Unmarshal(yamlFile, configValue); err != nil {
		panic(err)
	}
	return configValue
}

//set environment by env name
func SetEnv(env string) {
	if env != EnvData {
		panic("Invalid env")
	}
	err := os.Setenv("env", env)
	if err != nil {
		println(err.Error())
	}
}

//get and check env name correct
func GetEnvValue() *envValue {
	if configValue == nil {
		panic("Must run Loads first")
	}
	//env := GetEnv()
	return &configValue.Env.Data
}

//get secret
func GetSecret() string {
	if configValue == nil {
		panic("Must run Loads first")
	}
	return configValue.Env.Data.Secret.JwtSecretKey
}
