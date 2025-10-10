package database

import (
	"fmt"
	"testing"

	"github.com/JustinLi007/whatdoing/services/anime/internal/configs"
	"github.com/JustinLi007/whatdoing/services/anime/migrations"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupTestDb(t *testing.T) ServiceDb {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatalf("error: %v", err)
	}
	configs := configs.NewConfigs()
	configs.ModeEnv = "DEBUG"
	if err := configs.LoadEnv(); err != nil {
		t.Fatalf("error: %v", err)
	}
	connStr := configs.ConfigDb.PostgresConnStr()
	if connStr == "" {
		t.Fatalf("error: %v", fmt.Errorf("invalid conn str"))
	}
	db, err := NewDb(connStr)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if err := db.MigrateFS(migrations.Fs, "."); err != nil {
		t.Fatalf("error: %v", err)
	}

	query := `
	TRUNCATE
	anime,
	event_anime
	CASCADE
	`

	if _, err = db.Conn().Exec(query); err != nil {
		t.Fatalf("failed truncating test db: %v", err)
	}
	return db
}

func TestCreateAnime(t *testing.T) {
	db := SetupTestDb(t)
	serviceAnime := NewServiceAnime(db)

	reqAnime := &Anime{
		Name:     "anime 1",
		Episodes: 10,
	}
	newAnime, err := serviceAnime.CreateAnime(reqAnime)
	require.NoError(t, err)
	require.NotNil(t, newAnime)
	assert.Equal(t, newAnime.Name, "anime 1")
	assert.Equal(t, newAnime.Episodes, 10)
}

func TestGetAnimeById(t *testing.T) {
	db := SetupTestDb(t)
	serviceAnime := NewServiceAnime(db)

	reqAnime := &Anime{
		Name:     "anime 1",
		Episodes: 10,
	}
	newAnime, err := serviceAnime.CreateAnime(reqAnime)
	require.NoError(t, err)
	require.NotNil(t, newAnime)
	assert.Equal(t, newAnime.Name, "anime 1")
	assert.Equal(t, newAnime.Episodes, 10)

	existingAnime, err := serviceAnime.GetAnimeById(newAnime)
	require.NoError(t, err)
	require.NotNil(t, existingAnime)
	assert.Equal(t, existingAnime.Id, newAnime.Id)
	assert.Equal(t, existingAnime.CreatedAt, newAnime.CreatedAt)
	assert.Equal(t, existingAnime.UpdatedAt, newAnime.UpdatedAt)
	assert.Equal(t, existingAnime.Name, newAnime.Name)
	assert.Equal(t, existingAnime.Episodes, newAnime.Episodes)
}

func TestGetAnimeByName(t *testing.T) {
	db := SetupTestDb(t)
	serviceAnime := NewServiceAnime(db)

	reqAnime := &Anime{
		Name:     "anime 1",
		Episodes: 10,
	}
	newAnime, err := serviceAnime.CreateAnime(reqAnime)
	require.NoError(t, err)
	require.NotNil(t, newAnime)
	assert.Equal(t, newAnime.Name, "anime 1")
	assert.Equal(t, newAnime.Episodes, 10)

	existingAnime, err := serviceAnime.GetAnimeByName(newAnime)
	require.NoError(t, err)
	require.NotNil(t, existingAnime)
	assert.Equal(t, existingAnime.Id, newAnime.Id)
	assert.Equal(t, existingAnime.CreatedAt, newAnime.CreatedAt)
	assert.Equal(t, existingAnime.UpdatedAt, newAnime.UpdatedAt)
	assert.Equal(t, existingAnime.Name, newAnime.Name)
	assert.Equal(t, existingAnime.Episodes, newAnime.Episodes)
}

