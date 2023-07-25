package bun

import (
	"automator-go/adapters/repositories/bun/models"
	"automator-go/usecases/task"
	"context"
	"fmt"
	"github.com/nlepage/go-cuid2"
	"github.com/uptrace/bun"
)

type BunCaptureMedia struct {
	db  *bun.DB
	ctx *context.Context
}

func NewBunCaptureMedia(db *bun.DB, ctx *context.Context) *BunCaptureMedia {
	return &BunCaptureMedia{db: db, ctx: ctx}
}

func (b *BunCaptureMedia) Save(input task.NewMediaInput) error {
	mediaId, err := cuid2.CreateId()
	if err != nil {
		return fmt.Errorf("error generating media id: %w", err)
	}
	media := models.Media{
		ID:            mediaId,
		Attributes:    input.Attributes,
		Height:        input.Height,
		Width:         input.Width,
		X:             input.X,
		Y:             input.Y,
		Url:           input.Url,
		PHash:         input.PHash,
		Filename:      input.Filename,
		MediaUrl:      input.MediaUrl,
		ScreenshotUrl: input.ScreenshotUrl,
		ResourceUrl:   input.ResourceUrl,
		TaskId:        input.TaskId,
	}

	_, err = b.db.NewInsert().Model(&media).Exec(*b.ctx)
	if err != nil {
		return fmt.Errorf("error inserting media: %w", err)
	}

	return nil
}
