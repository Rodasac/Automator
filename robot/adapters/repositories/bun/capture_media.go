package bun

import (
	bunModels "automator-go/robot/adapters/repositories/bun/models"
	"automator-go/robot/entities/models"
	"automator-go/robot/usecases/task"
	"context"
	"fmt"
	"github.com/nlepage/go-cuid2"
	"github.com/uptrace/bun"
	"time"
)

type CaptureMedia struct {
	db *bun.DB
}

func NewBunCaptureMedia(db *bun.DB) *CaptureMedia {
	return &CaptureMedia{db: db}
}

func (b *CaptureMedia) Save(input task.NewMediaInput, ctx context.Context) error {
	mediaId, err := cuid2.CreateId()
	if err != nil {
		return fmt.Errorf("error generating media id: %w", err)
	}
	media := bunModels.Media{
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

	_, err = b.db.NewInsert().Model(&media).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error inserting media: %w", err)
	}

	return nil
}

func (b *CaptureMedia) GetMedia(mediaId string, ctx context.Context) (*models.Media, error) {
	media := &bunModels.Media{}
	err := b.db.NewSelect().Model(media).Where("id = ?", mediaId).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting media: %w", err)
	}

	return MapBunMediaToModel(media), nil
}

func (b *CaptureMedia) GetMediaByHash(hash string, ctx context.Context) (*models.Media, error) {
	media := &bunModels.Media{}
	err := b.db.NewSelect().Model(media).Where("p_hash = ?", hash).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting media: %w", err)
	}

	return MapBunMediaToModel(media), nil
}

func (b *CaptureMedia) GetMedias(filter *task.MediaFilter, ctx context.Context) ([]*models.Media, error) {
	medias := &[]bunModels.Media{}
	query := b.db.NewSelect().Model(medias)

	if filter.Hash != nil {
		query.Where("phash = ?", *filter.Hash)
	}

	if filter.CreatedAt != nil {
		query.Where("created_at > ?", filter.CreatedAt.Format(time.RFC3339))
	}

	if filter.TaskId != nil {
		query.Where("task_id = ?", *filter.TaskId)
	}

	if filter.Order != nil {
		query.Order("created_at " + string(*filter.Order))
	}

	if filter.Limit != nil {
		query.Limit(int(*filter.Limit))
	}

	err := query.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting medias: %w", err)
	}

	mediasModel := make([]*models.Media, 0, len(*medias))
	for _, media := range *medias {
		mediasModel = append(mediasModel, MapBunMediaToModel(&media))
	}

	return mediasModel, nil
}
