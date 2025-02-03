package main

import (
	"songs/config"
	"songs/database"
	"songs/internal/handlers"
	"songs/internal/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Log.Info("Старт приложениия")

	database.Init()

	router := gin.Default()
	router.Use(gin.Logger())

	router.GET("/songs", handlers.GetSongs)

	port := config.Get("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Log.Debugf("Сервер запущен на порту %s", port)
	if err := router.Run(":" + port); err != nil {
		logger.Log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	logger.Log.Info("Завершение приложения")
}
