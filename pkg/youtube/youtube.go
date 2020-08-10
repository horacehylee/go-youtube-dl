package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Client for downloading youtube videos
type Client struct {
}

// Download youtube video by video id
func (c Client) Download(w io.Writer, id string) error {
	v, err := getVideoInfo(id)
	if err != nil {
		return err
	}

	// fmt.Printf("%v\n", v["player_response"][0])

	var p PlayerResponse
	if err := json.Unmarshal([]byte(v["player_response"][0]), &p); err != nil {
		return err
	}

	// b, err := json.MarshalIndent(p, "", "  ")
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(b))

	// s := make([]StreamFormat, len(p.StreamingData.Formats)+len(p.StreamingData.AdaptiveFormats))
	var audio StreamFormat
	var found bool
	for _, s := range p.StreamingData.AdaptiveFormats {
		if strings.HasPrefix(s.MimeType, "audio/mp4") {
			audio = s.StreamFormat
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("audio format cannot be found")
	}
	if len(audio.URL) == 0 {
		return fmt.Errorf("audio url cannot be empty")
	}
	log.Printf("audio url: %v\n", audio.URL)

	length, err := getContentLength(audio.URL)
	if err != nil {
		return err
	}

	resp, err := getStream(audio.URL, length)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(w, resp.Body); err != nil {
		return err
	}
	return nil
}

func getVideoInfo(id string) (url.Values, error) {
	infoURL := fmt.Sprintf("https://youtube.com/get_video_info?video_id=%v&eurl=https://youtube.googleapis.com/v/%v", id, id)
	resp, err := http.Get(infoURL)
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

	p, err := url.ParseQuery(info)
	if err != nil {
		return nil, err
	}
	err = checkStatus(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func checkStatus(v url.Values) error {
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

func getContentLength(url string) (int64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get content length, status code returned: %v", resp.StatusCode)
	}
	return resp.ContentLength, nil
}

func getStream(url string, length int64) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// TODO: can split length into multiple goroutine to speed up http calls
	req.Header.Set("range", fmt.Sprintf("bytes=0-%v", length))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusPartialContent {
		return nil, fmt.Errorf("failed to get stream, status code returned: %v", resp.StatusCode)
	}
	return resp, nil
}
