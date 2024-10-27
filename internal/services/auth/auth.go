package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
	"usekit-auth/internal/domain/models"
	"usekit-auth/internal/lib/jwt"
	"usekit-auth/internal/storage"
)

var (
	ErrInvalidAppId       = errors.New("invalid app_id")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type Auth struct {
	logger      *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (id int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

// New возвращает новый инстанс сервиса Auth
func New(
	logger *slog.Logger,
	userSaver UserSaver,
	UserProvider UserProvider,
	AppProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		logger:      logger,
		usrSaver:    userSaver,
		usrProvider: UserProvider,
		appProvider: AppProvider,
		tokenTTL:    tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system.
//
// If user exists, but password is incorrect, returns error.
// If users doesn't exist, returns error.
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appId int,
) (string, error) {
	const op = "services/auth.Login"

	logger := a.logger.With(slog.String("operation", op))
	logger.Info("attempting to login")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			a.logger.Warn("user not found", err)
			return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		a.logger.Warn("failed to login", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.logger.Info("invalid credential", err)
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op)
	}

	logger.Info("successfully logged in")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.logger.Error("failed to generate token", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// RegisterNewUser registers new user in the system and returns user id.
// If user with given username exists, returns error.
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email,
	password string,
) (int64, error) {
	const op = "services/auth.RegisterNewUser"

	logger := a.logger.With(slog.String("operation", op))
	logger.Info("start register new user")
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to generate password hash", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		logger.Error("failed to save new user", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("user registered successfully")
	return id, nil
}

// IsAdmin returns true if user with given userId is admin.
func (a *Auth) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "services/auth.IsAdmin"

	logger := a.logger.With(slog.String("operation", op))

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			logger.Warn("app not found", err)
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}
		logger.Error("failed to check if user is admin", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))
	return isAdmin, nil
}
