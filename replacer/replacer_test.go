package replacer_test

import (
	"strings"
	"testing"

	r "github.com/alivingvendingmachine/frute/replacer"
)

func TestSearch(t *testing.T) {
	//t.Parallel()
	tests := []struct {
		inp      string
		sentinel string
		out1     int
		out2     int
	}{
		{"!test!", "!", 0, 6},
		{"!!!test!!!", "!!!", 0, 10},
		{"@@@test@@@", "@@@", 0, 10},
		{"@@@@@@@@@test@@@@@@@@@", "@@@@@@@@@", 0, 22},
		{"&&&test&&&", "&&&", 0, 10},
		{",,,test,,,", ",,,", 0, 10},
	}
	for i, test := range tests {
		o1, o2, err := r.Search(test.inp, test.sentinel)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		if o1 != test.out1 || o2 != test.out2 {
			t.Errorf("%d: got %d %d, expected %d %d", i, o1, o2, test.out1, test.out2)
		}
	}
}

func TestSearchErrors(t *testing.T) {
	//t.Parallel()
	_, _, err := r.Search("!!!test", "!!!")
	if err == nil {
		t.Errorf("expected error")
	}

	_, _, err = r.Search("!!!test!!!!!!", "!!!")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, _, err = r.Search("test", "!!!")
	if err == nil {
		t.Errorf("expected error")
	}

	_, _, err = r.Search("test", "^foo")
	if err == nil {
		t.Errorf("expected error")
	}

	_, _, err = r.Search("test", "???")
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestReplacer(t *testing.T) {
	//t.Parallel()
	tests := []struct {
		input    string
		replace  string
		sentinel string
		output   string
	}{
		{"this is a !!!test!!!", "game", "!!!", "this is a game"},
		{"this is a !!!!!!", "game", "!!!", "this is a game"},
		{"this is a @@@test@@@", "game", "@@@", "this is a game"},
		{"this is a !!!test!!!", "a", "!!!", "this is a a"},
	}

	for i, test := range tests {
		out, err := r.Replace(test.input, test.replace, test.sentinel)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		if strings.Compare(out, test.output) != 0 {
			t.Errorf("%d: expected %q, got %q", i, test.output, out)
		}
	}
}

func TestReplacerErrors(t *testing.T) {
	_, err := r.Replace("!!!test", "", "!!!")
	if err == nil {
		t.Errorf("expected error")
	}
}
