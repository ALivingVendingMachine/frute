package brute_test

import (
	"strings"
	"testing"

	"github.com/ALivingVendingMachine/frute/brute"
)

func TestBruteForce(t *testing.T) {
	s := []string{"!!!", "@@@", "###", "---", "~~~", "&&&", "===", ",,,", "<<<", "___"}
	i1 := []string{"smash", "no"}
	i2 := []string{"all cops are bad"}
	i3 := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	tests := []struct {
		sents  []string
		inputs []string
		str    string
		exp    string
	}{
		{s, i1, "!!!protect!!! the state @@@yes@@@ to kings", "smash the state no to kings"},
		{s, i2, "!!!everything is fine!!!", "all cops are bad"},
		{s, i3, "!!!0!!! @@@9@@@ ###8### ---7--- ~~~6~~~ &&&5&&& ===4=== ,,,3,,, <<<2<<< ___1___", "1 2 3 4 5 6 7 8 9 0"},
	}

	for i, test := range tests {
		out, err := brute.Forcer(test.str, test.inputs, test.sents)
		if err != nil {
			t.Errorf("error %d: %v", i, err)
		}
		t.Logf("test %d: %s", i, out)
		if strings.Compare(out, test.exp) != 0 {
			t.Errorf("test %d: expect %s, got %s", i, test.exp, out)
		}
	}
}
