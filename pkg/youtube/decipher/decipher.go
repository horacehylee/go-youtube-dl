package decipher

import (
	"github.com/go-playground/validator/v10"
	"github.com/horacehylee/go-youtube-dl/pkg/youtube/client"
)

// Decipher for Youtube API
type Decipher struct {
	validate *validator.Validate
	client   *client.Client
}

// NewDecipher returns new instance of decipher
func NewDecipher(c *client.Client) *Decipher {
	return &Decipher{
		validate: validator.New(),
		client:   c,
	}
}
