package util_test

import (
	"os"
	"testing"

	u "github.com/alivingvendingmachine/frute/util"
)

func TestGenerateRequest(t *testing.T) {
	err := u.GenerateRequest("GET", "http://www.google.com", "", nil, "test.out")
	if err != nil {
		t.Errorf("error: %v", err)
	}
	err = os.Remove("test.out")
	if err != nil {
		t.Errorf("error: %v", err)
	}
}

func TestGenerateRequestErrors(t *testing.T) {
	err := u.GenerateRequest("NO", "http://www.google.com", "", nil, "test.out")
	if err == nil {
		t.Error("expected error, got nothing")
	}

	err = u.GenerateRequest("post", "www.google.com", "", nil, "test.out")
	if err == nil {
		t.Error("expected error, got nothing")
	}
}

func TestBulkOpen(t *testing.T) {
	files := []string{
		"../testfiles/testFile0",
	}
	fps, err := u.BulkOpen(files)
	if err != nil {
		t.Errorf("error: %v\n", err)
	}
	for i, filename := range files {
		if fps[i] == nil {
			t.Errorf("file %d nil", i)
		}
		if fps[i].Name() != filename {
			t.Errorf("file %d: opened %s, expected %s", i, fps[i].Name(), filename)
		}
		t.Logf("file %d name: %s\n", i, fps[i].Name())
	}
}

func TestBulkOpenFull(t *testing.T) {
	files := []string{
		"../testfiles/testFile0",
		"../testfiles/testFile1",
		"../testfiles/testFile2",
		"../testfiles/testFile3",
		"../testfiles/testFile4",
		"../testfiles/testFile5",
		"../testfiles/testFile6",
		"../testfiles/testFile7",
		"../testfiles/testFile8",
		"../testfiles/testFile9",
	}
	fps, err := u.BulkOpen(files)
	if err != nil {
		t.Errorf("error: %v\n", err)
	}
	for i, filename := range files {
		if fps[i] == nil {
			t.Errorf("file %d nil", i)
		}
		if fps[i].Name() != filename {
			t.Errorf("file %d: opened %s, expected %s", i, fps[i].Name(), filename)
		}
		t.Logf("file %d name: %s\n", i, fps[i].Name())
	}
}
