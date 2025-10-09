package database

import (
	"database/sql"
	"log"
	"time"

	"service-user/internal/password"
	"service-user/internal/token"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID          `json:"id"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	Email        string             `json:"email"`
	Password     *password.Password `json:"-"`
	RefreshToken *token.Token       `json:"-"`
	Username     *string            `json:"username"`
	Role         string             `json:"role"`
}

type ServiceUsers interface {
	CreateUser(reqUser *User) (*User, error)
	GetUserById(reqUser *User) (*User, error)
	GetUserByEmailPassword(reqUser *User) (*User, error)
	UpdateUser(reqUser *User) (*User, error)
	DeleteUser(reqUser *User) error
}

type serviceUsers struct {
	db ServiceDb
}

var serviceUsersInstance *serviceUsers

func NewServiceUsers(db ServiceDb) ServiceUsers {
	if serviceUsersInstance != nil {
		return serviceUsersInstance
	}
	newServiceUsers := &serviceUsers{
		db: db,
	}
	serviceUsersInstance = newServiceUsers
	return serviceUsersInstance
}

func (s *serviceUsers) CreateUser(reqUser *User) (*User, error) {
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

	result, err := InsertUser(tx, reqUser)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *serviceUsers) GetUserById(reqUser *User) (*User, error) {
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

	result, err := SelectUserById(tx, reqUser)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *serviceUsers) GetUserByEmailPassword(reqUser *User) (*User, error) {
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

	result, err := SelectUserByEmailPassword(tx, reqUser)
	if err != nil {
		return nil, err
	}

	result.RefreshToken = token.NewToken(token.REFRESH_TOKEN_TTL)

	result, err = UpdateRefreshToken(tx, result)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *serviceUsers) UpdateUser(reqUser *User) (*User, error) {
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

	result, err := UpdateUser(tx, reqUser)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *serviceUsers) DeleteUser(reqUser *User) error {
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

	if err := DeleteUser(tx, reqUser); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func InsertUser(tx *sql.Tx, reqUser *User) (*User, error) {
	query := `
	INSERT INTO users (id, email, password_hash, refresh_token, expiry)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT DO NOTHING
	RETURNING id, created_at, updated_at, email, refresh_token, expiry, username, role 
	`

	result := &User{
		RefreshToken: &token.Token{},
	}

	if err := tx.QueryRow(
		query,
		uuid.New(),
		reqUser.Email,
		reqUser.Password.Hash,
		reqUser.RefreshToken.Hash,
		reqUser.RefreshToken.Expiry,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Email,
		&result.RefreshToken.Hash,
		&result.RefreshToken.Expiry,
		&result.Username,
		&result.Role,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func SelectUserById(tx *sql.Tx, reqUser *User) (*User, error) {
	query := `
	SELECT id, created_at, updated_at, email, refresh_token, expiry, username, role
	FROM users
	WHERE id = $1
	`

	result := &User{
		RefreshToken: &token.Token{},
	}

	if err := tx.QueryRow(
		query,
		reqUser.Id,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Email,
		&result.RefreshToken.Hash,
		&result.RefreshToken.Expiry,
		&result.Username,
		&result.Role,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func SelectUserByEmailPassword(tx *sql.Tx, reqUser *User) (*User, error) {
	query := `
	SELECT id, created_at, updated_at, email, password_hash, refresh_token, expiry, username, role
	FROM users
	WHERE email = $1
	`

	result := &User{
		Password:     &password.Password{},
		RefreshToken: &token.Token{},
	}

	if err := tx.QueryRow(
		query,
		reqUser.Email,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Email,
		&result.Password.Hash,
		&result.RefreshToken.Hash,
		&result.RefreshToken.Expiry,
		&result.Username,
		&result.Role,
	); err != nil {
		return nil, err
	}

	match, err := result.Password.Validate(reqUser.Password.PlainText)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, sql.ErrNoRows
	}

	return result, nil
}

func UpdateUser(tx *sql.Tx, reqUser *User) (*User, error) {
	// TODO:
	query := `
	UPDATE users
	SET
		updated_at = $1,
		username = $2
	WHERE id = $3
	RETURNING id, created_at, updated_at, username, email, role
	`

	result := &User{}

	if err := tx.QueryRow(
		query,
		time.Now().UTC(),
		reqUser.Username,
		reqUser.Id,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Username,
		&result.Email,
		&result.Role,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateRefreshToken(tx *sql.Tx, reqUser *User) (*User, error) {
	query := `
	UPDATE users
	SET
		updated_at = NOW(),
		refresh_token = $2
	WHERE id = $1
	RETURNING id, created_at, updated_at, email, refresh_token, expiry, username, role
	`

	result := &User{
		RefreshToken: &token.Token{},
	}

	if err := tx.QueryRow(
		query,
		reqUser.Id,
		reqUser.RefreshToken.Hash,
	).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Email,
		&result.RefreshToken.Hash,
		&result.RefreshToken.Expiry,
		&result.Username,
		&result.Role,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func DeleteUser(tx *sql.Tx, reqUser *User) error {
	query := `
	DELETE FROM users
	WHERE id = $1
	`

	queryResult, err := tx.Exec(
		query,
		reqUser.Id,
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
