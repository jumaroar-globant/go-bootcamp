package userservice

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	userrepository "github.com/jumaroar-globant/go-bootcamp/http/repository/user"
	"github.com/jumaroar-globant/go-bootcamp/shared"
)

// Service is the user service
type Service interface {
	Authenticate(ctx context.Context, username string, password string) (string, error)
	CreateUser(ctx context.Context, user shared.User) (shared.User, error)
	GetUser(ctx context.Context, userID string) (shared.User, error)
	UpdateUser(ctx context.Context, user shared.User) (shared.User, error)
	DeleteUser(ctx context.Context, userID string) (string, error)
}

type userService struct {
	repository userrepository.UserRepository
	logger     log.Logger
}

//NewService is the Service constructor
func NewService(repository userrepository.UserRepository, logger log.Logger) Service {
	return &userService{
		repository,
		logger,
	}
}

//Authenticate is a method to athenticate a user
func (s *userService) Authenticate(ctx context.Context, name string, password string) (string, error) {
	logger := log.With(s.logger, "method", "Authenticate")

	message, err := s.repository.Authenticate(ctx, name, password)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	return message, nil
}

//CreateUser is a method to create a user
func (s *userService) CreateUser(ctx context.Context, user shared.User) (shared.User, error) {
	logger := log.With(s.logger, "method", "CreateUser")

	userCreated, err := s.repository.CreateUser(ctx, user)
	if err != nil {
		level.Error(logger).Log("err", err)
		return shared.User{}, err
	}

	return userCreated, nil
}

//GetUser is a method to get a user by id
func (s *userService) GetUser(ctx context.Context, userID string) (shared.User, error) {
	logger := log.With(s.logger, "method", "GetUser")

	user, err := s.repository.GetUser(ctx, userID)
	if err != nil {
		level.Error(logger).Log("err", err)
		return shared.User{}, err
	}

	return user, nil
}

//UpdateUser is a method to update a user
func (s *userService) UpdateUser(ctx context.Context, user shared.User) (shared.User, error) {
	logger := log.With(s.logger, "method", "UpdateUser")

	userCreated, err := s.repository.UpdateUser(ctx, user)
	if err != nil {
		level.Error(logger).Log("err", err)
		return shared.User{}, err
	}

	return userCreated, nil
}

//DeleteUser is a method to delete a user
func (s *userService) DeleteUser(ctx context.Context, userID string) (string, error) {
	logger := log.With(s.logger, "method", "DeleteUser")

	message, err := s.repository.DeleteUser(ctx, userID)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	return message, nil
}
