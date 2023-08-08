package main

import (
	grpcDef "automator-go/grpc"
	grpcController "automator-go/robot/adapters/controllers/grpc"
	bunRepo "automator-go/robot/adapters/repositories/bun"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/extra/bunotel"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/uptrace/uptrace-go/uptrace"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	debug := os.Getenv("APP_DEBUG")

	uptraceDsn := os.Getenv("UPTRACE_DSN")
	if uptraceDsn == "" {
		log.Fatal("UPTRACE_DSN is required")
	}

	version := os.Getenv("APP_VERSION")

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(uptraceDsn),
		uptrace.WithServiceName("robot-rpc-server"),
		uptrace.WithServiceVersion(version),
	)
	defer func(ctx context.Context) {
		_ = uptrace.Shutdown(ctx)
	}(ctx)

	var zapLogger *zap.Logger
	if debug == "true" {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	otelLog := otelzap.New(zapLogger, otelzap.WithMinLevel(zapLogger.Level()))

	dsn := os.Getenv("DATABASE_URL")
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv("BUNDEBUG"),
	))
	db.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName("robot")))

	repo := bunRepo.NewBunCaptureMedia(db)

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	authInterceptor := grpcController.NewAuthInterceptor(otelLog)
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(authInterceptor.Unary(), otelgrpc.UnaryServerInterceptor()),
		grpc.ChainStreamInterceptor(authInterceptor.Stream(), otelgrpc.StreamServerInterceptor()),
	)
	grpcDef.RegisterMediaServiceServer(s, grpcController.NewGrpcServer(repo, otelLog))

	go func() {
		otelLog.Info("Starting server...", zap.Int("port", *port))
		if err := s.Serve(lis); err != nil {
			otelLog.Fatal("failed to serve: %v", zap.Error(err))
		}
	}()

	select {
	case <-ctx.Done():
		s.GracefulStop()
		_ = lis.Close()
		log.Printf("Exiting server...")

		return
	}
}
