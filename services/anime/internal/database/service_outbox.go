package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type EventStatus string
type EventType string

const (
	STATUS_INCOMPLETE = "incomplete"
	STATUS_PENDING    = "pending"
	STATUS_COMPLETED  = "completed"
)

const (
	EVENT_CREATE = "create"
	EVENT_UPDATE = "update"
	EVENT_DELETE = "create"
)

type Event struct {
	Id        uuid.UUID   `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	AnimeId   uuid.UUID   `json:"anime_id"`
	Name      string      `json:"name"`
	Episodes  int         `json:"episodes"`
	EventType EventType   `json:"event_type"`
	Status    EventStatus `json:"status"`
}

type ServiceOutbox interface {
	GetIncomplete() (*Event, error)
	MarkIncomplete(reqEvent *Event) (*Event, error)
	MarkCompleted(reqEvent *Event) (*Event, error)
}

type serviceOutbox struct {
	db ServiceDb
}

var serviceOutboxInstance *serviceOutbox

func NewServiceOutbox(db ServiceDb) ServiceOutbox {
	if serviceOutboxInstance != nil {
		return serviceOutboxInstance
	}
	newServiceOutbox := &serviceOutbox{
		db: db,
	}
	serviceOutboxInstance = newServiceOutbox
	return serviceOutboxInstance
}

func (o *serviceOutbox) GetIncomplete() (*Event, error) {
	tx, err := o.db.Conn().Begin()
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

	result, err := SelectIncomplete(tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (o *serviceOutbox) MarkIncomplete(reqEvent *Event) (*Event, error) {
	tx, err := o.db.Conn().Begin()
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

	result, err := UpdateEventStatusIncomplete(tx, reqEvent)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (o *serviceOutbox) MarkCompleted(reqEvent *Event) (*Event, error) {
	tx, err := o.db.Conn().Begin()
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

	result, err := UpdateEventStatusCompleted(tx, reqEvent)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func SelectIncomplete(tx *sql.Tx) (*Event, error) {
	query := `
	WITH next_incomplete AS (
		SELECT id FROM outbox
		WHERE status = 'incomplete'
		ORDER BY created_at ASC
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	)
	UPDATE outbox o
	SET
		status = $1
	FROM next_incomplete
	WHERE o.id = next_incomplete.id
	AND o.status = 'incomplete'
	RETURNING o.id, o.created_at, o.updated_at, o.anime_id, o.name, o.episodes, o.event_type, o.status
	`

	result := &Event{}

	if err := tx.QueryRow(
		query,
		STATUS_PENDING,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.AnimeId,
		&result.Name,
		&result.Episodes,
		&result.EventType,
		&result.Status,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateEventStatusIncomplete(tx *sql.Tx, reqEvent *Event) (*Event, error) {
	query := `
	UPDATE outbox
	SET
		updated_at = NOW(),
		status = $2
	WHERE id = $1
	RETURNING id, created_at, updated_at, anime_id, name, episodes, event_type, status
	`

	result := &Event{}

	if err := tx.QueryRow(
		query,
		reqEvent.Id,
		EVENT_STATUS_INCOMPLETE,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.AnimeId,
		&result.Name,
		&result.Episodes,
		&result.EventType,
		&result.Status,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateEventStatusCompleted(tx *sql.Tx, reqEvent *Event) (*Event, error) {
	query := `
	UPDATE outbox
	SET
		updated_at = NOW(),
		status = $2
	WHERE id = $1
	RETURNING id, created_at, updated_at, anime_id, name, episodes, event_type, status
	`

	result := &Event{}

	if err := tx.QueryRow(
		query,
		reqEvent.Id,
		STATUS_COMPLETED,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.AnimeId,
		&result.Name,
		&result.Episodes,
		&result.EventType,
		&result.Status,
	); err != nil {
		return nil, err
	}

	return result, nil
}
