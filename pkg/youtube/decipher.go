package youtube

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	urlKey = "url"
	spKey  = "sp"
	sigKey = "s"
)

func decryptCipher(videoID string, cipher string) (string, error) {
	p, err := url.ParseQuery(cipher)
	if err != nil {
		return "", err
	}
	if err := checkRequiredParam(p, urlKey); err != nil {
		return "", err
	}
	if err := checkRequiredParam(p, spKey); err != nil {
		return "", err
	}
	if err := checkRequiredParam(p, sigKey); err != nil {
		return "", err
	}

	sig, err := decryptSignature(videoID, p.Get(sigKey))
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%v&%v=%v", p.Get(urlKey), p.Get(spKey), sig)
	return url, nil
}

type decodeOp func(b []byte) []byte

type decodeOpFunc func(p interface{}) decodeOp

func reverseOpFunc(p interface{}) decodeOp {
	return func(b []byte) []byte {
		l, r := 0, len(b)-1
		for l < r {
			b[l], b[r] = b[r], b[l]
			l++
			r--
		}
		return b
	}
}

func spliceOpFunc(p interface{}) decodeOp {
	return func(b []byte) []byte {
		return b[p.(int):]
	}
}

func swapOpFunc(p interface{}) decodeOp {
	return func(b []byte) []byte {
		pos := p.(int) % len(b)
		b[0], b[pos] = b[pos], b[0]
		return b
	}
}

func decryptSignature(videoID string, sig string) (string, error) {
	playerURL, err := getPlayerURL(videoID)
	if err != nil {
		return "", err
	}
	fmt.Printf("player url: %v\n", playerURL)

	ops, err := getDecodeOps(playerURL)
	if err != nil {
		return "", err
	}

	b := []byte(sig)
	for _, op := range ops {
		b = op(b)
	}
	return string(b), nil
}

var (
	decodeOpsPattern = regexp.MustCompile(`function\(a\){a=a\.split\(""\);(.*);return a\.join\(""\)}`)
	reverseOpPattern = regexp.MustCompile(`([a-zA-Z_\\$][a-zA-Z_0-9]*):function\(a\){a\.reverse\(\)}`)
	spliceOpPattern  = regexp.MustCompile(`([a-zA-Z_\\$][a-zA-Z_0-9]*):function\(a,b\){a\.splice\(0,b\)}`)
	swapOpPattern    = regexp.MustCompile(`([a-zA-Z_\\$][a-zA-Z_0-9]*):function\(a,b\){var c=a\[0\];a\[0\]=a\[b%a\.length\];a\[b%a\.length\]=c}`)
)

func getDecodeOps(playerURL string) ([]decodeOp, error) {
	var ops []decodeOp
	resp, err := http.Get(playerURL)
	if err != nil {
		return ops, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ops, err
	}

	opsFuncMap := make(map[string]decodeOpFunc)
	if err := getOpFunc(opsFuncMap, reverseOpFunc, reverseOpPattern, b); err != nil {
		return ops, err
	}
	if err := getOpFunc(opsFuncMap, spliceOpFunc, spliceOpPattern, b); err != nil {
		return ops, err
	}
	if err := getOpFunc(opsFuncMap, swapOpFunc, swapOpPattern, b); err != nil {
		return ops, err
	}

	matches := decodeOpsPattern.FindSubmatch(b)
	if matches == nil || len(matches) < 2 {
		return ops, fmt.Errorf("failed to find ops with pattern: %v", decodeOpsPattern)
	}

	opsStrings := strings.Split(string(matches[1]), ";")
	if len(opsStrings) == 0 {
		return ops, fmt.Errorf("empty decode ops")
	}

	ops = make([]decodeOp, len(opsStrings))
	for i, opsString := range opsStrings {
		fmt.Printf("ops string: %v\n", opsString)
		f, err := parseJsFunction(opsString)
		if err != nil {
			return ops, err
		}
		opsFunc, ok := opsFuncMap[f.name]
		if !ok {
			return ops, fmt.Errorf("ops func cannot be found: %v", f.name)
		}
		ops[i] = opsFunc(f.param)
	}
	return ops, nil
}

func getOpFunc(ops map[string]decodeOpFunc, op decodeOpFunc, p *regexp.Regexp, b []byte) error {
	matches := p.FindSubmatch(b)
	if matches == nil || len(matches) < 2 {
		return fmt.Errorf("failed to find ops with pattern: %v", p.String())
	}
	k := string(matches[1])
	ops[k] = op
	return nil
}

var (
	playerURLPattern = regexp.MustCompile(`<script.*src="(.*)".*name="player_ias/base".*></script>`)
)

func getPlayerURL(videoID string) (string, error) {
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

	matches := playerURLPattern.FindSubmatch(b)
	if matches == nil || len(matches) < 2 {
		return "", fmt.Errorf("failed to find player url with pattern: %v", playerURLPattern)
	}
	return fmt.Sprintf("https://youtube.com%s", matches[1]), nil
}

func checkRequiredParam(v url.Values, k string) error {
	if v.Get(k) == "" {
		return fmt.Errorf("%v key cannot be found", k)
	}
	return nil
}
