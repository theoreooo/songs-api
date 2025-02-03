// internal/handlers/song_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"songs/database"
	"songs/internal/logger"
	"songs/internal/models"
	"songs/internal/services"

	"github.com/gin-gonic/gin"
)

// getSongs — получение списка песен с фильтрацией по полям и пагинацией
// @Summary Получение песен с фильтрацией и пагинацией
// @Description Получает список песен с фильтрацией по полям (например, group и song) и поддержкой пагинации
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Группа"
// @Param song query string false "Название песни"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param pageSize query int false "Размер страницы (по умолчанию 10)"
// @Success 200 {array} models.Song
// @Failure 500 {object} models.ErrorResponse
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	logger.Log.Info("Получение списка песен")
	var songs []models.Song
	query := database.DB.Model(&models.Song{})

	group := c.Query("group_name")
	if group != "" {
		query = query.Where("group ILIKE ?", "%"+group+"%")
		logger.Log.Debugf("Фильтрация по группе: %s", group)
	}

	songTitle := c.Query("song")
	if songTitle != "" {
		query = query.Where("song ILIKE ?", "%"+songTitle+"%")
		logger.Log.Debugf("Фильтрация по названию песни: %s", songTitle)
	}
	logger.Log.Debugf("Фильтры - group: %s, song: %s", group, songTitle)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize
	logger.Log.Debugf("Пагинация - страница: %d, размер: %d, offset: %d", page, pageSize, offset)

	if err := query.Limit(pageSize).Offset(offset).Find(&songs).Error; err != nil {
		logger.Log.Errorf("Ошибка при получении песен: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Log.Info("Песни успешно получены")
	c.JSON(http.StatusOK, songs)
}

// getSongText — получение текста песни с пагинацией по куплетам
// @Summary Получение текста песни с пагинацией по куплетам
// @Description Разбивает текст песни на куплеты и возвращает запрошенную страницу
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param pageSize query int false "Размер страницы (по умолчанию 5)"
// @Success 200 {object}  models.MessageResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /songs/{id}/text [get]
func GetSongText(c *gin.Context) {
	id := c.Param("id")
	logger.Log.Infof("Получение текста песни id: %s", id)
	var song models.Song
	if err := database.DB.First(&song, id).Error; err != nil {
		logger.Log.Errorf("Ошибка при получении текста песни: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
		return
	}

	verses := services.SplitVerses(song.Text)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(verses) {
		c.JSON(http.StatusOK, gin.H{"verses": []string{}})
		return
	}
	if end > len(verses) {
		end = len(verses)
	}
	logger.Log.Info("Успешно получен текст песни")
	c.JSON(http.StatusOK, gin.H{"verses": verses[start:end]})
}

// deleteSong — удаление песни по ID
// @Summary Удаление песни
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {object} models.MessageResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	id := c.Param("id")
	logger.Log.Infof("Удаление песни id: %s", id)
	if err := database.DB.Delete(&models.Song{}, id).Error; err != nil {
		logger.Log.Errorf("Песня не найдена: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
		return
	}
	logger.Log.Info("Песня успешно удалена в БД")
	c.JSON(http.StatusOK, gin.H{"message": "Песня удалена"})
}

// updateSong — изменение данных песни
// @Summary Изменение данных песни
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body models.Song true "Данные песни"
// @Success 200 {object} models.Song
// @Failure 400 {object} models.ErrorResponse
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	id := c.Param("id")
	logger.Log.Infof("Обновление песни id: %s", id)
	var song models.Song
	if err := database.DB.First(&song, id).Error; err != nil {
		logger.Log.Errorf("Песня не найдена: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
		return
	}
	logger.Log.Debugf("Песня: %v", song)

	var input models.Song
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Log.Errorf("Ошибка при биндинге JSON для новой песни: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Log.Debugf("input: %v", input)

	song.GroupName = input.GroupName
	song.Song = input.Song
	if input.ReleaseDate != "" {
		song.ReleaseDate = input.ReleaseDate
	}
	if input.Text != "" {
		song.Text = input.Text
	}
	if input.Link != "" {
		song.Link = input.Link
	}

	if err := database.DB.Save(&song).Error; err != nil {
		logger.Log.Errorf("Ошибка сохранения обновлений в БД: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Log.Info("Песня успешно обновлена в БД")
	c.JSON(http.StatusOK, song)
}

// addSong — добавление новой песни. При добавлении делается запрос во внешнее API для обогащения данных
// @Summary Добавление новой песни
// @Tags songs
// @Accept json
// @Produce json
// @Param song body map[string]string true "Данные песни (обязательные поля: group и song)"
// @Success 201 {object} models.Song
// @Failure 400 {object} models.ErrorResponse
// @Router /songs [post]
func AddSong(c *gin.Context) {
	logger.Log.Info("Добавление песни")
	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Log.Debugf("Инпат: %v", input)

	group, okGroup := input["group"]
	songTitle, okSong := input["song"]
	if !okGroup || !okSong {
		logger.Log.Errorf("Пустые поля group: %s и song: %s", group, songTitle)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поля group и song обязательны"})
		return
	}

	detail, err := services.FetchSongDetail(group, songTitle)
	if err != nil {
		logger.Log.Errorf("Ошибка получения данных с внешнего API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить информацию о песне"})
		return
	}
	logger.Log.Debugf("Данные полученные о песне с внешнего API: %v", detail)

	newSong := models.Song{
		GroupName:   group,
		Song:        songTitle,
		ReleaseDate: detail.ReleaseDate,
		Text:        detail.Text,
		Link:        detail.Link,
	}
	logger.Log.Debugf("Песня: %v", newSong)

	if err := database.DB.Create(&newSong).Error; err != nil {
		logger.Log.Errorf("Ошибка создания записи в БД о песне: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Log.Info("Песня успешно сохранена в БД")
	c.JSON(http.StatusCreated, newSong)
}
