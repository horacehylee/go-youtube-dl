package decipher

import (
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var mockClient = &youtubeClientMock{
	PlayerJSFunc: func(videoID string) (io.ReadCloser, error) {
		return os.Open("testdata/player.txt")
	},
}

func TestStreamUrl(t *testing.T) {
	d := NewDecipher(mockClient)

	videoID := "mKPbRIBCSHY"
	signatureCipher := "s=jjfPoalCVG4gcJ1loAz6Pkktg2tVIbs5zniAx2t9hOLEIC41rU53hAGLmIMZVxFVkYBoNrhFfhF4kixh4AWYM7QnCgIARw8JQ0qOdd&sp=sig&url=https://testing.com?test=1234"
	url, err := d.StreamURL(videoID, signatureCipher)
	if err != nil {
		t.Error(err)
	}

	expected := "https://testing.com?test=1234&sig=AOq0QJ8wRAIgCnQ7MYWA4hxik4FhfFhrNoBYkVFxVZMImLGdh35Ur14CIELOh9t2xAinz5sbIVt2gtkkP6zAol1Jcg4GVClaoPfj"
	if !cmp.Equal(expected, url) {
		t.Error(cmp.Diff(expected, url))
	}
}

func TestStreamUrlFailed(t *testing.T) {
	videoID := "mKPbRIBCSHY"
	signatureCiphers := []string{
		"sp=sig&url=https://testing.com?test=1234", // missing s
		"s=jjfPoalCVG4gcJ1loAz6Pkktg2tVIbs5zniAx2t9hOLEIC41rU53hAGLmIMZVxFVkYBoNrhFfhF4kixh4AWYM7QnCgIARw8JQ0qOdd&url=https://testing.com?test=1234",     // missing sp
		"s=jjfPoalCVG4gcJ1loAz6Pkktg2tVIbs5zniAx2t9hOLEIC41rU53hAGLmIMZVxFVkYBoNrhFfhF4kixh4AWYM7QnCgIARw8JQ0qOdd&sp=sig",                                // missing url
		"s=jjfPoalCVG4gcJ1loAz6Pkktg2tVIbs5zniAx2t9hOLEIC41rU53hAGLmIMZVxFVkYBoNrhFfhF4kixh4AWYM7QnCgIARw8JQ0qOdd&sp=sig&url=httpstesting.com?test=1234", // not url format
		"", // empty
		"s=ab%XXcd&sp=sig&url=https://testing.com?test=1234", // illegal format of URL query

		// "s=jjfPoalCVG4gcJ1loAz6Pkktg2tVIbs5zniAx2t9hOLEIC41rU53hAGLmIMZVxFVkYBoNrhFfhF4kixh4AWYM7QnCgIARw8JQ0qOdd&sp=sig&url=https://testing.com?test=1234", // missing s
	}
	d := NewDecipher(mockClient)

	for _, s := range signatureCiphers {
		_, err := d.StreamURL(videoID, s)
		if err == nil {
			t.Error("error not thrown for signatureCipher", s)
		}
	}
}
