package grpc

import (
	"automator-go/grpc"
	"automator-go/robot/usecases/task"
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type grpcServer struct {
	grpc.UnimplementedMediaServiceServer

	dbRepo task.CapturedMediaRepository
	logger *zap.Logger
}

func NewGrpcServer(dbRepo task.CapturedMediaRepository, logger *zap.Logger) grpc.MediaServiceServer {
	return &grpcServer{
		dbRepo: dbRepo,
		logger: logger,
	}
}

func (g *grpcServer) GetMediaById(ctx context.Context, param *grpc.MediaIdParam) (*grpc.MediaResponse, error) {
	g.logger.Debug("GetMediaById", zap.String("id", param.GetId()))
	mediaModel, err := g.dbRepo.GetMedia(param.GetId(), ctx)
	if err != nil {
		return nil, err
	}

	return MapMediaModelToRPC(mediaModel)
}

func (g *grpcServer) GetMediaByHash(ctx context.Context, param *grpc.MediaHashParam) (*grpc.MediaResponse, error) {
	g.logger.Debug("GetMediaByHash", zap.String("hash", param.GetPhash()))
	mediaModel, err := g.dbRepo.GetMediaByHash(param.GetPhash(), ctx)
	if err != nil {
		return nil, err
	}

	return MapMediaModelToRPC(mediaModel)
}

func (g *grpcServer) GetMediaList(ctx context.Context, param *grpc.MediaFiltersParam) (*grpc.MediaListResponse, error) {
	createdAt := new(time.Time)
	if param.CreatedAt != nil {
		parsedTime, err := time.Parse(time.RFC3339, param.GetCreatedAt())
		if err != nil {
			return nil, fmt.Errorf("error parsing time: %w", err)
		}
		*createdAt = parsedTime
	}
	orderBy := new(task.Order)
	if param.GetOrder() == grpc.MediaOrder_MEDIA_ORDER_ASC {
		*orderBy = task.ASC
	} else {
		*orderBy = task.DESC
	}

	filters := task.MediaFilter{
		Hash:      param.Hash,
		CreatedAt: createdAt,
		TaskId:    param.TaskId,
		Order:     orderBy,
		Limit:     param.Limit,
	}

	mediasModel, err := g.dbRepo.GetMedias(&filters, ctx)
	if err != nil {
		return nil, err
	}

	medias := make([]*grpc.Media, 0, len(mediasModel))
	for _, mediaModel := range mediasModel {
		media, err := MapMediaModelToMediaRPC(mediaModel)
		if err != nil {
			return nil, err
		}
		medias = append(medias, media)
	}

	return &grpc.MediaListResponse{
		Media: medias,
	}, nil
}
