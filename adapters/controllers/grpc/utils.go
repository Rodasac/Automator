package grpc

import (
	"automator-go/entities/models"
	"google.golang.org/protobuf/types/known/structpb"
	"time"
)

func MapMediaModelToMediaRPC(mediaModel *models.Media) (*Media, error) {
	attributes, err := structpb.NewStruct(mediaModel.Attributes)

	return &Media{
		Id:            mediaModel.Id,
		Attributes:    attributes,
		Height:        mediaModel.Height,
		Width:         mediaModel.Width,
		X:             mediaModel.X,
		Y:             mediaModel.Y,
		Url:           mediaModel.Url,
		Phash:         mediaModel.PHash,
		Filename:      mediaModel.Filename,
		MediaUrl:      mediaModel.MediaUrl,
		ScreenshotUrl: mediaModel.ScreenshotUrl,
		ResourceUrl:   mediaModel.ResourceUrl,
		TaskId:        mediaModel.TaskId,
		CreatedAt:     mediaModel.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     mediaModel.UpdatedAt.Format(time.RFC3339),
	}, err
}

func MapMediaModelToRPC(mediaModel *models.Media) (*MediaResponse, error) {
	media, err := MapMediaModelToMediaRPC(mediaModel)

	return &MediaResponse{
		Media: media,
	}, err
}
