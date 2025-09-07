package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                string
	IsProd                 string
	DB_HOST                string
	DB_USER                string
	DB_PASS                string
	DB_PORT                string
	DB_DATABASE            string
	GOOGLE_GEMINI_API_KEY  string
	JWT_SECRET             []byte
	STRIPE_SECRET_KEY      string
	AWS_IMAGES_BUCKET_NAME string
	AWS_ACCESS_KEY_ID      string
	AWS_SECRET_ACCESS_KEY  string
	PKG_30_LAUNCH          string
	PKG_100_LAUNCH         string
	PKG_200_LAUNCH         string
	STRIPE_WEBHOOK_SECRET  string
	FRONTEND_BASE_URL      string
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
	C.AWS_IMAGES_BUCKET_NAME = getEnv("AWS_IMAGES_BUCKET_NAME", "")
	C.AWS_ACCESS_KEY_ID = getEnv("AWS_ACCESS_KEY_ID", "")
	C.AWS_SECRET_ACCESS_KEY = getEnv("AWS_SECRET_ACCESS_KEY", "")
	C.PKG_30_LAUNCH = getEnv("PKG_30_LAUNCH", "")
	C.PKG_100_LAUNCH = getEnv("PKG_100_LAUNCH", "")
	C.PKG_200_LAUNCH = getEnv("PKG_200_LAUNCH", "")
	C.STRIPE_WEBHOOK_SECRET = getEnv("STRIPE_WEBHOOK_SECRET", "")
	C.FRONTEND_BASE_URL = getEnv("FRONTEND_BASE_URL", "")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
