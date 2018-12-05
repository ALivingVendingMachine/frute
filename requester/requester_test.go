package requester_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"

	r "github.com/ALivingVendingMachine/frute/requester"
)

func TestRequest(t *testing.T) {
	tests := []struct {
		filename string
	}{
		{"../testfiles/test_req"},
		{"../testfiles/test_file_2"},
	}

	for i, test := range tests {
		req, err := r.ReadReq(test.filename)

		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			t.Errorf("%d: err: %v\n", i, err)
		}

		tmpfile, err := ioutil.TempFile("", "req.*.out")
		if err != nil {
			t.Errorf("%d, error: %v\n", i, err)
		}
		defer os.Remove(tmpfile.Name())

		_, err = tmpfile.Write(dump)
		if err != nil {
			tmpfile.Close()
			t.Errorf("%d: error: %v\n", i, err)
		}

		testOldReq, err := ioutil.ReadFile(test.filename)
		if err != nil {
			t.Errorf("%d: error: %v\n", i, err)
		}
		testNewReq, err := ioutil.ReadFile(tmpfile.Name())
		if err != nil {
			t.Errorf("%d: error: %v\n", i, err)
		}

		if !bytes.Equal(testNewReq, testOldReq) {
			t.Errorf("read request and written request differ!\n")
		}

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("%d: error: %v\n", i, err)
		}

		dumpResp, err := httputil.DumpResponse(resp, true)
		if err != nil {
			t.Errorf("%d: error: %v\n", i, err)
		}
		t.Logf("%q\n", dumpResp)
	}
}
