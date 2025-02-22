package main

import (
	"songs/config"
	"songs/database"
	"songs/internal/handlers"
	"songs/internal/logger"

	_ "songs/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	logger.Log.Info("Старт приложениия")

	database.Init()

	router := gin.Default()
	router.Use(gin.Logger())

	router.GET("/songs", handlers.GetSongs)
	router.GET("/songs/:id", handlers.GetSong)
	router.GET("/songs/:id/text", handlers.GetSongText)
	router.POST("/songs", handlers.AddSong)
	router.PATCH("/songs/:id", handlers.PatchSong)
	router.DELETE("/songs/:id", handlers.DeleteSong)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
