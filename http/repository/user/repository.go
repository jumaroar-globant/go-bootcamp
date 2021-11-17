package userrepository

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/jumaroar-globant/go-bootcamp/user/pb"

	"google.golang.org/grpc"
)

type userRepository struct {
	conn   *grpc.ClientConn
	logger log.Logger
}

type UserRepository interface {
	Authenticate(ctx context.Context, username string, password string) (string, error)
}

func NewUserRepository(conn *grpc.ClientConn, logger log.Logger) UserRepository {
	return &userRepository{
		conn:   conn,
		logger: log.With(logger, "error", "grpc"),
	}
}

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
