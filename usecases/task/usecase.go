package task

import (
	"automator-go/entities/models"
	"automator-go/usecases/hasher"
	"context"
	"strings"
)

type Processor struct {
	automatorTaskAdapter AutomatorTaskAdapter
	capturedMediaRepo    CapturedMediaRepository
	storageMediaAdapter  StorageMediaAdapter
	imageHasher          hasher.ImageHasher
}

func NewProcessor(
	automatorTaskAdapter AutomatorTaskAdapter,
	capturedMediaRepo CapturedMediaRepository,
	storageMediaAdapter StorageMediaAdapter,
	imageHasher hasher.ImageHasher,
) *Processor {
	return &Processor{
		automatorTaskAdapter: automatorTaskAdapter,
		capturedMediaRepo:    capturedMediaRepo,
		storageMediaAdapter:  storageMediaAdapter,
		imageHasher:          imageHasher,
	}
}

func (p *Processor) Process(task *models.Task, ctx context.Context) error {
	mediaResult, err := p.automatorTaskAdapter.Run(task)
	if err != nil {
		return err
	}

	if mediaResult == nil {
		return nil
	}

	for _, mediaResult := range *mediaResult {
		hash, err := p.imageHasher.Hash(mediaResult.Media)
		if err != nil {
			return err
		}

		hashWithoutKind := strings.Split(hash, ":")[1]

		storageMedia, err := p.storageMediaAdapter.SaveMedia(hashWithoutKind, &mediaResult)
		if err != nil {
			return err
		}

		err = p.capturedMediaRepo.Save(NewMediaInput{
			Attributes:    mediaResult.Attributes,
			Height:        mediaResult.Height,
			Width:         mediaResult.Width,
			X:             mediaResult.X,
			Y:             mediaResult.Y,
			Url:           mediaResult.Url,
			PHash:         hash,
			Filename:      storageMedia.Filename,
			MediaUrl:      storageMedia.Media,
			ScreenshotUrl: storageMedia.Screenshot,
			ResourceUrl:   storageMedia.Resource,
			TaskId:        task.Id,
		}, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) GetMedia(mediaId string, ctx context.Context) (*models.Media, error) {
	media, err := p.capturedMediaRepo.GetMedia(mediaId, ctx)
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (p *Processor) GetMediaByHash(hash string, ctx context.Context) (*models.Media, error) {
	media, err := p.capturedMediaRepo.GetMediaByHash(hash, ctx)
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (p *Processor) GetMedias(filter *MediaFilter, ctx context.Context) ([]*models.Media, error) {
	medias, err := p.capturedMediaRepo.GetMedias(filter, ctx)
	if err != nil {
		return nil, err
	}

	return medias, nil
}
