package database

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/JustinLi007/whatdoing/services/users/internal/configs"
	"github.com/JustinLi007/whatdoing/services/users/migrations"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("error: %v", err)
	}
	configs := configs.NewConfigs()
	configs.Mode = "DEBUG"
	if err := configs.LoadEnv(); err != nil {
		log.Fatalf("error: %v", err)
	}
	connStr := configs.ConfigDb.PostgresConnStr()
	if connStr == "" {
		log.Fatalf("error: %v", fmt.Errorf("invalid conn str"))
	}
	db, err := NewDb(connStr)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := db.MigrateFS(migrations.Fs, "."); err != nil {
		log.Fatalf("error: %v", err)
	}

	code := m.Run()

	os.Exit(code)
}
