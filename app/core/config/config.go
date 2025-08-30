package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort               string
	IsProd                string
	DB_HOST               string
	DB_USER               string
	DB_PASS               string
	DB_PORT               string
	DB_DATABASE           string
	GOOGLE_GEMINI_API_KEY string
	JWT_SECRET            []byte
	STRIPE_SECRET_KEY     string
}

var C = &Config{}

func init() {

	if err := godotenv.Load(); err != nil {
		log.Println("no .env found")
	}

	C.AppPort = getEnv("APP_PORT", "8080")
	C.IsProd = getEnv("IS_PROD", "false")
	C.DB_HOST = getEnv("DB_HOST", "localhost")
	C.DB_USER = getEnv("DB_USER", "root")
	C.DB_PASS = getEnv("DB_PASS", "")
	C.DB_DATABASE = getEnv("DB_DATABASE", "timelinelove")
	C.DB_PORT = getEnv("DB_PORT", "")
	C.GOOGLE_GEMINI_API_KEY = getEnv("GOOGLE_GEMINI_API_KEY", "")
	C.JWT_SECRET = []byte(getEnv("JWT_SECRET", ""))
	C.STRIPE_SECRET_KEY = getEnv("STRIPE_SECRET_KEY", "")

}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
