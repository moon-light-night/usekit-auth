package sqlite

// слой работы с данными

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	uuidV4 "github.com/satori/go.uuid"
	"log/slog"
	"usekit-auth/internal/domain/models"
	"usekit-auth/internal/storage"
)

type Storage struct {
	db *sql.DB
}

// New creates a new instance of SQLite storage
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	// указываем путь к файлу с бд
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (string, error) {
	const op = "storage.sqlite.SaveUser"
	uuid := uuidV4.NewV4().String()
	stmt, err := s.db.Prepare(`INSERT INTO users (uuid, email, pass_hash) VALUES (?, ?, ?)`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, uuid, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	slog.Info("SaveUser res", res)

	// TODO: посмотреть, что возвращается в res и возвращать созданный uuid
	return uuid, nil
}

// User returns user by email
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"
	stmt, err := s.db.Prepare(`SELECT uuid, email, pass_hash FROM users WHERE email = ?`)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)
	var user models.User
	// получаем результат методом Scan и записываем значения из колонок найденной строки в поля объекта user
	err = row.Scan(&user.Uuid, &user.Email, &user.PassHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
