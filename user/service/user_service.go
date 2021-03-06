package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/jumaroar-globant/go-bootcamp/user/pb"
	"github.com/jumaroar-globant/go-bootcamp/user/repository"
	"github.com/jumaroar-globant/go-bootcamp/user/shared"

	sharedLib "github.com/jumaroar-globant/go-bootcamp/shared"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const (
	userDeletedString = "user deleted successfully"
)

var (
	ErrMissingUserName = errors.New("missing username")
	ErrMissingPassword = errors.New("missing password")
	ErrMissingUserID   = errors.New("missing user id")
)

type userService struct {
	repository repository.UserRepository
	logger     log.Logger
}

// UserService interface describes a user service
type UserService interface {
	Authenticate(ctx context.Context, authenticationRequest *pb.UserAuthRequest) (string, error)
	CreateUser(ctx context.Context, userRequest *pb.CreateUserRequest) (sharedLib.User, error)
	UpdateUser(context.Context, *pb.UpdateUserRequest) (sharedLib.User, error)
	GetUser(context.Context, *pb.GetUserRequest) (sharedLib.User, error)
	DeleteUser(context.Context, *pb.DeleteUserRequest) (string, error)
}

// NewService returns a Service with all of the expected dependencies
func NewUserService(userRep repository.UserRepository, logger log.Logger) UserService {
	return &userService{
		repository: userRep,
		logger:     logger,
	}
}

// Authenticate is the userService method to authenticate
func (s *userService) Authenticate(ctx context.Context, authenticationRequest *pb.UserAuthRequest) (string, error) {
	logger := log.With(s.logger, "method", "Authenticate")

	err := s.repository.Authenticate(ctx, authenticationRequest.Username, authenticationRequest.Password)
	if err != nil {
		level.Error(logger).Log("err", err)

		return "", err
	}

	return "User authenticated!", nil
}

// CreateUser is the userService method to create a user
func (s *userService) CreateUser(ctx context.Context, createUserRequest *pb.CreateUserRequest) (sharedLib.User, error) {
	logger := log.With(s.logger, "method", "CreateUser")

	if createUserRequest.Name == "" {
		return sharedLib.User{}, ErrMissingUserName
	}

	if createUserRequest.Password == "" {
		return sharedLib.User{}, ErrMissingPassword
	}

	age, err := strconv.Atoi(createUserRequest.Age)
	if err != nil {
		level.Error(logger).Log("error_converting_age_to_integer", err)

		return sharedLib.User{}, err
	}

	passwordHash, err := shared.HashPassword(createUserRequest.Password)
	if err != nil {
		level.Error(logger).Log("error_hashing_password", err)

		return sharedLib.User{}, err
	}

	user := sharedLib.User{
		ID:                    shared.GenerateID("USR"),
		Name:                  createUserRequest.Name,
		Password:              passwordHash,
		Age:                   age,
		AdditionalInformation: createUserRequest.AdditionalInformation,
		Parents:               createUserRequest.Parent,
	}

	err = s.repository.CreateUser(ctx, user)
	if err != nil {
		level.Error(logger).Log("error_creating_user_in_database", err)

		return sharedLib.User{}, err
	}

	return user, nil
}

// UpdateUser is the userService method to update a user
func (s *userService) UpdateUser(ctx context.Context, updateUserRequest *pb.UpdateUserRequest) (sharedLib.User, error) {
	logger := log.With(s.logger, "method", "UpdateUser")

	if updateUserRequest.Id == "" {
		return sharedLib.User{}, ErrMissingUserID
	}

	age, err := strconv.Atoi(updateUserRequest.Age)
	if err != nil {
		level.Error(logger).Log("error_converting_age_to_integer", err)

		return sharedLib.User{}, err
	}

	user := sharedLib.User{
		ID:                    updateUserRequest.Id,
		Name:                  updateUserRequest.Name,
		Age:                   age,
		AdditionalInformation: updateUserRequest.AdditionalInformation,
		Parents:               updateUserRequest.Parent,
	}

	updatedUser, err := s.repository.UpdateUser(ctx, user)
	if err != nil {
		level.Error(logger).Log("error_updating_user_in_database", err)

		return sharedLib.User{}, err
	}

	return updatedUser, nil
}

// GetUser is the userService method to get a user
func (s *userService) GetUser(ctx context.Context, getUserRequest *pb.GetUserRequest) (sharedLib.User, error) {
	logger := log.With(s.logger, "method", "GetUser")

	if getUserRequest.Id == "" {
		return sharedLib.User{}, ErrMissingUserID
	}

	user, err := s.repository.GetUser(ctx, getUserRequest.Id)
	if err != nil {
		level.Error(logger).Log("error_getting_user_from_database", err)

		return sharedLib.User{}, err
	}

	return user, nil
}

// DeleteUser is the userService method to delete a user
func (s *userService) DeleteUser(ctx context.Context, deleteUserRequest *pb.DeleteUserRequest) (string, error) {
	logger := log.With(s.logger, "method", "DeleteUser")

	if deleteUserRequest.Id == "" {
		return "", ErrMissingUserID
	}

	err := s.repository.DeleteUser(ctx, deleteUserRequest.Id)
	if err != nil {
		level.Error(logger).Log("error_deleting_user_from_database", err)

		return "", err
	}

	return userDeletedString, nil
}
