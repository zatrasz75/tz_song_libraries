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
