package handlers

import (
	"github.com/gorilla/mux"
	"zatrasz75/tz_song_libraries/configs"
	_ "zatrasz75/tz_song_libraries/docs"
	"zatrasz75/tz_song_libraries/internal/middleware"
	"zatrasz75/tz_song_libraries/internal/repository"
	"zatrasz75/tz_song_libraries/pkg/logger"
)

type Api struct {
	l    logger.LoggersInterface
	repo *repository.Store
	cfg  *configs.Config
}

func New(l logger.LoggersInterface, repo *repository.Store, cfg *configs.Config) *Api {
	return &Api{
		l:    l,
		repo: repo,
		cfg:  cfg,
	}
}

//	@title			Реализация онлайн библиотеки песен
//	@version		1.0

// @contact.name Михаил Токмачев
// @contact.url https://t.me/Zatrasz
// @contact.email zatrasz@ya.ru

// @host						localhost:8586
// @BasePath					/

// NewRouter -.
func NewRouter(cfg *configs.Config, a *Api) *mux.Router {
	r := mux.NewRouter()

	s := r.PathPrefix("/songs").Subrouter()
	s.Use(middleware.CORS(cfg.Server.CORSAllowedOrigins))

	registerSongsHandlers(s, a)

	return r
}
