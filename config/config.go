package config

import (
	"os"

	"songs/internal/logger"

	"github.com/joho/godotenv"
)

var AppConfig map[string]string

func init() {
	if err := godotenv.Load(".env"); err != nil {
		logger.Log.Errorf("Ошибка загрузки .env файла: %v", err)
	}
	AppConfig = map[string]string{
		"PORT":          os.Getenv("PORT"),
		"DATABASE_URL":  os.Getenv("DATABASE_URL"),
		"MUSIC_API_URL": os.Getenv("MUSIC_API_URL"),
	}
	logger.Log.Debugf("Переменные окружения: %v", AppConfig)
}

func Get(key string) string {
	return AppConfig[key]
}
