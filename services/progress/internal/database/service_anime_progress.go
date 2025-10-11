package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type AnimeProgress struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Episode   int       `json:"episode"`
	UserId    uuid.UUID `json:"user_id"`
	AnimeId   uuid.UUID `json:"anime_id"`
}

type ServiceAnimeProgress interface {
	CreateAnimeProgress(reqProgress *AnimeProgress) (*AnimeProgress, error)
	GetAnimeProgress(reqProgress *AnimeProgress) (*AnimeProgress, error)
	UpdateAnimeProgress(reqProgress *AnimeProgress) (*AnimeProgress, error)
	DeleteAnimeProgress(reqProgress *AnimeProgress) error
}

type serviceAnimeProgress struct {
	db ServiceDb
}

var serviceAnimeProgressInstance *serviceAnimeProgress

func NewServiceProgress(db ServiceDb) ServiceAnimeProgress {
	if serviceAnimeProgressInstance != nil {
		return serviceAnimeProgressInstance
	}
	newServiceProgress := &serviceAnimeProgress{
		db: db,
	}
	serviceAnimeProgressInstance = newServiceProgress
	return serviceAnimeProgressInstance
}

func (s *serviceAnimeProgress) CreateAnimeProgress(reqProgress *AnimeProgress) (*AnimeProgress, error) {
	tx, err := s.db.Conn().Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if err.Error() == "sql: transaction has already been committed or rolled back" {
				return
			}
			log.Printf("error: %v", err)
		}
	}()

	result, err := InsertAnimeProgress(tx, reqProgress)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *serviceAnimeProgress) GetAnimeProgress(reqProgress *AnimeProgress) (*AnimeProgress, error) {
	tx, err := s.db.Conn().Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if err.Error() == "sql: transaction has already been committed or rolled back" {
				return
			}
			log.Printf("error: %v", err)
		}
	}()

	result, err := SelectAnimeProgress(tx, reqProgress)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *serviceAnimeProgress) UpdateAnimeProgress(reqProgress *AnimeProgress) (*AnimeProgress, error) {
	tx, err := s.db.Conn().Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if err.Error() == "sql: transaction has already been committed or rolled back" {
				return
			}
			log.Printf("error: %v", err)
		}
	}()

	result, err := UpdateAnimeProgress(tx, reqProgress)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *serviceAnimeProgress) DeleteAnimeProgress(reqProgress *AnimeProgress) error {
	tx, err := s.db.Conn().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if err.Error() == "sql: transaction has already been committed or rolled back" {
				return
			}
			log.Printf("error: %v", err)
		}
	}()

	if err := DeleteAnimeProgress(tx, reqProgress); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func InsertAnimeProgress(tx *sql.Tx, reqProgress *AnimeProgress) (*AnimeProgress, error) {
	query := `
	INSERT INTO anime_progress (id, episode, user_id, anime_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at, episode, user_id, anime_id
	`

	result := &AnimeProgress{}

	if err := tx.QueryRow(
		query,
		uuid.New(),
		reqProgress.Episode,
		reqProgress.UserId,
		reqProgress.AnimeId,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Episode,
		&result.UserId,
		&result.AnimeId,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func SelectAnimeProgress(tx *sql.Tx, reqProgress *AnimeProgress) (*AnimeProgress, error) {
	query := `
	SELECT id, created_at, updated_at, episode, user_id, anime_id
	FROM anime_progress
	WHERE id = $1
	`

	result := &AnimeProgress{}

	if err := tx.QueryRow(
		query,
		reqProgress.Id,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Episode,
		&result.UserId,
		&result.AnimeId,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateAnimeProgress(tx *sql.Tx, reqProgress *AnimeProgress) (*AnimeProgress, error) {
	query := `
	UPDATE anime_progress
	SET
		updated_at = NOW(),
		episode = $2
	WHERE id = $1
	RETURNING id, created_at, updated_at, episode, user_id, anime_id
	`

	result := &AnimeProgress{}

	if err := tx.QueryRow(
		query,
		reqProgress.Id,
		reqProgress.Episode,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Episode,
		&result.UserId,
		&result.AnimeId,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func DeleteAnimeProgress(tx *sql.Tx, reqProgress *AnimeProgress) error {
	query := `
	DELETE FROM anime_progress
	WHERE id = $1
	`

	queryResult, err := tx.Exec(
		query,
		reqProgress.Id,
	)
	if err != nil {
		return err
	}

	n, err := queryResult.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}
