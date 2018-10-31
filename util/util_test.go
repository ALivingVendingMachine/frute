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
		t.Fatalf("error: %v\n", err)
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

func TestFirstHalfOpen(t *testing.T) {
	files := []string{
		"../testfiles/testFile0",
		"../testfiles/testFile1",
		"../testfiles/testFile2",
		"../testfiles/testFile3",
		"../testfiles/testFile4",
		"",
		"",
		"",
		"",
		"",
	}

	fps, err := u.BulkOpen(files)
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
	for i, filename := range files {
		if fps[i] == nil && files[i] != "" {
			t.Errorf("file %d nil", i)
		}
		if fps[i] != nil {
			t.Logf("file %d fp name %s\tfilename %s", i, fps[i].Name(), filename)
		} else {
			t.Logf("file %d fp name: nil\tfilename %s", i, filename)
		}
	}
}

func TestSecondHalfOpen(t *testing.T) {
	files := []string{
		"",
		"",
		"",
		"",
		"",
		"../testfiles/testFile5",
		"../testfiles/testFile6",
		"../testfiles/testFile7",
		"../testfiles/testFile8",
		"../testfiles/testFile9",
	}

	fps, err := u.BulkOpen(files)
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
	for i, filename := range files {
		if fps[i] == nil && files[i] != "" {
			t.Errorf("file %d nil", i)
		}
		if fps[i] != nil {
			t.Logf("file %d fp name %s\tfilename %s", i, fps[i].Name(), filename)
		} else {
			t.Logf("file %d fp name: nil\tfilename %s", i, filename)
		}
	}
}

func TestReadInputsNoIterNext(t *testing.T) {
	files := []string{
		"../testfiles/testFile0",
		"../testfiles/testFile1",
		"../testfiles/testFile2",
	}

	offs := make([]int64, len(files))

	fps, err := u.BulkOpen(files)
	if err != nil {
		t.Fatalf("error in BulkOpen: %v\n", err)
	}

	expected := [][]string{
		[]string{"tf0-0", "tf1-0", "tf2-0"},
		[]string{"tf0-1", "tf1-0", "tf2-0"},
		[]string{"tf0-2", "tf1-0", "tf2-0"},
		[]string{"tf0-3", "tf1-0", "tf2-0"},
		[]string{"tf0-4", "tf1-0", "tf2-0"},
	}

	for i := 0; i < 5; i++ {
		out, _, err := u.ReadInputs(fps, offs)
		if err != nil {
			t.Fatalf("error in ReadInputs: %v\n", err)
		}
		for j := 0; j < len(out); j++ {
			t.Logf("%d: %s", i, out[j])
		}
		for j := 0; j < len(out); j++ {
			if out[j] != expected[i][j] {
				t.Errorf("%d failed, out[%d] = %s, expected[%d][%d] = %s", i, j, out[j], i, j, expected[i][j])
			}
		}
	}
}

func TestReadInputsIterNextSmall(t *testing.T) {
	files := []string{
		"../testfiles/small",
		"../testfiles/testFile1",
		"../testfiles/testFile2",
	}

	offs := make([]int64, len(files))

	fps, err := u.BulkOpen(files)
	if err != nil {
		t.Fatalf("error in BulkOpen: %v\n", err)
	}

	expected := [][]string{
		[]string{"same", "tf1-0", "tf2-0"},
		[]string{"same", "tf1-1", "tf2-0"},
	}

	for i := 0; i < 2; i++ {
		out, _, err := u.ReadInputs(fps, offs)
		if err != nil {
			t.Fatalf("error in ReadInputs: %v\n", err)
		}
		for j := 0; j < len(out); j++ {
			t.Logf("%d: %s", i, out[j])
		}
		for j := 0; j < len(out); j++ {
			if out[j] != expected[i][j] {
				t.Errorf("%d failed, out[%d] = %s, expected[%d][%d] = %s", i, j, out[j], i, j, expected[i][j])
			}
		}
	}
}

