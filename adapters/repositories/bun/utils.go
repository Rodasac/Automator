package bun

import (
	bunModels "automator-go/adapters/repositories/bun/models"
	"automator-go/entities/models"
	"time"
)

func MapBunMediaToModel(media *bunModels.Media) *models.Media {
	var deletedAt *time.Time
	if !media.DeletedAt.IsZero() {
		deletedAt = &media.DeletedAt.Time
	}

	return &models.Media{
		Id:            media.ID,
		Attributes:    media.Attributes,
		Height:        media.Height,
		Width:         media.Width,
		X:             media.X,
		Y:             media.Y,
		Url:           media.Url,
		PHash:         media.PHash,
		Filename:      media.Filename,
		MediaUrl:      media.MediaUrl,
		ScreenshotUrl: media.ScreenshotUrl,
		ResourceUrl:   media.ResourceUrl,
		TaskId:        media.TaskId,
		CreatedAt:     media.CreatedAt,
		UpdatedAt:     media.UpdatedAt,
		DeletedAt:     deletedAt,
	}
}
