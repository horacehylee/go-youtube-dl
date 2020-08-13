package decipher

import (
	"fmt"
	"net/url"
)

type decodedSignatureCipher struct {
	URL string `validate:"required,url"`
	Sp  string `validate:"required"`
	Sig string `validate:"required"`
}

//StreamURL is used for getting youtube url from signature cipher
func (d *Decipher) StreamURL(videoID string, signatureCipher string) (string, error) {
	decoded, err := d.decodeSignatureCipher(signatureCipher)
	if err != nil {
		return "", err
	}

	sig, err := d.decryptSignature(videoID, decoded.Sig)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%v&%v=%v", decoded.URL, decoded.Sp, sig)
	return url, nil
}

func (d *Decipher) decodeSignatureCipher(signatureCipher string) (decodedSignatureCipher, error) {
	v, err := url.ParseQuery(signatureCipher)
	if err != nil {
		return decodedSignatureCipher{}, err
	}
	decoded := decodedSignatureCipher{
		URL: v.Get("url"),
		Sp:  v.Get("sp"),
		Sig: v.Get("s"),
	}
	return decoded, d.validate.Struct(decoded)
}
