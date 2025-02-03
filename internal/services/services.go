package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"songs/config"
	"songs/internal/logger"
	"songs/internal/models"
)

func FetchSongDetail(group, songTitle string) (*models.SongDetail, error) {
	logger.Log.Info("Получение информации о песне с внешнего API")
	apiBase := config.Get("MUSIC_API_URL")
	if apiBase == "" {
		return nil, fmt.Errorf("MUSIC_API_URL не задан")
	}

	url := fmt.Sprintf("%s/info?group=%s&song=%s", apiBase, group, songTitle)
	logger.Log.Debugf("URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к внешнему API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("внешнее API вернуло статус %d", resp.StatusCode)
	}

	var detail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}
	logger.Log.Info("Успешно получены данные о песне с внешнего API")
	return &detail, nil
}

func SplitVerses(text string) []string {
	verses := strings.Split(text, "\n\n")
	var result []string
	for _, v := range verses {
		v = strings.TrimSpace(v)
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}
