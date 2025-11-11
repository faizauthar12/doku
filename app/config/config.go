package config

import (
	"github.com/faizauthar12/doku/app/utils/helper"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Doku struct {
		ClientID  string
		SecretKey string
	}
}

var cfg Configuration

func GetEnvString(key string, dflt string) string {
	value := os.Getenv(key)
	if value == "" {
		return dflt
	}
	return value
}

func GetEnvSliceString(key string, dflt string) []string {
	values := os.Getenv(key)
	if values == "" {
		return strings.Split(dflt, ",")
	}
	return strings.Split(values, ",")
}

func GetEnvInt(key string, dflt int) int {
	value := os.Getenv(key)
	i, err := strconv.ParseInt(value, 10, 64)
	if value == "" && err != nil {
		return dflt
	}
	return int(i)
}

func GetEnvBool(key string, dflt bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return dflt
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return dflt
	}
	return b
}

func Get() Configuration {
	return cfg
}

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		logData := helper.WriteLog(err, http.StatusInternalServerError, helper.DefaultStatusText[http.StatusInternalServerError])
		log.Fatalf("Error loading .env file")
		log.Fatalf("Error: %v", logData)
	}

	// Doku Config
	cfg.Doku.ClientID = GetEnvString("DOKU_API_CLIENT_ID", "")
	cfg.Doku.SecretKey = GetEnvString("DOKU_API_SECRET_KEY", "")
}
