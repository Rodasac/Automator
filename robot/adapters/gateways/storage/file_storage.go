package storage

import (
	"automator-go/robot/usecases/task"
	"fmt"
	"github.com/nlepage/go-cuid2"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"os"
)

type FileStorage struct {
	MediaExtension string
	logger         *otelzap.LoggerWithCtx
}

func NewFileStorage(extension string, logger *otelzap.LoggerWithCtx) *FileStorage {
	if extension == "" {
		extension = "png"
	}

	return &FileStorage{
		MediaExtension: extension,
		logger:         logger,
	}
}

func (fsm *FileStorage) SaveMedia(hashWithoutKind string, media *task.RawMedia) (task.StorageMedia, error) {
	fsm.logger.Debug("Saving media files")
	filenameId, err := cuid2.CreateId()
	if err != nil {
		return task.StorageMedia{}, fmt.Errorf("error generating files id: %w", err)
	}

	mediaFilename := hashWithoutKind + "_" + filenameId + "." + fsm.MediaExtension
	screenshotFilename := hashWithoutKind + "_" + filenameId + "." + fsm.MediaExtension
	mediaPath := "./media/media_" + mediaFilename
	screenshotPath := "./media/screenshot_" + screenshotFilename

	fileMedia, err := os.Create(mediaPath)
	if err != nil {
		return task.StorageMedia{}, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			println(err.Error())
		}
	}(fileMedia)

	_, err = fileMedia.Write(media.Media)
	if err != nil {
		return task.StorageMedia{}, err
	}
	fsm.logger.Debug("Saved media to file storage")

	fileScreenshot, err := os.Create(screenshotPath)
	if err != nil {
		return task.StorageMedia{}, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			println(err.Error())
		}
	}(fileScreenshot)

	_, err = fileScreenshot.Write(media.Screenshot)
	if err != nil {
		return task.StorageMedia{}, err
	}
	fsm.logger.Debug("Saved screenshot to file storage")

	var resourcePath string
	if media.Resource != nil && len(media.Resource) > 0 {
		resourceFilename := hashWithoutKind + "_" + filenameId + "." + media.Ext
		resourcePath = "./media/resource_" + resourceFilename

		fileResource, err := os.Create(resourcePath)
		if err != nil {
			return task.StorageMedia{}, err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				println(err.Error())
			}
		}(fileResource)

		_, err = fileResource.Write(media.Resource)
		if err != nil {
			return task.StorageMedia{}, err
		}
		fsm.logger.Debug("Saved resource to file storage")
	}

	return task.StorageMedia{
		Filename:   mediaFilename,
		Media:      mediaPath,
		Screenshot: screenshotPath,
		Resource:   resourcePath,
	}, nil
}
