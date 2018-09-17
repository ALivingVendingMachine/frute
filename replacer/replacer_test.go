package replacer

import (
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	//t.Parallel()
	ret1, ret2, err := search("!!!test!!!", "!!!")
	if err != nil {
		t.Errorf("error: %v", err)
	}
	t.Logf("Got: %d %d", ret1, ret2)
}

func TestSearchErrors(t *testing.T) {
	//t.Parallel()
	_, _, err := search("!!!test", "!!!")
	if err == nil {
		t.Errorf("expected error")
	}

	_, _, err = search("!!!test!!!!!!", "!!!")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, _, err = search("test", "!!!")
	if err == nil {
		t.Errorf("expected error")
	}

	_, _, err = search("test", "^foo")
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
	}

	for i, test := range tests {
		out, err := Replace(test.input, test.replace, test.sentinel)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		if strings.Compare(out, test.output) != 0 {
			t.Errorf("%d: expected %q, got %q", i, test.output, out)
		}
	}
}
