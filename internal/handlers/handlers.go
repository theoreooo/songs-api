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

// GetSongs godoc
// @Summary Получение списка песен с фильтрацией
// @Description Возвращает список песен. Можно фильтровать по названию песни, группе, дате релиза и другим полям.
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Название группы для фильтрации (регистр не важен)"
// @Param song query string false "Название песни для фильтрации (регистр не важен)"
// @Param releaseDate query string false "Дата релиза (в формате YYYY-MM-DD) для фильтрации"
// @Param text query string false "Фрагмент текста песни для поиска"
// @Param link query string false "Фрагмент URL ссылки для поиска"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param pageSize query int false "Размер страницы (по умолчанию 10)"
// @Success 200 {array} models.Song
// @Failure 500 {object} models.ErrorResponse
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	logger.Log.Info("Получение списка песен")
	var songs []models.Song
	query := database.DB.Preload("Artist").Model(&models.Song{})

	if group := c.Query("group"); group != "" {
		query = query.Joins("JOIN artists ON artists.id = songs.artist_id").
			Where("artists.name ILIKE ?", "%"+group+"%")
		logger.Log.Debugf("Фильтрация по группе: %s", group)
	}

	if songTitle := c.Query("song"); songTitle != "" {
		query = query.Where("songs.song ILIKE ?", "%"+songTitle+"%")
		logger.Log.Debugf("Фильтрация по названию песни: %s", songTitle)
	}

	if releaseDate := c.Query("releaseDate"); releaseDate != "" {
		query = query.Where("DATE(songs.release_date) = ?", releaseDate)
		logger.Log.Debugf("Фильтрация по дате релиза: %s", releaseDate)
	}

	if text := c.Query("text"); text != "" {
		query = query.Where("songs.text ILIKE ?", "%"+text+"%")
		logger.Log.Debugf("Фильтрация по тексту: %s", text)
	}

	if link := c.Query("link"); link != "" {
		query = query.Where("songs.link ILIKE ?", "%"+link+"%")
		logger.Log.Debugf("Фильтрация по ссылке: %s", link)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize
	logger.Log.Debugf("Пагинация - страница: %d, размер: %d, offset: %d", page, pageSize, offset)

	if err := query.Limit(pageSize).Offset(offset).Find(&songs).Error; err != nil {
		logger.Log.Errorf("Ошибка при получении песен: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	logger.Log.Info("Песни успешно получены")
	c.JSON(http.StatusOK, songs)
}

// GetSong godoc
// @Summary Получение детальной информации о песне
// @Description Возвращает информацию о песне по указанному ID, включая данные артиста.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {object} models.Song "Данные песни, включая артиста"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена"
// @Router /songs/{id} [get]
func GetSong(c *gin.Context) {
	id := c.Param("id")
	logger.Log.Infof("Получение песни id: %s", id)
	var song models.Song

	if err := database.DB.Preload("Artist").First(&song, id).Error; err != nil {
		logger.Log.Errorf("Песня не найдена: %v", err)
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Песня не найдена"})
		return
	}
	logger.Log.Infof("Песня успешно получена: %+v", song)
	c.JSON(http.StatusOK, song)
}

// GetSongText godoc
// @Summary Получение текста песни с пагинацией по куплетам
// @Description Разбивает текст песни на куплеты и возвращает запрошенную страницу.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param pageSize query int false "Размер страницы (по умолчанию 5)"
// @Success 200 {object} models.MessageResponse "Ответ содержит массив строк куплетов"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена"
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

// DeleteSong godoc
// @Summary Удаление песни
// @Description Удаляет песню по указанному ID.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {object} models.MessageResponse "Песня успешно удалена"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена"
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

// PatchSong godoc
// @Summary Частичное обновление данных песни
// @Description Обновляет указанные поля песни по ID. Если поле не передано, оно не изменяется.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body models.SongUpdate true "Данные для обновления песни"
// @Success 200 {object} models.Song "Обновлённые данные песни"
// @Failure 400 {object} models.ErrorResponse "Ошибка в запросе или данные невалидны"
// @Router /songs/{id} [patch]
func PatchSong(c *gin.Context) {
	id := c.Param("id")
	logger.Log.Infof("Частичное обновление песни id: %s", id)

	var song models.Song
	if err := database.DB.First(&song, id).Error; err != nil {
		logger.Log.Errorf("Песня не найдена: %v", err)
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Песня не найдена"})
		return
	}
	logger.Log.Debugf("Исходная песня: %+v", song)

	var input models.SongUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Log.Errorf("Ошибка при биндинге JSON для обновления песни: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	logger.Log.Debugf("Полученные данные для обновления: %+v", input)

	if input.GroupName != nil {
		var artist models.Artist
		if err := database.DB.First(&artist, song.ArtistID).Error; err != nil {
			logger.Log.Errorf("Артист не найден: %v", err)
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Артист не найден"})
			return
		}
		artist.Name = *input.GroupName
		if err := database.DB.Save(&artist).Error; err != nil {
			logger.Log.Errorf("Ошибка обновления артиста: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
	}
	if input.Song != nil {
		song.Song = *input.Song
	}
	if input.ReleaseDate != nil {
		song.ReleaseDate = *input.ReleaseDate
	}
	if input.Text != nil {
		song.Text = *input.Text
	}
	if input.Link != nil {
		song.Link = *input.Link
	}

	if err := database.DB.Save(&song).Error; err != nil {
		logger.Log.Errorf("Ошибка сохранения обновлений в БД: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	logger.Log.Info("Песня успешно обновлена в БД")
	c.JSON(http.StatusOK, song)
}

// AddSong godoc
// @Summary Добавление новой песни
// @Description Добавляет новую песню, обогащая данные через внешний API. Если артист с указанным именем не существует, он создается.
// @Tags songs
// @Accept json
// @Produce json
// @Param song body map[string]string true "Данные песни (обязательные поля: group и song)"
// @Success 201 {object} models.Song "Созданная песня с данными из внешнего API"
// @Failure 400 {object} models.ErrorResponse "Ошибка валидации входных данных"
// @Failure 500 {object} models.ErrorResponse "Ошибка при получении данных с внешнего API или сохранении в БД"
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

	var artist models.Artist
	if err := database.DB.Where("name ILIKE ?", group).First(&artist).Error; err != nil {
		artist = models.Artist{
			Name: group,
		}
		if err := database.DB.Create(&artist).Error; err != nil {
			logger.Log.Errorf("Ошибка создания артиста: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		logger.Log.Infof("Создан новый артист: %v", artist)
	} else {
		logger.Log.Infof("Найден существующий артист: %v", artist)
	}

	newSong := models.Song{
		ArtistID:    artist.ID,
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
