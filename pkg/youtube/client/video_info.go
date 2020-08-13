package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// VideoPlayerInfo provides player details
func (c *Client) VideoPlayerInfo(videoID string) (PlayerResponse, error) {
	v, err := c.videoInfo(videoID)
	if err != nil {
		return PlayerResponse{}, err
	}

	var p PlayerResponse
	if err := json.Unmarshal([]byte(v["player_response"][0]), &p); err != nil {
		return PlayerResponse{}, err
	}
	return p, nil
}

// videoInfo from youtube info API
func (c *Client) videoInfo(videoID string) (url.Values, error) {
	infoURL := fmt.Sprintf("https://youtube.com/get_video_info?video_id=%v&eurl=https://youtube.googleapis.com/v/%v", videoID, videoID)
	resp, err := c.client.Get(infoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get video info api returned: %v", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	info := string(body)

	v, err := url.ParseQuery(info)
	if err != nil {
		return nil, err
	}
	err = checkVideoInfoStatus(v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func checkVideoInfoStatus(v url.Values) error {
	status, ok := v["status"]
	if !ok {
		return fmt.Errorf("no response status found in the server's answer")
	}
	if status[0] == "fail" {
		reason, ok := v["reason"]
		if ok {
			return fmt.Errorf("'fail' response status found in the server's answer, reason: '%s'", reason[0])
		}
		return fmt.Errorf("'fail' response status found in the server's answer, no reason given")
	}
	if status[0] != "ok" {
		return fmt.Errorf("non-success response status found in the server's answer (status: '%s')", status)
	}
	return nil
}
