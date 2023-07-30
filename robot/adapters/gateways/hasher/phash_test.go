package hasher

import (
	"bytes"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestPHashHandler(t *testing.T) {
	logger := zap.NewExample()
	imageFile, err := os.Open("../../../testing_resources/test.png")
	if err != nil {
		t.Fatalf("error opening image: %v", err)
	}
	defer func(imageFile *os.File) {
		err := imageFile.Close()
		if err != nil {
			t.Fatalf("error closing image: %v", err)
		}
	}(imageFile)

	// Image to bytes
	imageBuff := bytes.Buffer{}
	_, err = imageBuff.ReadFrom(imageFile)
	if err != nil {
		t.Fatalf("error reading image: %v", err)
	}
	base64Image := imageBuff.Bytes()

	tests := []struct {
		name    string
		image   []byte
		wantErr bool
		want    string
	}{
		{
			name:    "should return error when image is invalid",
			image:   []byte("invalid"),
			wantErr: true,
		},
		{
			name:    "should return error when image is nil",
			image:   nil,
			wantErr: true,
		},
		{
			name:    "should return hash when image is valid",
			image:   base64Image,
			wantErr: false,
			want:    "p:8f28f6d738680f07",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewPHashHandler(logger)
			got, err := handler.Hash(tt.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("PHashHandler.Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == "" && !tt.wantErr {
				t.Errorf("PHashHandler.Hash() = %v, want %v", got, tt.want)
			}
			if tt.want != "" && got != tt.want {
				t.Errorf("PHashHandler.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
