package configs

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type EnvConfig struct {
	ServerHost        string `yaml:"server_host" envconfig:"SERVER_HOST"`
	ServerPort        string `yaml:"server_port" envconfig:"SERVER_PORT"`
	ServerExpireToken int    `yaml:"server_expire_token" envconfig:"SERVER_EXPIRE_TOKEN"`
	DBName            string `yaml:"db_name" envconfig:"DB_NAME"`
	DBHost            string `yaml:"db_host" envconfig:"DB_HOST"`
	DBPort            string `yaml:"db_port" envconfig:"DB_PORT"`
	DBUser            string `yaml:"db_user" envconfig:"DB_USER"`
	DBPassword        string `yaml:"db_password" envconfig:"DB_PASSWORD"`
	MailAdmin         string `yaml:"mail_admin" envconfig:"MAIL_ADMIN"`
	MailPassword      string `yaml:"mail_password" envconfig:"MAIL_PASSWORD"`
	MailSmtp          string `yaml:"mail_smtp" envconfig:"MAIL_SMTP"`
	MailPort          string `yaml:"mail_port" envconfig:"MAIL_PORT"`
	JWTSalt           string `yaml:"jwt_salt" envconfig:"JWT_SALT"`
	JWTBuffer         string `yaml:"jwt_buffer" envconfig:"JWT_BUFFER"`
	JWTSecretKey      string `yaml:"jwt-secret-key" envconfig:"JWT_SECRET_KEY"`
}

var (
	env EnvConfig
)

func InitEnvConfig() {
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func GetEnvConfig() EnvConfig {
	return env
}
