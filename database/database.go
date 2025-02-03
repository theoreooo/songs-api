package database

import (
	"songs/config"
	"songs/internal/logger"
	"songs/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	logger.Log.Info("Инициализация БД")
	dsn := config.Get("DATABASE_URL")
	if dsn == "" {
		logger.Log.Fatal("DATABASE_URL не задан")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	if err := db.AutoMigrate(&models.Song{}); err != nil {
		logger.Log.Fatalf("Ошибка миграции: %v", err)
	}
	DB = db
	logger.Log.Info("Успешное подключение к БД")
}
