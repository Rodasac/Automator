package grpc

import (
	"automator-go/usecases/task"
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type grpcServer struct {
	UnimplementedMediaServiceServer

	dbRepo task.CapturedMediaRepository
	logger *zap.Logger
}

func NewGrpcServer(dbRepo task.CapturedMediaRepository, logger *zap.Logger) MediaServiceServer {
	return &grpcServer{
		dbRepo: dbRepo,
		logger: logger,
	}
}

func (g *grpcServer) GetMediaById(ctx context.Context, param *MediaIdParam) (*MediaResponse, error) {
	g.logger.Debug("GetMediaById", zap.String("id", param.GetId()))
	mediaModel, err := g.dbRepo.GetMedia(param.GetId(), ctx)
	if err != nil {
		return nil, err
	}

	return MapMediaModelToRPC(mediaModel)
}

func (g *grpcServer) GetMediaByHash(ctx context.Context, param *MediaHashParam) (*MediaResponse, error) {
	g.logger.Debug("GetMediaByHash", zap.String("hash", param.GetPhash()))
	mediaModel, err := g.dbRepo.GetMediaByHash(param.GetPhash(), ctx)
	if err != nil {
		return nil, err
	}

	return MapMediaModelToRPC(mediaModel)
}

func (g *grpcServer) GetMediaList(ctx context.Context, param *MediaFiltersParam) (*MediaListResponse, error) {
	createdAt := new(time.Time)
	if param.CreatedAt != nil {
		parsedTime, err := time.Parse(time.RFC3339, param.GetCreatedAt())
		if err != nil {
			return nil, fmt.Errorf("error parsing time: %w", err)
		}
		*createdAt = parsedTime
	}
	orderBy := new(task.Order)
	if param.GetOrder() == MediaOrder_MEDIA_ORDER_ASC {
		*orderBy = task.ASC
	} else {
		*orderBy = task.DESC
	}

	filters := task.MediaFilter{
		Hash:      param.Hash,
		CreatedAt: createdAt,
		TaskId:    param.TaskId,
		Order:     orderBy,
	}

	mediasModel, err := g.dbRepo.GetMedias(&filters, ctx)
	if err != nil {
		return nil, err
	}

	medias := make([]*Media, 0, len(mediasModel))
	for _, mediaModel := range mediasModel {
		media, err := MapMediaModelToMediaRPC(mediaModel)
		if err != nil {
			return nil, err
		}
		medias = append(medias, media)
	}

	return &MediaListResponse{
		Media: medias,
	}, nil
}
