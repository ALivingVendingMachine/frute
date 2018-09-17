package replacer_test

import (
	"testing"

	r "github.com/alivingvendingmachine/frute/replacer"
)

func TestMain(t *testing.T) {
	ret, err := r.Search("This is !!!a test!!!", "!!!")
	if err != nil {
		t.Errorf("error: %v")
	}
	t.Logf("GOT %q", ret)
}
