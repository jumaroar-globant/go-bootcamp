package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	userendpoints "github.com/jumaroar-globant/go-bootcamp/http/endpoints/user"
	userrepository "github.com/jumaroar-globant/go-bootcamp/http/repository/user"
	userservice "github.com/jumaroar-globant/go-bootcamp/http/service/user"
	"github.com/jumaroar-globant/go-bootcamp/http/transport"
	"google.golang.org/grpc"
)

func main() {
	var (
		grpcUserServiceAddr = flag.String("addr", "localhost:50051", "The gprcUserServer address in the format of host:port")
		httpAddr            = flag.String("http", ":8080", "http listen address")
	)
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "httpService",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "http service started")
	defer level.Info(logger).Log("msg", "http service ended")

	flag.Parse()

	var err error
	var grpcUserServiceConn *grpc.ClientConn
	{
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())
		grpcUserServiceConn, err = grpc.Dial(*grpcUserServiceAddr, opts...)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	var srv userservice.Service
	{
		repository := userrepository.NewUserRepository(grpcUserServiceConn, logger)
		srv = userservice.NewService(repository, logger)
	}

	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	endpoints := userendpoints.MakeEndpoints(srv)
	go func() {
		httpHandler := transport.NewHTTPServer(endpoints, logger)
		errChan <- http.ListenAndServe(*httpAddr, httpHandler)
	}()

	level.Error(logger).Log("exit", <-errChan)
}
