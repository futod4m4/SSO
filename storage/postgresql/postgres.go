package postgresql

import (
	"SSO/internal/domain/models"
	"SSO/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/gopsql/psql"
	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

var connectionString = fmt.Sprintf("postgres://%s:@%s:%d/%s",
	"fedor",
	"localhost",
	5432,
	"sso_for_app",
)

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"

	// Указываем путь до бд
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SaveUser saves user to database.
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte, username, sex, location, dateOfBirth string) (int64, error) {
	const op = "storage.postgresql.SaveUser"

	var id int64
	err := s.db.QueryRowContext(
		ctx,
		"INSERT INTO users(email, pass_hash, username, location, birth_date, sex) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		email, passHash, username, location, dateOfBirth, sex,
	).Scan(&id)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			// Код уникального ограничения
			if pgErr.Code == "23505" {
				return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns user by email.
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgresql.User"

	var user models.User

	err := s.db.QueryRowContext(
		ctx,
		"SELECT id, email, pass_hash, username, location, sex, birth_date FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PassHash, &user.Username, &user.Location, &user.Sex, &user.DateOfBirth)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgresql.IsAdmin"
	var isAdmin bool

	err := s.db.QueryRowContext(
		ctx,
		"SELECT is_admin FROM users WHERE id = $1",
		userID,
	).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (s *Storage) IsExists(ctx context.Context, email string) (bool, error) {
	const op = "storage.postgresql.IsExists"
	var isExists = true

	err := s.db.QueryRowContext(
		ctx,
		"SELECT email FROM users WHERE email = $1",
		email,
	).Scan(&email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			isExists = false
			return isExists, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return isExists, fmt.Errorf("%s: %w", op, err)
	}

	return isExists, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.postgresql.App"
	var app models.App

	err := s.db.QueryRowContext(
		ctx,
		"SELECT id, name, secret FROM apps WHERE id = $1",
		appID,
	).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
