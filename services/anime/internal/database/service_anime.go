package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

const (
	EVENT_TYPE_ANIME_CREATE = "create"
)

const (
	EVENT_STATUS_INCOMPLETE = "incomplete"
)

type Anime struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Episodes  int       `json:"episodes"`
}

type ServiceAnime interface {
	CreateAnime(reqAnime *Anime) (*Anime, error)
	GetAnimeById(reqAnime *Anime) (*Anime, error)
	GetAnimeByName(reqAnime *Anime) (*Anime, error)
	UpdateAnime(reqAnime *Anime) (*Anime, error)
	DeleteAnimeById(reqAnime *Anime) error
}

type serviceAnime struct {
	db ServiceDb
}

var serviceAnimeInstance *serviceAnime

func NewServiceAnime(db ServiceDb) ServiceAnime {
	if serviceAnimeInstance != nil {
		return serviceAnimeInstance
	}
	newServiceAnime := &serviceAnime{
		db: db,
	}
	serviceAnimeInstance = newServiceAnime
	return serviceAnimeInstance
}

func (s *serviceAnime) CreateAnime(reqAnime *Anime) (*Anime, error) {
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

	result, err := InsertAnime(tx, reqAnime)
	if err != nil {
		return nil, err
	}

	if err := InsertAnimeCreateEvent(tx, result); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *serviceAnime) GetAnimeById(reqAnime *Anime) (*Anime, error) {
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

	result, err := SelectAnimeById(tx, reqAnime)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *serviceAnime) GetAnimeByName(reqAnime *Anime) (*Anime, error) {
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

	result, err := SelectAnimeByName(tx, reqAnime)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *serviceAnime) UpdateAnime(reqAnime *Anime) (*Anime, error) {
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

	result, err := UpdateAnime(tx, reqAnime)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *serviceAnime) DeleteAnimeById(reqAnime *Anime) error {
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

	if err := DeleteAnimeById(tx, reqAnime); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func InsertAnime(tx *sql.Tx, reqAnime *Anime) (*Anime, error) {
	query := `
	INSERT INTO anime (id, name, episodes)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at, name, episodes
	`

	result := &Anime{}

	if err := tx.QueryRow(
		query,
		uuid.New(),
		reqAnime.Name,
		reqAnime.Episodes,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Name,
		&result.Episodes,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func SelectAnimeById(tx *sql.Tx, reqAnime *Anime) (*Anime, error) {
	query := `
	SELECT id, created_at, updated_at, name, episodes
	FROM anime
	WHERE id = $1
	`

	result := &Anime{}

	if err := tx.QueryRow(
		query,
		reqAnime.Id,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Name,
		&result.Episodes,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func SelectAnimeByName(tx *sql.Tx, reqAnime *Anime) (*Anime, error) {
	query := `
	SELECT id, created_at, updated_at, name, episodes
	FROM anime
	WHERE name = $1
	`

	result := &Anime{}

	if err := tx.QueryRow(
		query,
		reqAnime.Name,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Name,
		&result.Episodes,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateAnime(tx *sql.Tx, reqAnime *Anime) (*Anime, error) {
	query := `
	UPDATE anime
	SET
		updated_at = NOW(),
		name = $1,
		episodes = $2
	WHERE id = $3
	RETURNING id, created_at, updated_at, name, episodes
	`

	result := &Anime{}

	if err := tx.QueryRow(
		query,
		reqAnime.Name,
		reqAnime.Episodes,
		reqAnime.Id,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Name,
		&result.Episodes,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func DeleteAnimeById(tx *sql.Tx, reqAnime *Anime) error {
	query := `
	DELETE FROM anime
	WHERE id = $1
	`

	queryResult, err := tx.Exec(
		query,
		reqAnime.Id,
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

func InsertAnimeCreateEvent(tx *sql.Tx, reqAnime *Anime) error {
	query := `
	INSERT INTO event_anime (id, anime_id, name, episodes, event_type, status)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	queryResult, err := tx.Exec(
		query,
		uuid.New(),
		reqAnime.Id,
		reqAnime.Name,
		reqAnime.Episodes,
		EVENT_TYPE_ANIME_CREATE,
		EVENT_STATUS_INCOMPLETE,
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
