package main

import (
	grpc3 "automator-go/grpc"
	grpc2 "automator-go/robot/adapters/controllers/grpc"
	bun2 "automator-go/robot/adapters/repositories/bun"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	debug := os.Getenv("APP_DEBUG")

	var zapLogger *zap.Logger
	if debug == "true" {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv("BUNDEBUG"),
	))

	repo := bun2.NewBunCaptureMedia(db)

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	grpc3.RegisterMediaServiceServer(s, grpc2.NewGrpcServer(repo, zapLogger))

	go func() {
		zapLogger.Info("Starting server...", zap.Int("port", *port))
		if err := s.Serve(lis); err != nil {
			zapLogger.Fatal("failed to serve: %v", zap.Error(err))
		}
	}()

	select {
	case <-ctx.Done():
		s.GracefulStop()
		_ = lis.Close()
		zapLogger.Info("Exiting server...")

		return
	}
}
