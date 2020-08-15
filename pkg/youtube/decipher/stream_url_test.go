package decipher

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStreamUrl(t *testing.T) {
	mockClient := &youtubeClientMock{
		PlayerJSFunc: func(videoID string) (io.ReadCloser, error) {
			return os.Open("testdata/player.txt")
		},
	}
	d := New(mockClient)

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
	mockClient := &youtubeClientMock{
		PlayerJSFunc: func(videoID string) (io.ReadCloser, error) {
			return os.Open("testdata/player.txt")
		},
	}
	d := New(mockClient)

	videoID := "mKPbRIBCSHY"
	signatureCiphers := []string{
		"sp=sig&url=https://testing.com?test=1234", // missing s
		"s=jjfPoalCVG4gcJ1loAz6Pkktg2tVIbs5zniAx2t9hOLEIC41rU53hAGLmIMZVxFVkYBoNrhFfhF4kixh4AWYM7QnCgIARw8JQ0qOdd&url=https://testing.com?test=1234",     // missing sp
		"s=jjfPoalCVG4gcJ1loAz6Pkktg2tVIbs5zniAx2t9hOLEIC41rU53hAGLmIMZVxFVkYBoNrhFfhF4kixh4AWYM7QnCgIARw8JQ0qOdd&sp=sig",                                // missing url
		"s=jjfPoalCVG4gcJ1loAz6Pkktg2tVIbs5zniAx2t9hOLEIC41rU53hAGLmIMZVxFVkYBoNrhFfhF4kixh4AWYM7QnCgIARw8JQ0qOdd&sp=sig&url=httpstesting.com?test=1234", // not url format
		"s=ab%XXcd&sp=sig&url=https://testing.com?test=1234", // illegal format of URL query
		"", // empty
	}

	for _, s := range signatureCiphers {
		_, err := d.StreamURL(videoID, s)
		if err == nil {
			t.Error("error not thrown for signatureCipher", s)
		}
	}
}

func TestFailedFetchPlayerJs(t *testing.T) {
	type testCase struct {
		name          string
		decipher      *Decipher
		expectedError error
	}
	cases := []testCase{
		{
			name: "failed fetching player js",
			decipher: New(&youtubeClientMock{
				PlayerJSFunc: func(videoID string) (io.ReadCloser, error) {
					return nil, errors.New("testing error")
				},
			}),
			expectedError: errors.New("testing error"),
		},
		{
			name: "failed to read player js",
			decipher: New(&youtubeClientMock{
				PlayerJSFunc: func(videoID string) (io.ReadCloser, error) {
					r := &mockReadCloser{
						ReadFunc: func(p []byte) (n int, err error) {
							return 0, errors.New("testing failed to read")
						},
						CloseFunc: func() error {
							return nil
						},
					}
					return r, nil
				},
			}),
			expectedError: errors.New("testing failed to read"),
		},
		{
			name: "failed to read empty player js",
			decipher: New(&youtubeClientMock{
				PlayerJSFunc: func(videoID string) (io.ReadCloser, error) {
					return getMockReadCloser([]byte("")), nil
				},
			}),
			expectedError: errors.New("failed to find reverse decrypt function"),
		},
		{
			name: "failed to read player js with empty decrypt ops",
			decipher: New(&youtubeClientMock{
				PlayerJSFunc: func(videoID string) (io.ReadCloser, error) {
					return os.Open("testdata/empty_decrypt_ops.txt")
				},
			}),
			expectedError: errors.New("empty decrypt ops"),
		},
		{
			name: "failed to read player js with no decrypt ops",
			decipher: New(&youtubeClientMock{
				PlayerJSFunc: func(videoID string) (io.ReadCloser, error) {
					return os.Open("testdata/no_decrypt_ops.txt")
				},
			}),
			expectedError: errors.New("failed to find decrypt ops"),
		},
		{
			name: "failed to read player js with invalid decrypt function",
			decipher: New(&youtubeClientMock{
				PlayerJSFunc: func(videoID string) (io.ReadCloser, error) {
					return os.Open("testdata/invalid_decrypt_function.txt")
				},
			}),
			expectedError: errors.New("failed to split JS function name"),
		},
		{
			name: "failed to read player js with unknown decrypt op",
			decipher: New(&youtubeClientMock{
				PlayerJSFunc: func(videoID string) (io.ReadCloser, error) {
					return os.Open("testdata/unknown_decrypt_ops.txt")
				},
			}),
			expectedError: errors.New("ops function cannot be found: UnkownFunction"),
		},
	}

	for _, c := range cases {
		videoID := "mKPbRIBCSHY"
		signatureCipher := "s=jjfPoalCVG4gcJ1loAz6Pkktg2tVIbs5zniAx2t9hOLEIC41rU53hAGLmIMZVxFVkYBoNrhFfhF4kixh4AWYM7QnCgIARw8JQ0qOdd&sp=sig&url=https://testing.com?test=1234"
		_, err := c.decipher.StreamURL(videoID, signatureCipher)
		if err == nil {
			t.Error("failed case: ", c.name)
		} else {
			if c.expectedError != nil {
				if !cmp.Equal(c.expectedError.Error(), err.Error()) {
					t.Error(cmp.Diff(c.expectedError.Error(), err.Error()))
				}
			}
		}
	}
}

type mockReadCloser struct {
	ReadFunc  func(p []byte) (n int, err error)
	CloseFunc func() error
}

func (r *mockReadCloser) Read(p []byte) (n int, err error) {
	return r.ReadFunc(p)
}

func (r *mockReadCloser) Close() error {
	return r.CloseFunc()
}

func getMockReadCloser(data []byte) *mockReadCloser {
	idx := int64(0)
	return &mockReadCloser{
		ReadFunc: func(p []byte) (n int, err error) {
			if idx >= int64(len(data)) {
				err = io.EOF
				return
			}

			n = copy(p, data[idx:])
			idx += int64(n)
			return
		},
		CloseFunc: func() error {
			return nil
		},
	}
}
