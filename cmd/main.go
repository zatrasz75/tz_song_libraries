package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"zatrasz75/tz_song_libraries/configs"
	"zatrasz75/tz_song_libraries/internal/app"
	"zatrasz75/tz_song_libraries/pkg/logger"
)

// Init Загружает значения из файла .env в систему.
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Файл .env не найден.", err)
	}
}

func main() {
	l := logger.NewLogger()

	// Получаем текущий рабочий каталог
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка при получении текущего рабочего каталога:", err)
		return
	}
	// Построение абсолютного пути к файлу configs.yml
	configPath := filepath.Join(cwd, "configs", "configs.yml")

	// Configuration
	cfg, err := configs.NewConfig(configPath)
	if err != nil {
		l.Fatal("ошибка при разборе конфигурационного файла", err)
	}

	app.Run(cfg, l)
}
