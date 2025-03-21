package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"strings"
	"testing"
	"testovoe/internal/models"
	"time"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx, "postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithSQLDriver("pgx"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("не получилось запустить контейнер с postgres: %v", err)
	}

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("не вышло взять строку подключения: %v", err)
	}
	t.Logf("строка подключения: %s", dsn)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("не удалось подключиться к базе: %v", err)
	}

	_, err = pool.Exec(ctx, `
		CREATE TABLE users (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("не получилось создать таблицу: %v", err)
	}

	cleanup := func() {
		pool.Close()
		pgContainer.Terminate(ctx)
	}

	return pool, cleanup
}

func TestUserRepository_CreateUser(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(pool)

	user := &models.User{Name: "Test User", Email: "test@example.com"}
	err := repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	var createdUser models.User
	err = pool.QueryRow(context.Background(), "SELECT id, name, email FROM users WHERE id = $1", user.ID).
		Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)

	duplicateUser := &models.User{Name: "Another User", Email: "test@example.com"}
	err = repo.CreateUser(context.Background(), duplicateUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key value violates unique constraint")
}

func TestUserRepository_GetUserByID(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(pool)

	var userID int64
	err := pool.QueryRow(context.Background(), "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", "Test User", "test@example.com").Scan(&userID)
	assert.NoError(t, err)

	user, err := repo.GetUserByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)

	_, err = repo.GetUserByID(context.Background(), 999)
	assert.Equal(t, ErrUserNotFound, err)
}

func TestUserRepository_UpdateUserByID(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(pool)

	var userID int64
	err := pool.QueryRow(context.Background(), "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", "Old User", "old@example.com").Scan(&userID)
	assert.NoError(t, err)

	updateUser := &models.User{Name: "New User", Email: "new@example.com"}
	err = repo.UpdateUserByID(context.Background(), userID, updateUser)
	assert.NoError(t, err)

	var updatedUser models.User
	err = pool.QueryRow(context.Background(), "SELECT id, name, email FROM users WHERE id = $1", userID).
		Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email)
	assert.NoError(t, err)
	assert.Equal(t, "New User", updatedUser.Name)
	assert.Equal(t, "new@example.com", updatedUser.Email)

	// обновляем несуществующего юзера
	err = repo.UpdateUserByID(context.Background(), 999, updateUser)
	assert.Equal(t, ErrUserNotFound, err)

	_, err = pool.Exec(context.Background(), "INSERT INTO users (name, email) VALUES ($1, $2)", "Another User", "another@example.com")
	assert.NoError(t, err)
	err = repo.UpdateUserByID(context.Background(), userID, &models.User{Name: "New User", Email: "another@example.com"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key value violates unique constraint")
}

func TestUserRepository_DeleteUserByID(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(pool)

	var userID int64
	err := pool.QueryRow(context.Background(), "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", "Test User", "test@example.com").Scan(&userID)
	assert.NoError(t, err)

	err = repo.DeleteUserByID(context.Background(), userID)
	assert.NoError(t, err)

	_, err = repo.GetUserByID(context.Background(), userID)
	assert.Equal(t, ErrUserNotFound, err)

	err = repo.DeleteUserByID(context.Background(), 999)
	assert.Equal(t, ErrUserNotFound, err)
}

func TestUserRepository_CreateUser_LongFields(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(pool)

	longName := strings.Repeat("a", 101)
	user := &models.User{Name: longName, Email: "test@example.com"}
	err := repo.CreateUser(context.Background(), user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value too long for type character varying(100)")

	longEmail := strings.Repeat("b", 101) + "@example.com"
	user = &models.User{Name: "Test User", Email: longEmail}
	err = repo.CreateUser(context.Background(), user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "value too long for type character varying(100)")
}
