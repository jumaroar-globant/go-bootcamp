package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/jumaroar-globant/go-bootcamp/user/config"
	"github.com/jumaroar-globant/go-bootcamp/user/endpoints"
	"github.com/jumaroar-globant/go-bootcamp/user/pb"
	"github.com/jumaroar-globant/go-bootcamp/user/repository"
	"github.com/jumaroar-globant/go-bootcamp/user/service"
	"github.com/jumaroar-globant/go-bootcamp/user/shared"

	transport "github.com/jumaroar-globant/go-bootcamp/user/transports"
	"google.golang.org/grpc"
)

var (
	dbDriver = shared.GetStringEnvVar("MICROSERVICE_DATABASE_DRIVER", "mysql")
)

func initDatabase() (interface{}, error) {
	if dbDriver == "mongo" {
		return config.ConnectToMongoDB()
	}

	return config.Connect()
}

func initRepository(logger log.Logger) (repository.UserRepository, error) {
	db, err := initDatabase()
	if err != nil {
		return nil, err
	}

	if dbDriver == "mongo" {
		return repository.NewMongoUserRepository(db.(config.MongoDatabaseHelper), logger), nil
	}

	return repository.NewUserRepository(db.(*sql.DB), logger), nil
}

func main() {
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	userRepository, err := initRepository(logger)
	if err != nil {
		level.Error(logger).Log("error_connecting_to_database", err)
		return
	}

	userService := service.NewUserService(userRepository, logger)
	userEndpoints := endpoints.MakeEndpoints(userService)
	grpcServer := transport.NewGRPCServer(userEndpoints, logger)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	grpcListener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log("during", "Listen", "err", err)
		os.Exit(1)
	}

	go func() {
		baseServer := grpc.NewServer()
		pb.RegisterUserServiceServer(baseServer, grpcServer)
		level.Info(logger).Log("msg", "Server started successfully 🚀")
		baseServer.Serve(grpcListener)
	}()

	level.Error(logger).Log("exit", <-errs)
}
