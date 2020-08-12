package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	playerJsURLPattern = regexp.MustCompile(`<script.*src="(.*)".*name="player_ias/base".*></script>`)
)

// PlayerJS to fetch Javascript of player for video
func PlayerJS(videoID string) (io.ReadCloser, error) {
	playerURL, err := playerJsURL(videoID)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(playerURL)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func playerJsURL(videoID string) (string, error) {
	url := fmt.Sprintf("https://youtube.com/embed/%v?hl=en", videoID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	matches := playerJsURLPattern.FindSubmatch(b)
	if matches == nil || len(matches) < 2 {
		return "", fmt.Errorf("failed to find player url with pattern: %v", playerJsURLPattern)
	}
	return fmt.Sprintf("https://youtube.com%s", matches[1]), nil
}
