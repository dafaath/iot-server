package configs

import (
	"bytes"
	"log"
	"os"
	"path"
	"strings"

	_ "embed"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		Env  string `json:"env"`
		Domain  string `json:"domain"`
	} `json:"server"`
	Database struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
	} `json:"database"`
	JWT struct {
		SecretKey string `json:"secretKey"`
	} `json:"jwt"`
	Mail struct {
		SMTPHost               string `json:"smtpHost"`
		SMTPPort               int    `json:"smtpPort"`
		SenderName             string `json:"senderName"`
		AuthenticationMail     string `json:"authenticationMail"`
		AuthenticationPassword string `json:"authenticationPassword"`
	} `json:"mail"`
	Account struct {
		AdminUsername string `json:"adminUsername"`
		AdminEmail    string `json:"adminEmail"`
		AdminPassword string `json:"adminPassword"`
		UserUsername  string `json:"userUsername"`
		UserEmail     string `json:"userEmail"`
		UserPassword  string `json:"userPassword"`
	} `json:"account"`
}

//go:embed config.json
var configFile []byte
var cfg Config

func init() {
	working_directory, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory, %s", err.Error())
	}

	configSettings := viper.New()

	env_path := path.Join(working_directory, ".env")
	err = godotenv.Load(env_path)
	if err != nil {
		log.Fatalf("Error getting reading env, %s", err.Error())
	}

	// Environment variables
	configSettings.AutomaticEnv()
	configSettings.SetEnvPrefix("APP")
	configSettings.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Configuration file
	configSettings.SetConfigType("json")

	err = configSettings.ReadConfig(bytes.NewBuffer(configFile))
	if err != nil {
		log.Fatalf("Error reading log, %s", err.Error())
	}

	err = configSettings.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Error unmarshal config, %s", err.Error())
	}
}

func GetConfig() *Config {
	return &cfg
}
