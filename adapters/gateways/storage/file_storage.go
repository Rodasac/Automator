package storage

import (
	"automator-go/usecases/task"
	"fmt"
	"github.com/nlepage/go-cuid2"
	"os"
)

type FileStorage struct {
	ext string
}

func NewFileStorage() *FileStorage {
	return &FileStorage{
		ext: ".png",
	}
}

func (fsm *FileStorage) SaveMedia(hashWithoutKind string, media []byte, screenshot []byte) (task.StorageMedia, error) {
	filenameId, err := cuid2.CreateId()
	if err != nil {
		return task.StorageMedia{}, fmt.Errorf("error generating files id: %w", err)
	}

	mediaFilename := hashWithoutKind + "_" + filenameId + fsm.ext
	screenshotFilename := hashWithoutKind + "_" + filenameId + fsm.ext
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

	_, err = fileMedia.Write(media)
	if err != nil {
		return task.StorageMedia{}, err
	}

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

	_, err = fileScreenshot.Write(screenshot)
	if err != nil {
		return task.StorageMedia{}, err
	}

	return task.StorageMedia{
		Filename:   mediaFilename,
		Media:      mediaPath,
		Screenshot: screenshotPath,
	}, nil
}
