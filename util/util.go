package util

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

// GenerateRequest takes a method, url, body, a list of headers, and a filepath,
// all strings.  It then generates a request for the given url and method, with
// the given body, and writes it out to the given filepath.
func GenerateRequest(method string, url string, body string, headerFlags []string, filename string) error {
	method = strings.ToUpper(method)
	if method != "GET" &&
		method != "HEAD" &&
		method != "POST" &&
		method != "PUT" &&
		method != "DELETE" &&
		method != "OPTIONS" &&
		method != "TRACE" &&
		method != "PATCH" {
		return errors.New("generate request: invalid method")
	}

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return err
	}

	for _, header := range headerFlags {
		h := strings.Split(header, ":")
		if h[1][0] == ' ' {
			req.Header.Add(h[0], h[1][1:])
		} else {
			req.Header.Add(h[0], h[1])
		}
	}

	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(dump)
	if err != nil {
		return err
	}

	return nil
}

// BulkOpen takes a slice of strings of filenames, and returns a slice of
// file handlers to the caller.
func BulkOpen(infiles []string) ([]*os.File, error) {
	ret := make([]*os.File, len(infiles))

	for i := 0; i < len(infiles); i++ {
		if infiles[i] != "" {
			fp, err := os.Open(infiles[i])
			if err != nil {
				return nil, fmt.Errorf("error opening file %d", i)
			}
			ret[i] = fp
		} else {
			ret[i] = nil
		}
	}

	return ret, nil
}
