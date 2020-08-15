package youtube

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strings"

	"github.com/horacehylee/go-youtube-dl/pkg/youtube/client"
	"github.com/horacehylee/go-youtube-dl/pkg/youtube/decipher"
	"golang.org/x/sync/errgroup"
)

const (
	streamChunkSize = int64(10 * 1024 * 1024) // 10MB for stream URL chuck
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

	length, err := d.client.StreamLength(url)
	if err != nil {
		return err
	}

	bytes, err := d.readStreamChunks(url, length)
	if err != nil {
		return err
	}
	bufWriter := bufio.NewWriterSize(w, int(streamChunkSize))
	for _, b := range bytes {
		bufWriter.Write(b)
	}
	return bufWriter.Flush()
}

func (d *Downloader) getURL(videoID string, s client.StreamFormat) (string, error) {
	if s.URL == "" && s.SignatureCipher == "" {
		return "", fmt.Errorf("Both url and signature cipher is empty")
	}
	if s.URL != "" {
		return s.URL, nil
	}
	cipher := decipher.NewDecipher(d.client)
	return cipher.StreamURL(videoID, s.SignatureCipher)
}

func (d *Downloader) readStreamChunks(url string, length int64) ([][]byte, error) {
	chunks := int(math.Ceil(float64(length) / float64(streamChunkSize)))
	bytes := make([][]byte, chunks)
	var g errgroup.Group

	from := int64(0)
	for i := 0; i < chunks; i++ {
		to := from + streamChunkSize
		if to > length {
			to = length
		}

		g.Go(func(i int, from int64, to int64) func() error {
			return func() error {
				r, err := d.client.Stream(url, from, to)
				if err != nil {
					return err
				}
				defer r.Close()

				b, err := ioutil.ReadAll(r)
				if err != nil {
					return err
				}
				bytes[i] = b
				return nil
			}
		}(i, from, to))

		from = to + 1
	}

	err := g.Wait()
	return bytes, err
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
