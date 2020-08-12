package youtube

import (
	"fmt"
	"io"
	"strings"

	"github.com/horacehylee/go-youtube-dl/pkg/youtube/client"
	"github.com/horacehylee/go-youtube-dl/pkg/youtube/decipher"
)

// Download youtube video by video id
func Download(w io.Writer, videoID string) error {
	p, err := client.VideoPlayerInfo(videoID)
	if err != nil {
		return err
	}

	audioStreamFilter := func(s client.StreamFormat) bool {
		return strings.HasPrefix(s.MimeType, "audio/mp4")
	}
	audio, err := findStream(p.StreamingData.AdaptiveFormats, audioStreamFilter)
	if err != nil {
		return err
	}

	url, err := getURL(videoID, audio)
	if err != nil {
		return err
	}
	if len(url) == 0 {
		return fmt.Errorf("url cannot be empty")
	}
	fmt.Printf("url: %v\n", url)

	length, err := client.StreamLength(url)
	if err != nil {
		return err
	}

	r, err := client.Stream(url, 0, length)
	if err != nil {
		return err
	}
	defer r.Close()

	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	return nil
}

func getURL(videoID string, s client.StreamFormat) (string, error) {
	if s.URL == "" && s.SignatureCipher == "" {
		return "", fmt.Errorf("Both url and signature cipher is empty")
	}
	if s.URL != "" {
		return s.URL, nil
	}
	return decipher.DecryptStreamURL(videoID, s.SignatureCipher)
}

type streamPredicate = func(s client.StreamFormat) bool

func findStream(streams []client.StreamFormat, predicate streamPredicate) (client.StreamFormat, error) {
	var stream client.StreamFormat
	var found bool
	for _, s := range streams {
		if predicate(s) {
			stream = s
			found = true
			break
		}
	}
	if !found {
		return client.StreamFormat{}, fmt.Errorf("stream cannot be found")
	}
	return stream, nil
}
