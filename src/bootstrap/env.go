package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	AnonUserId             string `mapstructure:"MONGO_ANON_USER_ID"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	BasePath               string `mapstructure:"BASE_PATH"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	MqttHost               string `mapstructure:"MQTT_HOST"`
	MqttPort               int    `mapstructure:"MQTT_PORT"`
	MqttUser               string `mapstructure:"MQTT_USER"`
	MqttPass               string `mapstructure:"MQTT_PASSWORD"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	GmailSender            string `mapstructure:"GMAIL_SENDER"`
	GmailSenderAppPassword string `mapstructure:"GMAIL_SENDER_APP_PASSWORD"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile("/envs/.env")
	// to run locally
	// viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	// fmt.Println(env)
	return &env
}
