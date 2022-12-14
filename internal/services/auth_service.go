package services

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"short_url/internal/models"
	"short_url/internal/security"
	log "short_url/pkg/logger"
)

// AuthServiceConfig Конфигурация к AuthService
type AuthServiceConfig struct {
	AuthRepo	authRepository
	SubRepo		subRepository
	Logger		*log.Log
}

// AuthService Управляет регистрацией и аутентификацией пользователей
type AuthService struct {
	authRepo	authRepository
	subRepo		subRepository
	logger		*log.Log
}

// Конструктор для AuthService
func NewAuthService(c *AuthServiceConfig) *AuthService {
	return &AuthService{
		authRepo:	c.AuthRepo,
		subRepo:	c.SubRepo,
		logger:		c.Logger,
	}
}

// SignInUserByName вызывает методы других слоев, которые позволят войти пользователю по его имени
func (s *AuthService) SignInUserByName(ctx context.Context, dto models.SignInUserDTO) (models.SignInUserDTO, error) {
	ctx = log.ContextWithSpan(ctx, "SignInUserByName")
	l := s.logger.WithContext(ctx)

	l.Debug("SignInUserByName() started")
	defer l.Debug("SignInUserByName() done")

	// Ищем зарегистрированного пользователя
	u, err := s.authRepo.FindByUsername(ctx, dto.Username)
	if err != nil {
		l.Errorf("Unable to find user. Error: %e", err)

		return dto, err
	}

	// Сравниваем пароли пользователя
	ok, err := security.ComparePasswords(u.Password, dto.Password)
	if err != nil {
		l.Errorf("Unable to compare password. Error: %e", err)

		return dto, err
	}
	if !ok {
		return dto, errors.New("invalid password")
	}

	// Проверяем, есть ли у пользователя подписка
	_, ok = s.subRepo.FindSubscribe(ctx, dto.Username)
	if ok {
		dto.Subscribe	= 1
	} else {
		dto.Subscribe	= 2
	}

	// Мапим данные из db в dto структуру
	dto.FirstName	= u.FirstName
	dto.LastName	= u.LastName
	

	return dto, nil
}

// SignUpUser вызывает методы для создания пользователя и хеширования пароля
func (s *AuthService) SignUpUser(ctx context.Context, dto models.SignUpUserDTO) error {
	ctx = log.ContextWithSpan(ctx, "SignUpUser")
	l := s.logger.WithContext(ctx)

	l.Debug("SignUpUser() started")
	defer l.Debug("SignUpUser() done")

	// Проверяем, зарегистрирован ли пользователь
	_, err := s.authRepo.FindByUsername(ctx, dto.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			l.Info("user not found")
		} else {
			l.Errorf("Unable to find user. Error: %e", err)
			return err
		}
	} else {
		return errors.New("user exist")
	}

	// Хешируем пароль
	hashPassword, err := security.HashPassword(dto.Password)
	if err != nil {
		l.Errorf("Unable to hash password. Error: %e", err)
		return err
	}

	// Создаем пользователя в базе данных
	err = s.authRepo.CreateUser(ctx, models.UserDB{
		Username:	dto.Username,
		FirstName:	dto.FirstName,
		LastName:	dto.LastName,
		Password:	hashPassword,
	})
	if err != nil {
		l.Errorf("Unable to create user. Error: %e", err)
		return err
	}

	return nil
}
