package task

import (
	models2 "automator-go/robot/entities/models"
	"automator-go/robot/usecases/hasher"
	"context"
	"errors"
	"testing"
)

type MockAutomatorTaskAdapter struct {
	Media *RawMedia
	Error error
}

func (m *MockAutomatorTaskAdapter) Run(*models2.Task) (*[]RawMedia, error) {
	if m.Error != nil || m.Media == nil {
		return nil, m.Error
	}

	return &[]RawMedia{*m.Media}, m.Error
}

type MockStorageMediaAdapter struct {
	Error error
}

func (m *MockStorageMediaAdapter) SaveMedia(string, *RawMedia) (StorageMedia, error) {
	return StorageMedia{}, m.Error
}

type MockCapturedMediaRepository struct {
	Error error
}

func (m *MockCapturedMediaRepository) Save(NewMediaInput, context.Context) error {
	return m.Error
}

func (m *MockCapturedMediaRepository) GetMedia(string, context.Context) (*models2.Media, error) {
	return &models2.Media{}, m.Error
}

func (m *MockCapturedMediaRepository) GetMediaByHash(string, context.Context) (*models2.Media, error) {
	return &models2.Media{}, m.Error
}

func (m *MockCapturedMediaRepository) GetMedias(*MediaFilter, context.Context) ([]*models2.Media, error) {
	return []*models2.Media{}, m.Error
}

type MockImageHasher struct {
	Error error
}

func (m *MockImageHasher) Hash([]byte) (string, error) {
	return "h:filename", m.Error
}

func TestProcessor(t *testing.T) {
	task := &models2.Task{
		Id:          "1",
		Title:       "Test",
		Description: "Test",
		Url:         "https://google.com",
		Country:     "US",
		WithProxy:   false,
		Actions: []models2.TaskAction{
			{
				Id:    "1",
				Label: "Navigate",
				Type:  models2.Navigate,
				Value: "input[name='q']",
			},
		},
	}
	media := &RawMedia{
		Ext:        "png",
		Media:      []byte("test"),
		Screenshot: []byte("test"),
		Attributes: map[string]interface{}{
			"test": "test",
		},
		Height: 100,
		Width:  100,
		X:      100,
		Y:      100,
		Url:    "https://google.com",
	}

	tests := []struct {
		name                 string
		automatorTaskAdapter AutomatorTaskAdapter
		capturedMediaRepo    CapturedMediaRepository
		storageMediaAdapter  StorageMediaAdapter
		imageHasher          hasher.ImageHasher
		task                 *models2.Task
		wantErr              bool
	}{
		{
			name: "success with media",
			automatorTaskAdapter: &MockAutomatorTaskAdapter{
				Media: media,
			},
			capturedMediaRepo:   &MockCapturedMediaRepository{},
			storageMediaAdapter: &MockStorageMediaAdapter{},
			imageHasher:         &MockImageHasher{},
			task:                task,
			wantErr:             false,
		},
		{
			name:                 "success without media",
			automatorTaskAdapter: &MockAutomatorTaskAdapter{},
			capturedMediaRepo:    &MockCapturedMediaRepository{},
			storageMediaAdapter:  &MockStorageMediaAdapter{},
			imageHasher:          &MockImageHasher{},
			task:                 task,
			wantErr:              false,
		},
		{
			name: "error automator task adapter",
			automatorTaskAdapter: &MockAutomatorTaskAdapter{
				Error: errors.New("error"),
			},
			capturedMediaRepo:   &MockCapturedMediaRepository{},
			storageMediaAdapter: &MockStorageMediaAdapter{},
			imageHasher:         &MockImageHasher{},
			task:                task,
			wantErr:             true,
		},
		{
			name: "error captured media repo",
			automatorTaskAdapter: &MockAutomatorTaskAdapter{
				Media: media},
			capturedMediaRepo: &MockCapturedMediaRepository{
				Error: errors.New("error"),
			},
			storageMediaAdapter: &MockStorageMediaAdapter{},
			imageHasher:         &MockImageHasher{},
			task:                task,
			wantErr:             true,
		},
		{
			name: "error storage media adapter",
			automatorTaskAdapter: &MockAutomatorTaskAdapter{
				Media: media,
			},
			capturedMediaRepo: &MockCapturedMediaRepository{},
			storageMediaAdapter: &MockStorageMediaAdapter{
				Error: errors.New("error"),
			},
			imageHasher: &MockImageHasher{},
			task:        task,
			wantErr:     true,
		},
		{
			name: "error image hasher",
			automatorTaskAdapter: &MockAutomatorTaskAdapter{
				Media: media,
			},
			capturedMediaRepo:   &MockCapturedMediaRepository{},
			storageMediaAdapter: &MockStorageMediaAdapter{},
			imageHasher: &MockImageHasher{
				Error: errors.New("error"),
			},
			task:    task,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewProcessor(tt.automatorTaskAdapter, tt.capturedMediaRepo, tt.storageMediaAdapter, tt.imageHasher)
			err := processor.Process(tt.task, context.TODO())
			if (err != nil) != tt.wantErr {
				t.Errorf("Processor.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
