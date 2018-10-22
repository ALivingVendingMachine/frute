package decoder

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/kothar/brotli-go.v0/dec"
)

func gunzip(res *http.Response) (string, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("gunzip: %v", err)
	}

	read, err := gzip.NewReader(bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("gunzip: %v", err)
	}
	defer read.Close()

	data, err := ioutil.ReadAll(read)
	if err != nil {
		return "", fmt.Errorf("gunzip: %v", err)
	}

	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return string(data), nil
}

func unbro(res *http.Response) (string, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("unbro: %v", err)
	}

	dec, err := dec.DecompressBuffer(body, nil)
	if err != nil {
		return "", fmt.Errorf("unbro: %v", err)
	}

	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return string(dec), nil
}

// Decode takes an http response, reads the Content-Encoding header, then
// decompresses the body if needed and returns it
func Decode(res *http.Response) (string, error) {
	enc := res.Header.Get("Content-Encoding")
	switch enc {
	case "gzip":
		return gunzip(res)
	case "br":
		return unbro(res)
	default:
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("Decode: %v", err)
		}
		return string(data), nil
	}
}
