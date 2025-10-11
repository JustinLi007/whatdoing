package database

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type ServiceDb interface {
	Conn() *sql.DB
	MigrateFS(migrationFS fs.FS, dir string) error
}

type serviceDb struct {
	db *sql.DB
}

var serviceDbInstance *serviceDb

func NewDb(connStr string) (ServiceDb, error) {
	if serviceDbInstance != nil {
		return serviceDbInstance, nil
	}

	conn, err := Open(connStr)
	if err != nil {
		return nil, err
	}

	newServiceDb := &serviceDb{
		db: conn,
	}
	serviceDbInstance = newServiceDb

	return serviceDbInstance, nil
}

func (s *serviceDb) Conn() *sql.DB {
	return s.db
}

func (s *serviceDb) MigrateFS(migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(s.db, dir)
}

func Open(connStr string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}

func Migrate(db *sql.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("migrate: %v", err)
	}

	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("migrate: %v", err)
	}

	return nil
}
