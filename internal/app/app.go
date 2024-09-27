package app

import (
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"zatrasz75/tz_song_libraries/configs"
	"zatrasz75/tz_song_libraries/internal/handlers"
	"zatrasz75/tz_song_libraries/internal/middleware"
	"zatrasz75/tz_song_libraries/internal/repository"
	"zatrasz75/tz_song_libraries/pkg/logger"
	"zatrasz75/tz_song_libraries/pkg/postgres"
	"zatrasz75/tz_song_libraries/pkg/server"
)

func Run(cfg *configs.Config, l logger.LoggersInterface) {
	pg, err := postgres.New(cfg.DataBase.ConnStr, l, postgres.OptionSet(cfg.DataBase.PoolMax, cfg.DataBase.ConnAttempts, cfg.DataBase.ConnTimeout))
	if err != nil {
		l.Fatal("ошибка запуска - postgres.New:", err)
	}
	defer pg.Close()

	err = pg.Migrate(l)
	if err != nil {
		l.Fatal("ошибка миграции", err)
	}

	repo := repository.New(pg)

	// Создание экземпляра Api
	api := handlers.New(l, repo, cfg)
	router := handlers.NewRouter(cfg, api)

	router.Use(middleware.LoggingResponse)
	router.Use(middleware.SetHeader)

	// Swagger UI
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs/"))))
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	srv := server.New(router, server.OptionSet(cfg.Server.AddrHost, cfg.Server.AddrPort, cfg.Server.ReadTimeout, cfg.Server.WriteTimeout, cfg.Server.IdleTimeout, cfg.Server.ShutdownTime))
	go func() {
		err = srv.Start()
		if err != nil {
			l.Error("Остановка сервера:", err)
		}
	}()

	l.Info("Запуск сервера на http://" + cfg.Server.AddrHost + ":" + cfg.Server.AddrPort)
	l.Info("Документация Swagger API: http://" + cfg.Server.AddrHost + ":" + cfg.Server.AddrPort + "/swagger/index.html")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("принят сигнал прерывания прерывание %s", s.String())
	case err = <-srv.Notify():
		l.Error("получена ошибка сигнала прерывания сервера", err)
	}

	err = srv.Shutdown()
	if err != nil {
		l.Error("не удалось завершить работу сервера", err)
	}
}
