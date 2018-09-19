package main

import (
	"testing"

	f "github.com/alivingvendingmachine/frute/fuzzer"
	r "github.com/alivingvendingmachine/frute/replacer"
)

func TestFuzzAndReplace(t *testing.T) {
	fuzz, err := f.MutateString("hello world!", 12345, 1)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	ret, err := r.Replace("This is a !!!test!!!", fuzz, "!!!")
	if err != nil {
		t.Errorf("error: %v", err)
	}
	t.Logf("GOT %q", ret)
}
