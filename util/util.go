package util

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

type scanByteCounter struct {
	BytesRead int64
}

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

func (s *scanByteCounter) wrap(split bufio.SplitFunc) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (int, []byte, error) {
		adv, tok, err := split(data, atEOF)
		s.BytesRead += int64(adv)
		return adv, tok, err
	}
}

// ReadInputs takes an array of file pointers and an array of int64, whic it uses
// as offsets into those files.  It then reads each file pointer, incrementing
// read heads only as needed.
func ReadInputs(inFiles []*os.File, offsets []int64) ([]string, bool, error) {
	ret := make([]string, len(inFiles))
	iter := true

	if len(inFiles) != len(offsets) {
		return nil, false, fmt.Errorf("ReadInputs: length of fp list != length of offsets")
	}

	for i := 0; i < len(inFiles); i++ {
		if inFiles[i] != nil {
			counter := scanByteCounter{}
			inFiles[i].Seek(offsets[i], 0)
			scanner := bufio.NewScanner(inFiles[i])
			splitFunc := counter.wrap(bufio.ScanLines)
			scanner.Split(splitFunc)

			ok := scanner.Scan()
			if ok {
				if iter {
					offsets[i] += counter.BytesRead
					iter = false
					if i != 0 {
						ok = scanner.Scan()
						if !ok { // next file needs to iterate
							if i == len(inFiles)-1 {
								return ret, true, nil
							}
							iter = true
							offsets[i] = 0
							counter.BytesRead = 0

							inFiles[i].Seek(offsets[i], 0)
							scanner = bufio.NewScanner(inFiles[i])
							scanner.Split(splitFunc)

							ok = scanner.Scan()
							if !ok {
								return nil, false, fmt.Errorf("ReadInputs: file %d empty?", i)
							}
						}
					}
					ret[i] = scanner.Text()
				} else {
					ret[i] = scanner.Text()
				}
			} else { // !ok
				if scanner.Err() != nil {
					return nil, false, scanner.Err()
				}
				// almost definitely EOF
				iter = true
				offsets[i] = 0
				counter.BytesRead = 0

				inFiles[i].Seek(offsets[i], 0)
				scanner = bufio.NewScanner(inFiles[i])
				scanner.Split(splitFunc)

				ok = scanner.Scan()
				if ok {
					offsets[i] += counter.BytesRead
					ret[i] = scanner.Text()
					if i == len(inFiles)-1 {
						return ret, true, nil
					}
				} else {
					return nil, false, fmt.Errorf("ReadInputs: file %d empty?", i)
				}
			}
		} else {
			ret[i] = ""
		}
	}

	return ret, false, nil
}
