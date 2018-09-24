package util

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

// GenerateRequest takes
func GenerateRequest(method string, url string, body string, filename string) error {
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
