package task

import (
	"automator-go/entities/models"
)

type RawMedia struct {
	Media      []byte
	Screenshot []byte
	Attributes map[string]interface{}
	Height     float64
	Width      float64
	X          float64
	Y          float64
	Url        string
}

type AutomatorTaskAdapter interface {
	Run(task *models.Task) (*[]RawMedia, error)
}

type StorageMedia struct {
	Filename   string
	Media      string
	Screenshot string
}

type StorageMediaAdapter interface {
	SaveMedia(hashName string, media []byte, screenshot []byte) (StorageMedia, error)
}

type NewMediaInput struct {
	Attributes    map[string]interface{}
	Height        float64
	Width         float64
	X             float64
	Y             float64
	Url           string
	PHash         string
	Filename      string
	MediaUrl      string
	ScreenshotUrl string
	TaskId        string
}

type CapturedMediaRepository interface {
	Save(input NewMediaInput) error
}

type ProcessorUseCase interface {
	Process(task *models.Task) error
}
