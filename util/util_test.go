package util_test

import (
	"os"
	"testing"

	u "github.com/alivingvendingmachine/frute/util"
)

func TestGenerateRequest(t *testing.T) {
	err := u.GenerateRequest("GET", "http://www.google.com", "", "test.out")
	if err != nil {
		t.Errorf("error: %v", err)
	}
	err = os.Remove("test.out")
	if err != nil {
		t.Errorf("error: %v", err)
	}
}

func TestGenerateRequestErrors(t *testing.T) {
	err := u.GenerateRequest("NO", "http://www.google.com", "", "test.out")
	if err == nil {
		t.Error("expected error, got nothing")
	}

	err = u.GenerateRequest("post", "www.google.com", "", "test.out")
	if err == nil {
		t.Error("expected error, got nothing")
	}
}
