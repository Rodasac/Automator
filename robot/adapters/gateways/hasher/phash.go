package hasher

import (
	"bytes"
	"fmt"
	"github.com/corona10/goimagehash"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"image/png"
)

type PHashHandler struct {
	logger *otelzap.LoggerWithCtx
}

func NewPHashHandler(logger *otelzap.LoggerWithCtx) *PHashHandler {
	return &PHashHandler{
		logger: logger,
	}
}

func (p *PHashHandler) Hash(image []byte) (string, error) {
	p.logger.Debug("Hashing image")
	decoded, err := png.Decode(bytes.NewReader(image))
	if err != nil {
		return "", fmt.Errorf("error decoding image: %w", err)
	}

	// Only return error if the image is not a valid image
	hash, _ := goimagehash.PerceptionHash(decoded)
	p.logger.Debug("Finished hashing image")

	return hash.ToString(), nil
}