func TestReadInputsIterNext(t *testing.T) {
	files := []string{
		"../testfiles/testFile0",
		"../testfiles/testFile1",
		"../testfiles/testFile2",
	}

	offs := make([]int64, len(files))

	fps, err := u.BulkOpen(files)
	if err != nil {
		t.Fatalf("error in BulkOpen: %v\n", err)
	}

	expected := [][]string{
		[]string{"tf0-0", "tf1-0", "tf2-0"},
		[]string{"tf0-1", "tf1-0", "tf2-0"},
		[]string{"tf0-2", "tf1-0", "tf2-0"},
		[]string{"tf0-3", "tf1-0", "tf2-0"},
		[]string{"tf0-4", "tf1-0", "tf2-0"},
		[]string{"tf0-0", "tf1-1", "tf2-0"},
		[]string{"tf0-1", "tf1-1", "tf2-0"},
		[]string{"tf0-2", "tf1-1", "tf2-0"},
		[]string{"tf0-3", "tf1-1", "tf2-0"},
		[]string{"tf0-4", "tf1-1", "tf2-0"},
		[]string{"tf0-0", "tf1-2", "tf2-0"},
		[]string{"tf0-1", "tf1-2", "tf2-0"},
		[]string{"tf0-2", "tf1-2", "tf2-0"},
		[]string{"tf0-3", "tf1-2", "tf2-0"},
		[]string{"tf0-4", "tf1-2", "tf2-0"},
		[]string{"tf0-0", "tf1-3", "tf2-0"},
		[]string{"tf0-1", "tf1-3", "tf2-0"},
		[]string{"tf0-2", "tf1-3", "tf2-0"},
		[]string{"tf0-3", "tf1-3", "tf2-0"},
		[]string{"tf0-4", "tf1-3", "tf2-0"},
		[]string{"tf0-0", "tf1-4", "tf2-0"},
		[]string{"tf0-1", "tf1-4", "tf2-0"},
		[]string{"tf0-2", "tf1-4", "tf2-0"},
		[]string{"tf0-3", "tf1-4", "tf2-0"},
		[]string{"tf0-4", "tf1-4", "tf2-0"},
		[]string{"tf0-0", "tf1-0", "tf2-1"},
		[]string{"tf0-1", "tf1-0", "tf2-1"},
		[]string{"tf0-2", "tf1-0", "tf2-1"},
		[]string{"tf0-3", "tf1-0", "tf2-1"},
		[]string{"tf0-4", "tf1-0", "tf2-1"},
	}

	for i := 0; i < 30; i++ {
		out, _, err := u.ReadInputs(fps, offs)
		if err != nil {
			t.Fatalf("error in ReadInputs: %v\n", err)
		}
		for j := 0; j < len(out); j++ {
			t.Logf("%d: %s", i, out[j])
		}
		for j := 0; j < len(out); j++ {
			if out[j] != expected[i][j] {
				t.Errorf("%d failed, out[%d] = %s, expected[%d][%d] = %s", i, j, out[j], i, j, expected[i][j])
			}
		}
	}
}

