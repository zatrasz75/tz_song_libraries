package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	_ "zatrasz75/tz_song_libraries/docs"
	"zatrasz75/tz_song_libraries/internal/models"
)

func registerSongsHandlers(s *mux.Router, a *Api) {
	s.HandleFunc("", a.creatSongHandler).Methods(http.MethodPost)
}

// createUserHandler godoc
//
// @Summary Добавление новой песни
// @Tags		Songs
// @Description Принимает обязательные поля group , song .
// @Accept  json
// @Produce  json
// @Param   songs body models.Songs true "Данные структуры песни"
// @Success 200 {string} string "Успешно созданная запись < Ok >"
// @Failure 400 {string} string "Неверный формат запроса или не верно заполнены обязательные поля"
// @Failure 500 {string} string "Не получены детальные данные из API или ошибка при сохранении в бд"
// @Router /songs [post]
func (a *Api) creatSongHandler(w http.ResponseWriter, r *http.Request) {
	var newSong models.Songs
	if err := json.NewDecoder(r.Body).Decode(&newSong); err != nil {
		a.l.Error("не удалось проанализировать запрос JSON", err)
		http.Error(w, "не удалось проанализировать запрос JSON", http.StatusBadRequest)
		return
	}

	// Проверка обязательных полей
	if newSong.Group == "" || newSong.Song == "" {
		http.Error(w, "Поля Группа и Песня обязательны для заполнения", http.StatusBadRequest)
		return
	}

	// Вызов внешнего API для получения дополнительной информации о песне
	songDetail, err := _fetchSongDetail(newSong.Group, newSong.Song, a.cfg.External.Url)
	if err != nil {
		a.l.Error("Ошибка при получении сведений о песне", err)
		http.Error(w, fmt.Sprintf("Ошибка при получении сведений о песне: %v", err), http.StatusInternalServerError)
		return
	}
	newSong.Detail.ReleaseDate = songDetail.ReleaseDate
	newSong.Detail.Text = songDetail.Text
	newSong.Detail.Link = songDetail.Link

	// Сохранения песни в базе данных
	id, err := a.repo.CreatSong(newSong)
	if err != nil {
		a.l.Error("Ошибка при сохранении песни в базе данных", err)
		http.Error(w, fmt.Sprintf("Ошибка при сохранении песни: %v", err), http.StatusInternalServerError)
		return
	}
	idStr := strconv.Itoa(id)

	// Ответ клиенту
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Ok",
		"id":      idStr,
	})
}

func _fetchSongDetail(group, song, baseUrl string) (models.SongDetail, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/info?group=%s&song=%s", baseUrl, group, song)

	var httpClient *http.Client
	if strings.HasPrefix(baseUrl, "https") {
		// На случай https , отключаем проверку сертификата
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else {
		// Для HTTP используем стандартного клиента
		httpClient = http.DefaultClient
	}

	var songDetail models.SongDetail

	// Отправляем GET-запрос
	resp, err := httpClient.Get(url)
	if err != nil {
		return songDetail, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		return songDetail, err
	}

	return songDetail, nil
}
