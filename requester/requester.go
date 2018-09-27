package requester

import (
	"bufio"
	"net/http"
	"net/url"
	"os"
)

// ReadReq reads a request in from a file path, then returns the read request in
// a useable form,
func ReadReq(filepath string) (*http.Request, error) {
	fp, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	reader := bufio.NewReader(fp)
	req, err := http.ReadRequest(reader)
	if err != nil {
		return nil, err
	}

	var u *url.URL
	if req.Proto == "HTTP/1.0" || req.Proto == "HTTP/1.1" || req.Proto == "HTTP/2" {
		u, err = url.Parse("http://" + req.Host + req.RequestURI)
		if err != nil {
			return nil, err
		}
	} else {
		u, err = url.Parse("http://" + req.Host + req.RequestURI)
		if err != nil {
			return nil, err
		}
	}

	req.URL = u
	req.RequestURI = ""

	return req, nil
}
