package requester_test

import (
	"net/http"
	"net/http/httputil"
	"testing"

	r "github.com/alivingvendingmachine/frute/requester"
)

func TestRequest(t *testing.T) {
	file := "../testfiles/test_req"
	c := make(chan *http.Response)
	go func() {
		err := r.ReadReqAndDo(c, file)
		if err != nil {
			t.Errorf("err: %v", err)
			t.FailNow()
		}
	}()
	res := <-c
	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	t.Logf("%q\n", dump)
}

func TestRequest2(t *testing.T) {
	file := "../testfiles/test_file_2"
	c := make(chan *http.Response)
	go func() {
		err := r.ReadReqAndDo(c, file)
		if err != nil {
			t.Errorf("err: %v", err)
			t.FailNow()
		}
	}()
	res := <-c
	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	t.Logf("%q\n", dump)
}
