package userrepository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/jumaroar-globant/go-bootcamp/user/pb"

	sharedLib "github.com/jumaroar-globant/go-bootcamp/shared"

	"google.golang.org/grpc"
)

var (
	ErrBadAge = errors.New("age is not a number")
)

type userRepository struct {
	client pb.UserServiceClient
	logger log.Logger
}

// UserRepository is the user repository
type UserRepository interface {
	Authenticate(ctx context.Context, username string, password string) (string, error)
	CreateUser(ctx context.Context, user sharedLib.User) (sharedLib.User, error)
	GetUser(ctx context.Context, userID string) (sharedLib.User, error)
	UpdateUser(ctx context.Context, user sharedLib.User) (sharedLib.User, error)
	DeleteUser(ctx context.Context, userID string) (string, error)
}

// NewUserRepository is the UserRepository constructor
func NewUserRepository(conn *grpc.ClientConn, logger log.Logger) UserRepository {
	return &userRepository{
		client: pb.NewUserServiceClient(conn),
		logger: log.With(logger, "error", "grpc"),
	}
}

// Authenticate is the userRepository authentication method
func (r *userRepository) Authenticate(ctx context.Context, username string, pwdHash string) (string, error) {
	logger := log.With(r.logger, "method", "Authenticate")

	request := &pb.UserAuthRequest{
		Username: username,
		Password: pwdHash,
	}

	reply, err := r.client.Authenticate(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	return reply.Message, nil
}

// CreateUser is the userRepository user creation method
func (r *userRepository) CreateUser(ctx context.Context, user sharedLib.User) (sharedLib.User, error) {
	logger := log.With(r.logger, "method", "CreateUser")

	request := &pb.CreateUserRequest{
		Name:                  user.Name,
		Password:              user.Password,
		Age:                   strconv.Itoa(user.Age),
		AdditionalInformation: user.AdditionalInformation,
		Parent:                user.Parents,
	}

	reply, err := r.client.CreateUser(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return sharedLib.User{}, err
	}

	fmt.Println(reply.Age)

	intAge, err := strconv.Atoi(reply.Age)
	if err != nil {
		level.Error(logger).Log("err", ErrBadAge)
		return sharedLib.User{}, ErrBadAge
	}

	return sharedLib.User{
		ID:                    reply.Id,
		Name:                  reply.Name,
		Age:                   intAge,
		AdditionalInformation: reply.AdditionalInformation,
		Parents:               reply.Parent,
	}, nil
}

// GetUser is the userRepository method to retrieve an user by id
func (r *userRepository) GetUser(ctx context.Context, userID string) (sharedLib.User, error) {
	logger := log.With(r.logger, "method", "GetUser")

	request := &pb.GetUserRequest{
		Id: userID,
	}

	reply, err := r.client.GetUser(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return sharedLib.User{}, err
	}

	intAge, err := strconv.Atoi(reply.Age)
	if err != nil {
		level.Error(logger).Log("err", ErrBadAge)
		return sharedLib.User{}, ErrBadAge
	}

	return sharedLib.User{
		ID:                    reply.Id,
		Name:                  reply.Name,
		Age:                   intAge,
		AdditionalInformation: reply.AdditionalInformation,
		Parents:               reply.Parent,
	}, nil
}

// UpdateUser is the userRepository method to update an user
func (r *userRepository) UpdateUser(ctx context.Context, user sharedLib.User) (sharedLib.User, error) {
	logger := log.With(r.logger, "method", "UpdateUser")

	request := &pb.UpdateUserRequest{
		Id:                    user.ID,
		Name:                  user.Name,
		Age:                   strconv.Itoa(user.Age),
		AdditionalInformation: user.AdditionalInformation,
		Parent:                user.Parents,
	}

	reply, err := r.client.UpdateUser(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return sharedLib.User{}, err
	}

	intAge, err := strconv.Atoi(reply.Age)
	if err != nil {
		level.Error(logger).Log("err", ErrBadAge)
		return sharedLib.User{}, ErrBadAge
	}

	return sharedLib.User{
		ID:                    reply.Id,
		Name:                  reply.Name,
		Age:                   intAge,
		AdditionalInformation: reply.AdditionalInformation,
		Parents:               reply.Parent,
	}, nil
}

// DeleteUser is the userRepository method to delete an user by id
func (r *userRepository) DeleteUser(ctx context.Context, userID string) (string, error) {
	logger := log.With(r.logger, "method", "DeleteUser")

	request := &pb.DeleteUserRequest{
		Id: userID,
	}

	reply, err := r.client.DeleteUser(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	return reply.Message, nil
}
