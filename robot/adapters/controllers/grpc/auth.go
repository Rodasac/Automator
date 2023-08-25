package grpc

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
)

type AuthInterceptor struct {
	logger         *otelzap.Logger
	apiUser        string
	apiKey         string
	isAuthRequired bool
}

func NewAuthInterceptor(logger *otelzap.Logger) *AuthInterceptor {
	isAuthRequired := os.Getenv("API_AUTH_REQUIRED")
	apiUser := os.Getenv("API_USER")
	apiKey := os.Getenv("API_KEY")

	if isAuthRequired == "true" && (apiUser == "" || apiKey == "") {
		logger.Fatal("API_USER or API_KEY is not set")
	}

	return &AuthInterceptor{
		logger:         logger,
		apiUser:        apiUser,
		apiKey:         apiKey,
		isAuthRequired: isAuthRequired == "true",
	}
}

func (a *AuthInterceptor) authorize(ctx context.Context, method string) error {
	a.logger.Ctx(ctx).Debug("Authorize", zap.String("method", method))

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	header := md.Get("authorization")
	if len(header) == 0 {
		return status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	headerValue := header[0]
	if len(headerValue) < 6 {
		return status.Error(codes.Unauthenticated, "authorization token is not valid")
	}
	headerTokenValue := headerValue[6:]

	rawToken := fmt.Sprintf("%s:%s", a.apiUser, a.apiKey)
	base64TokenString := base64.StdEncoding.EncodeToString([]byte(rawToken))

	if base64TokenString != headerTokenValue {
		return status.Error(codes.Unauthenticated, "authorization token is not valid")
	}

	return nil
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !a.isAuthRequired {
			return handler(ctx, req)
		}

		err := a.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if !a.isAuthRequired {
			return handler(srv, ss)
		}

		err := a.authorize(ss.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, ss)
	}
}
