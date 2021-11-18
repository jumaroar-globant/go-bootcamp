package userrepository

import (
	"context"
	"errors"
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
	conn   *grpc.ClientConn
	logger log.Logger
}

// UserRepository is the user repository
type UserRepository interface {
	Authenticate(ctx context.Context, username string, password string) (string, error)
	CreateUser(ctx context.Context, user *sharedLib.User) (*sharedLib.User, error)
	GetUser(ctx context.Context, userID string) (*sharedLib.User, error)
}

// NewUserRepository is the UserRepository constructor
func NewUserRepository(conn *grpc.ClientConn, logger log.Logger) UserRepository {
	return &userRepository{
		conn:   conn,
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

	client := pb.NewUserServiceClient(r.conn)

	reply, err := client.Authenticate(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	return reply.Message, nil
}

// CreateUser is the userRepository user creation method
func (r *userRepository) CreateUser(ctx context.Context, user *sharedLib.User) (*sharedLib.User, error) {
	logger := log.With(r.logger, "method", "CreateUser")

	request := &pb.CreateUserRequest{
		Name:                  user.Name,
		Password:              user.Password,
		Age:                   strconv.Itoa(user.Age),
		AdditionalInformation: user.AdditionalInformation,
		Parent:                user.Parents,
	}

	client := pb.NewUserServiceClient(r.conn)

	reply, err := client.CreateUser(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	intAge, err := strconv.Atoi(reply.Age)
	if err != nil {
		level.Error(logger).Log("err", ErrBadAge)
		return nil, ErrBadAge
	}

	return &sharedLib.User{
		ID:                    reply.Id,
		Name:                  reply.Name,
		Age:                   intAge,
		AdditionalInformation: reply.AdditionalInformation,
		Parents:               reply.Parent,
	}, nil
}

// GetUser is the userRepository method to retrieve an user by id
func (r *userRepository) GetUser(ctx context.Context, userID string) (*sharedLib.User, error) {
	logger := log.With(r.logger, "method", "GetUser")

	request := &pb.GetUserRequest{
		Id: userID,
	}

	client := pb.NewUserServiceClient(r.conn)

	reply, err := client.GetUser(ctx, request)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, err
	}

	intAge, err := strconv.Atoi(reply.Age)
	if err != nil {
		level.Error(logger).Log("err", ErrBadAge)
		return nil, ErrBadAge
	}

	return &sharedLib.User{
		ID:                    reply.Id,
		Name:                  reply.Name,
		Age:                   intAge,
		AdditionalInformation: reply.AdditionalInformation,
		Parents:               reply.Parent,
	}, nil
}