func TestUpdateAnime(t *testing.T) {
	db := SetupTestDb(t)
	serviceAnime := NewServiceAnime(db)

	reqAnime := &Anime{
		Name:     "anime 1",
		Episodes: 10,
	}
	newAnime, err := serviceAnime.CreateAnime(reqAnime)
	require.NoError(t, err)
	require.NotNil(t, newAnime)
	assert.Equal(t, newAnime.Name, "anime 1")
	assert.Equal(t, newAnime.Episodes, 10)

	existingAnime, err := serviceAnime.GetAnimeByName(newAnime)
	require.NoError(t, err)
	require.NotNil(t, existingAnime)
	assert.Equal(t, existingAnime.Id, newAnime.Id)
	assert.Equal(t, existingAnime.CreatedAt, newAnime.CreatedAt)
	assert.Equal(t, existingAnime.UpdatedAt, newAnime.UpdatedAt)
	assert.Equal(t, existingAnime.Name, newAnime.Name)
	assert.Equal(t, existingAnime.Episodes, newAnime.Episodes)

	existingAnime.Name = "something"
	existingAnime.Episodes = 1243
	updatedAnime, err := serviceAnime.UpdateAnime(existingAnime)
	require.NoError(t, err)
	require.NotNil(t, updatedAnime)
	assert.Equal(t, updatedAnime.Id, newAnime.Id)
	assert.Equal(t, updatedAnime.CreatedAt, newAnime.CreatedAt)
	assert.NotEqual(t, updatedAnime.UpdatedAt, newAnime.UpdatedAt)
	assert.NotEqual(t, updatedAnime.Name, newAnime.Name)
	assert.NotEqual(t, updatedAnime.Episodes, newAnime.Episodes)
	assert.Equal(t, updatedAnime.Name, "something")
	assert.Equal(t, updatedAnime.Episodes, 1243)
}

func TestDeleteAnimeById(t *testing.T) {
	db := SetupTestDb(t)
	serviceAnime := NewServiceAnime(db)

	reqAnime := &Anime{
		Name:     "anime 1",
		Episodes: 10,
	}
	newAnime, err := serviceAnime.CreateAnime(reqAnime)
	require.NoError(t, err)
	require.NotNil(t, newAnime)
	assert.Equal(t, newAnime.Name, "anime 1")
	assert.Equal(t, newAnime.Episodes, 10)

	existingAnime, err := serviceAnime.GetAnimeByName(newAnime)
	require.NoError(t, err)
	require.NotNil(t, existingAnime)
	assert.Equal(t, existingAnime.Id, newAnime.Id)
	assert.Equal(t, existingAnime.CreatedAt, newAnime.CreatedAt)
	assert.Equal(t, existingAnime.UpdatedAt, newAnime.UpdatedAt)
	assert.Equal(t, existingAnime.Name, newAnime.Name)
	assert.Equal(t, existingAnime.Episodes, newAnime.Episodes)

	existingAnime.Name = "something"
	existingAnime.Episodes = 1243
	updatedAnime, err := serviceAnime.UpdateAnime(existingAnime)
	require.NoError(t, err)
	require.NotNil(t, updatedAnime)
	assert.Equal(t, updatedAnime.Id, newAnime.Id)
	assert.Equal(t, updatedAnime.CreatedAt, newAnime.CreatedAt)
	assert.NotEqual(t, updatedAnime.UpdatedAt, newAnime.UpdatedAt)
	assert.NotEqual(t, updatedAnime.Name, newAnime.Name)
	assert.NotEqual(t, updatedAnime.Episodes, newAnime.Episodes)
	assert.Equal(t, updatedAnime.Name, "something")
	assert.Equal(t, updatedAnime.Episodes, 1243)

	err = serviceAnime.DeleteAnimeById(existingAnime)
	require.NoError(t, err)

	deletedAnime, err := serviceAnime.GetAnimeByName(existingAnime)
	require.Error(t, err)
	require.Nil(t, deletedAnime)
	deletedAnime, err = serviceAnime.GetAnimeById(existingAnime)
	require.Error(t, err)
	require.Nil(t, deletedAnime)
}
