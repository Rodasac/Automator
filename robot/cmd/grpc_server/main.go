package main

import (
	grpcDef "automator-go/grpc"
	grpcController "automator-go/robot/adapters/controllers/grpc"
	bunRepo "automator-go/robot/adapters/repositories/bun"
	utils2 "automator-go/utils"
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

const serviceName = "robot-grpc-automator"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	debugEnv := os.Getenv("APP_DEBUG")
	debug := debugEnv == "true"
	version := os.Getenv("APP_VERSION")

	utils2.StartTrace(serviceName, version, debug)
	defer utils2.ShutdownTrace(ctx)

	ctx, span := utils2.StartSpan(ctx, serviceName, "root")

	logger := utils2.StartLogger(debug)

	db := utils2.OpenDb()

	repo := bunRepo.NewBunCaptureMedia(db)

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	authInterceptor := grpcController.NewAuthInterceptor(logger)
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(authInterceptor.Unary(), otelgrpc.UnaryServerInterceptor()),
		grpc.ChainStreamInterceptor(authInterceptor.Stream(), otelgrpc.StreamServerInterceptor()),
	)
	grpcDef.RegisterMediaServiceServer(s, grpcController.NewGrpcServer(repo, logger))

	go func() {
		logger.Ctx(ctx).Info("Starting server...", zap.Int("port", *port))
		if err := s.Serve(lis); err != nil {
			logger.Ctx(ctx).Fatal("failed to serve: %v", zap.Error(err))
		}
	}()

	<-stopSignal
	span.End()
	s.GracefulStop()
	_ = lis.Close()
	logger.Ctx(ctx).Info("Exiting server...")
}
