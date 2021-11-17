package userservice

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	userrepository "github.com/jumaroar-globant/go-bootcamp/http/repository/user"
)

type Service interface {
	Authenticate(ctx context.Context, username string, password string) (string, error)
}

type userService struct {
	repository userrepository.UserRepository
	logger     log.Logger
}

func NewService(repository userrepository.UserRepository, logger log.Logger) Service {
	return &userService{
		repository,
		logger,
	}
}

func (s *userService) Authenticate(ctx context.Context, name string, password string) (string, error) {
	logger := log.With(s.logger, "method", "Authenticate")

	userId, err := s.repository.Authenticate(ctx, name, password)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	return userId, nil
}
