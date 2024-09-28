package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
	"strconv"
	"strings"
	_ "zatrasz75/tz_song_libraries/docs"
	"zatrasz75/tz_song_libraries/internal/models"
)

func registerSongsHandlers(s *mux.Router, a *Api) {
	s.HandleFunc("", a.creatSongHandler).Methods(http.MethodPost)
	s.HandleFunc("", a.getLibraryDataHandler).Methods(http.MethodGet)
	s.HandleFunc("", a.deleteSongHandler).Methods(http.MethodDelete)
	s.HandleFunc("/lyrics", a.getSongLyricsHandler).Methods(http.MethodGet)

	// Swagger UI
	s.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs/"))))
	s.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}

// creatSongHandler godoc
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

// getLibraryDataHandler godoc
//
// @Summary Получение данных с фильтрацией по всем полям и пагинацией
// @Tags		Songs
// @Description Принимает поля group , song , releaseDate , text , link .
// @Accept  json
// @Produce  json
// @Param filter query string false "Фильтр по названию"
// @Param limit query integer false "Количество записей на странице для пагинации" default(10)
// @Param offset query integer false "Номер страницы" default(1)
// @Success 200 {array} models.Songs "Список песен"
// @Failure 400 {string} string "Недопустимый параметр смещения"
// @Failure 500 {string} string "Ошибка при получение списка песен или Ошибка при обработке запроса"
// @Router /songs [get]
func (a *Api) getLibraryDataHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	filter := queryParams.Get("filter")
	limit := queryParams.Get("limit")
	offset := queryParams.Get("offset")

	var lmt, off int
	var err error
	if limit != "" {
		lmt, err = strconv.Atoi(limit)
		if err != nil || lmt <= 0 {
			http.Error(w, "Недопустимый параметр смещения", http.StatusBadRequest)
			return
		}
	} else {
		lmt = 10
	}

	if offset != "" {
		off, err = strconv.Atoi(offset)
		if err != nil || off < 0 {
			http.Error(w, "Недопустимый параметр смещения", http.StatusBadRequest)
			return
		}
	} else {
		off = 1
	}

	list, err := a.repo.GetLibraryData(filter, off, lmt)
	if err != nil {
		a.l.Error("Ошибка при получение списка песен", err)
		http.Error(w, fmt.Sprintf("Ошибка при получение списка песен: %v", err), http.StatusInternalServerError)
		return
	}

	listJSON, err := json.Marshal(list)
	if err != nil {
		a.l.Error("Ошибка при преобразовании данных в JSON", err)
		http.Error(w, "Ошибка при обработке запроса", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(listJSON)
	if err != nil {
		http.Error(w, "ошибка при отправке данных", http.StatusInternalServerError)
		a.l.Error("ошибка при отправке данных: ", err)
		return
	}
}

// getSongLyricsHandler godoc
//
// @Summary Получение текста песни с пагинацией по куплетам
// @Tags		Songs
// @Description Принимает поля songId , limit , offset .
// @Accept  json
// @Produce  json
// @Param songId query integer false "ID записи"
// @Param limit query integer false "Количество записей на странице для пагинации" default(10)
// @Param offset query integer false "Номер страницы" default(1)
// @Success 200 {array} models.LyricResponse "Список куплетов"
// @Failure 400 {string} string "Недопустимый параметр смещения или Недопустимый параметр смещения"
// @Failure 500 {string} string "Ошибка при получение списка песен или Ошибка при обработке запроса"
// @Router /songs/lyrics [get]
func (a *Api) getSongLyricsHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	songID := queryParams.Get("songId")
	limit := queryParams.Get("limit")
	offset := queryParams.Get("offset")

	var lmt, off int
	var err error
	if limit != "" {
		lmt, err = strconv.Atoi(limit)
		if err != nil || lmt <= 0 {
			http.Error(w, "Недопустимый параметр смещения", http.StatusBadRequest)
			return
		}
	} else {
		lmt = 10
	}

	if offset != "" {
		off, err = strconv.Atoi(offset)
		if err != nil || off < 0 {
			http.Error(w, "Недопустимый параметр смещения", http.StatusBadRequest)
			return
		}
	} else {
		off = 1
	}

	// Предполагаем, что у нас есть метод GetSongLyrics в репозитории
	lyrics, nextID, err := a.repo.GetSongLyrics(songID, off, lmt)
	if err != nil {
		a.l.Error("Ошибка при получении текста песни", err)
		http.Error(w, "Ошибка при получении текста песни", http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	response := models.LyricResponse{
		Chunks:     lyrics,
		NextPageID: int(nextID),
	}

	// Преобразуем response в JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		a.l.Error("Ошибка при формировании JSON ответа", err)
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		a.l.Error("Ошибка при записи данных в ответ", err)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
		return
	}
}

// deleteSongHandler godoc
//
// @Summary Удаление песни из библиотеки по ID
// @Tags		Songs
// @Description Принимает поля songId .
// @Accept  json
// @Produce  json
// @Param songId query integer false "ID записи"
// @Success 200 {string} string "Данные песни успешно удалены"
// @Failure 400 {string} string "Отсутствует идентификатор в запросе"
// @Failure 500 {string} string "не удалось преобразовать строку в число или Ошибка при удалении данных"
// @Router /songs [delete]
// ]
func (a *Api) deleteSongHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	songID := queryParams.Get("songId")
	if songID == "" {
		a.l.Debug("Отсутствует идентификатор в запросе")
		http.Error(w, "Отсутствует идентификатор в запросе", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(songID)
	if err != nil {
		a.l.Error("не удалось преобразовать строку в число", err)
		http.Error(w, "не удалось преобразовать строку в число", http.StatusInternalServerError)
		return
	}

	err = a.repo.DeleteSongById(id)
	if err != nil {
		a.l.Error("Ошибка при удалении данных", err)
		http.Error(w, "Ошибка при удалении данных", http.StatusInternalServerError)
		return
	}
	a.l.Info("Данные песни c id %d успешно удалены", id)

	// Устанавливаем правильный Content-Type для HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Данные песни успешно удалены"))
	if err != nil {
		http.Error(w, "ошибка при отправке данных", http.StatusInternalServerError)
		a.l.Error("ошибка при отправке данных: ", err)
		return
	}
}
