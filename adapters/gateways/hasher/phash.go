package hasher

import (
	"bytes"
	"fmt"
	"github.com/corona10/goimagehash"
	"image/png"
	//"golang.org/x/image/webp"
)

type PHashHandler struct{}

func NewPHashHandler() *PHashHandler {
	return &PHashHandler{}
}

func (p *PHashHandler) Hash(image []byte) (string, error) {
	decoded, err := png.Decode(bytes.NewReader(image))
	if err != nil {
		return "", fmt.Errorf("error decoding image: %w", err)
	}

	// Only return error if the image is not a valid image
	hash, _ := goimagehash.PerceptionHash(decoded)

	return hash.ToString(), nil
}
