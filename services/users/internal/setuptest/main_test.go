package setuptest

import (
	"fmt"
	"log"
	"os"
	"service-user/internal/configs"
	"service-user/internal/database"
	"service-user/migrations"
	"testing"

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
	db, err := database.NewDb(connStr)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if err := db.MigrateFS(migrations.Fs, "."); err != nil {
		log.Fatalf("error: %v", err)
	}

	query := `DELETE FROM users`
	_, err = db.Conn().Exec(query)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	code := m.Run()

	os.Exit(code)
}