func TestExhaustFile(t *testing.T) {
	files := []string{
		"../testfiles/testFile0",
		"../testfiles/testFile1",
		"../testfiles/testFile2",
	}

	offs := make([]int64, len(files))

	fps, err := u.BulkOpen(files)
	if err != nil {
		t.Fatalf("error in BulkOpen: %v\n", err)
	}

	expected := [][]string{
		[]string{"tf0-0", "tf1-0", "tf2-0"},
		[]string{"tf0-1", "tf1-0", "tf2-0"},
		[]string{"tf0-2", "tf1-0", "tf2-0"},
		[]string{"tf0-3", "tf1-0", "tf2-0"},
		[]string{"tf0-4", "tf1-0", "tf2-0"},
		[]string{"tf0-0", "tf1-1", "tf2-0"},
		[]string{"tf0-1", "tf1-1", "tf2-0"},
		[]string{"tf0-2", "tf1-1", "tf2-0"},
		[]string{"tf0-3", "tf1-1", "tf2-0"},
		[]string{"tf0-4", "tf1-1", "tf2-0"},
		[]string{"tf0-0", "tf1-2", "tf2-0"},
		[]string{"tf0-1", "tf1-2", "tf2-0"},
		[]string{"tf0-2", "tf1-2", "tf2-0"},
		[]string{"tf0-3", "tf1-2", "tf2-0"},
		[]string{"tf0-4", "tf1-2", "tf2-0"},
		[]string{"tf0-0", "tf1-3", "tf2-0"},
		[]string{"tf0-1", "tf1-3", "tf2-0"},
		[]string{"tf0-2", "tf1-3", "tf2-0"},
		[]string{"tf0-3", "tf1-3", "tf2-0"},
		[]string{"tf0-4", "tf1-3", "tf2-0"},
		[]string{"tf0-0", "tf1-4", "tf2-0"},
		[]string{"tf0-1", "tf1-4", "tf2-0"},
		[]string{"tf0-2", "tf1-4", "tf2-0"},
		[]string{"tf0-3", "tf1-4", "tf2-0"},
		[]string{"tf0-4", "tf1-4", "tf2-0"},
		[]string{"tf0-0", "tf1-0", "tf2-1"},
		[]string{"tf0-1", "tf1-0", "tf2-1"},
		[]string{"tf0-2", "tf1-0", "tf2-1"},
		[]string{"tf0-3", "tf1-0", "tf2-1"},
		[]string{"tf0-4", "tf1-0", "tf2-1"},
		[]string{"tf0-0", "tf1-1", "tf2-1"},
		[]string{"tf0-1", "tf1-1", "tf2-1"},
		[]string{"tf0-2", "tf1-1", "tf2-1"},
		[]string{"tf0-3", "tf1-1", "tf2-1"},
		[]string{"tf0-4", "tf1-1", "tf2-1"},
		[]string{"tf0-0", "tf1-2", "tf2-1"},
		[]string{"tf0-1", "tf1-2", "tf2-1"},
		[]string{"tf0-2", "tf1-2", "tf2-1"},
		[]string{"tf0-3", "tf1-2", "tf2-1"},
		[]string{"tf0-4", "tf1-2", "tf2-1"},
		[]string{"tf0-0", "tf1-3", "tf2-1"},
		[]string{"tf0-1", "tf1-3", "tf2-1"},
		[]string{"tf0-2", "tf1-3", "tf2-1"},
		[]string{"tf0-3", "tf1-3", "tf2-1"},
		[]string{"tf0-4", "tf1-3", "tf2-1"},
		[]string{"tf0-0", "tf1-4", "tf2-1"},
		[]string{"tf0-1", "tf1-4", "tf2-1"},
		[]string{"tf0-2", "tf1-4", "tf2-1"},
		[]string{"tf0-3", "tf1-4", "tf2-1"},
		[]string{"tf0-4", "tf1-4", "tf2-1"},
		[]string{"tf0-0", "tf1-0", "tf2-2"},
		[]string{"tf0-1", "tf1-0", "tf2-2"},
		[]string{"tf0-2", "tf1-0", "tf2-2"},
		[]string{"tf0-3", "tf1-0", "tf2-2"},
		[]string{"tf0-4", "tf1-0", "tf2-2"},
		[]string{"tf0-0", "tf1-1", "tf2-2"},
		[]string{"tf0-1", "tf1-1", "tf2-2"},
		[]string{"tf0-2", "tf1-1", "tf2-2"},
		[]string{"tf0-3", "tf1-1", "tf2-2"},
		[]string{"tf0-4", "tf1-1", "tf2-2"},
		[]string{"tf0-0", "tf1-2", "tf2-2"},
		[]string{"tf0-1", "tf1-2", "tf2-2"},
		[]string{"tf0-2", "tf1-2", "tf2-2"},
		[]string{"tf0-3", "tf1-2", "tf2-2"},
		[]string{"tf0-4", "tf1-2", "tf2-2"},
		[]string{"tf0-0", "tf1-3", "tf2-2"},
		[]string{"tf0-1", "tf1-3", "tf2-2"},
		[]string{"tf0-2", "tf1-3", "tf2-2"},
		[]string{"tf0-3", "tf1-3", "tf2-2"},
		[]string{"tf0-4", "tf1-3", "tf2-2"},
		[]string{"tf0-0", "tf1-4", "tf2-2"},
		[]string{"tf0-1", "tf1-4", "tf2-2"},
		[]string{"tf0-2", "tf1-4", "tf2-2"},
		[]string{"tf0-3", "tf1-4", "tf2-2"},
		[]string{"tf0-4", "tf1-4", "tf2-2"},
		[]string{"tf0-0", "tf1-0", "tf2-3"},
		[]string{"tf0-1", "tf1-0", "tf2-3"},
		[]string{"tf0-2", "tf1-0", "tf2-3"},
		[]string{"tf0-3", "tf1-0", "tf2-3"},
		[]string{"tf0-4", "tf1-0", "tf2-3"},
		[]string{"tf0-0", "tf1-1", "tf2-3"},
		[]string{"tf0-1", "tf1-1", "tf2-3"},
		[]string{"tf0-2", "tf1-1", "tf2-3"},
		[]string{"tf0-3", "tf1-1", "tf2-3"},
		[]string{"tf0-4", "tf1-1", "tf2-3"},
		[]string{"tf0-0", "tf1-2", "tf2-3"},
		[]string{"tf0-1", "tf1-2", "tf2-3"},
		[]string{"tf0-2", "tf1-2", "tf2-3"},
		[]string{"tf0-3", "tf1-2", "tf2-3"},
		[]string{"tf0-4", "tf1-2", "tf2-3"},
		[]string{"tf0-0", "tf1-3", "tf2-3"},
		[]string{"tf0-1", "tf1-3", "tf2-3"},
		[]string{"tf0-2", "tf1-3", "tf2-3"},
		[]string{"tf0-3", "tf1-3", "tf2-3"},
		[]string{"tf0-4", "tf1-3", "tf2-3"},
		[]string{"tf0-0", "tf1-4", "tf2-3"},
		[]string{"tf0-1", "tf1-4", "tf2-3"},
		[]string{"tf0-2", "tf1-4", "tf2-3"},
		[]string{"tf0-3", "tf1-4", "tf2-3"},
		[]string{"tf0-4", "tf1-4", "tf2-3"},
		[]string{"tf0-0", "tf1-0", "tf2-4"},
		[]string{"tf0-1", "tf1-0", "tf2-4"},
		[]string{"tf0-2", "tf1-0", "tf2-4"},
		[]string{"tf0-3", "tf1-0", "tf2-4"},
		[]string{"tf0-4", "tf1-0", "tf2-4"},
		[]string{"tf0-0", "tf1-1", "tf2-4"},
		[]string{"tf0-1", "tf1-1", "tf2-4"},
		[]string{"tf0-2", "tf1-1", "tf2-4"},
		[]string{"tf0-3", "tf1-1", "tf2-4"},
		[]string{"tf0-4", "tf1-1", "tf2-4"},
		[]string{"tf0-0", "tf1-2", "tf2-4"},
		[]string{"tf0-1", "tf1-2", "tf2-4"},
		[]string{"tf0-2", "tf1-2", "tf2-4"},
		[]string{"tf0-3", "tf1-2", "tf2-4"},
		[]string{"tf0-4", "tf1-2", "tf2-4"},
		[]string{"tf0-0", "tf1-3", "tf2-4"},
		[]string{"tf0-1", "tf1-3", "tf2-4"},
		[]string{"tf0-2", "tf1-3", "tf2-4"},
		[]string{"tf0-3", "tf1-3", "tf2-4"},
		[]string{"tf0-4", "tf1-3", "tf2-4"},
		[]string{"tf0-0", "tf1-4", "tf2-4"},
		[]string{"tf0-1", "tf1-4", "tf2-4"},
		[]string{"tf0-2", "tf1-4", "tf2-4"},
		[]string{"tf0-3", "tf1-4", "tf2-4"},
		[]string{"tf0-4", "tf1-4", "tf2-4"},
	}

	for i := 0; i < 125; i++ {
		out, test, err := u.ReadInputs(fps, offs)
		if err != nil {
			t.Fatalf("error in ReadInputs: %v\n", err)
		}
		for j := 0; j < len(out); j++ {
			t.Logf("%d: %s", i, out[j])
		}
		t.Logf("%t", test)
		if test != false {
			t.Errorf("%d: test should have been false, test is %t", i, test)
		}
		for j := 0; j < len(out); j++ {
			if out[j] != expected[i][j] {
				t.Errorf("%d failed, out[%d] = %s, expected[%d][%d] = %s", i, j, out[j], i, j, expected[i][j])
			}
		}
	}
	// file should be exhausted, this should return true
	out, test, err := u.ReadInputs(fps, offs)
	t.Logf("out: %q\n", out)
	t.Logf("%t\n", test)
	if test != true {
		t.Fatalf("test should have returned true, returned %t", test)
	}
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
}

func TestCountPerms(t *testing.T) {
	files := []string{
		"../testfiles/testFile0",
		"../testfiles/testFile1",
		"../testfiles/testFile2",
	}

	fps, err := u.BulkOpen(files)
	if err != nil {
		t.Fatalf("error in BulkOpen: %v\n", err)
	}

	ret, err := u.CountPerms(fps)
	if err != nil {
		t.Fatalf("some sort of error in countperms: %v", err)
	}
	t.Logf("got %d", ret)
	if ret != 125 {
		t.Fatalf("expected 125, got %d", ret)
	}
}
