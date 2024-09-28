package repository

import (
	"context"
	"fmt"
	"time"
	"zatrasz75/tz_song_libraries/internal/models"
	"zatrasz75/tz_song_libraries/pkg/postgres"
)

type Store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Store {
	return &Store{pg}
}

// CreatSong Добавление песен
func (s *Store) CreatSong(m models.Songs) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Начать транзакцию
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("не удалось запустить транзакцию: %w", err)
	}
	defer tx.Rollback(ctx)

	query := "INSERT INTO songs (s_group, song, release_date, text, link) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var id int
	err = tx.QueryRow(ctx, query, m.Group, m.Song, m.Detail.ReleaseDate, m.Detail.Text, m.Detail.Link).Scan(&id)
	if err != nil {
		return 0, err
	}

	// Фиксация транзакции
	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("не удалось зафиксировать транзакцию: %w", err)
	}

	return id, nil
}

func (s *Store) GetLibraryData(filter string, offset, limit int) ([]models.Songs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, s_group, song, release_date, text, link FROM songs"
	var args []interface{}

	// Добавляем фильтрацию, если параметр filter задан
	if filter != "" {
		query += " WHERE s_group LIKE $1 OR song LIKE $1 OR release_date LIKE $1 OR text LIKE $1 OR link LIKE $1"
		args = append(args, "%"+filter+"%")
	}

	// Ограничиваем результаты по лимиту и смещению
	query += " ORDER BY id LIMIT $2 OFFSET $3"
	args = append(args, limit, (offset-1)*limit)

	rows, err := s.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось запросить список песен: %w", err)
	}
	defer rows.Close()

	var songs []models.Songs
	for rows.Next() {
		var song models.Songs
		err = rows.Scan(&song.ID, &song.Group, &song.Song, &song.Detail.ReleaseDate, &song.Detail.Text, &song.Detail.Link)
		if err != nil {
			return nil, fmt.Errorf("не удалось отсканировать строку: %w", err)
		}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе строк: %w", err)
	}

	return songs, nil
}
