package decipher

import (
	"io"

	"github.com/go-playground/validator/v10"
)

//go:generate moq -out decipher_test.go . youtubeClient
type youtubeClient interface {
	PlayerJS(videoID string) (io.ReadCloser, error)
}

// Decipher for Youtube API
type Decipher struct {
	validate *validator.Validate
	client   youtubeClient
}

// New returns new instance of decipher
func New(c youtubeClient) *Decipher {
	return &Decipher{
		validate: validator.New(),
		client:   c,
	}
}
