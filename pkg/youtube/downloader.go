package youtube

import (
	"fmt"
	"io"
	"strings"

	"github.com/horacehylee/go-youtube-dl/pkg/youtube/client"
	"github.com/horacehylee/go-youtube-dl/pkg/youtube/decipher"
)

// Downloader for downloading Youtube videos
type Downloader struct {
	client *client.Client
}

// NewDownloader returns new instance of downloader
func NewDownloader() *Downloader {
	return &Downloader{
		client: client.NewClient(),
	}
}

// DownloadVideo of Youtube video with video ID
func (d *Downloader) DownloadVideo(videoID string, w io.Writer) error {
	return nil
}

// DownloadAudio of Youtube video with video ID
func (d *Downloader) DownloadAudio(videoID string, w io.Writer) error {
	p, err := d.client.VideoPlayerInfo(videoID)
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

	url, err := d.getURL(videoID, audio)
	if err != nil {
		return err
	}
	fmt.Printf("url: %v\n", url)

	length, err := d.client.StreamLength(url)
	if err != nil {
		return err
	}

	r, err := d.client.Stream(url, 0, length)
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = io.Copy(w, r)
	return err
}

func (d *Downloader) getURL(videoID string, s client.StreamFormat) (string, error) {
	if s.URL == "" && s.SignatureCipher == "" {
		return "", fmt.Errorf("Both url and signature cipher is empty")
	}
	if s.URL != "" {
		return s.URL, nil
	}
	cipher := decipher.NewDecipher(d.client)
	url, err := cipher.StreamURL(videoID, s.SignatureCipher)
	if err != nil {
		return "", err
	}
	if url == "" {
		return "", fmt.Errorf("url cannot be empty")
	}
	return url, nil
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
