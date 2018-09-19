package fuzzer_test

import (
	"strings"
	"testing"

	f "github.com/alivingvendingmachine/frute/fuzzer"
)

func TestMutateString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		inp   string
		seed  int64
		iters int
		out   string
	}{
		{"test", 482914785, 1, "est"},                                      // 0
		{"test", 111111111, 1, "est"},                                      // 1
		{"test", 111111111, 5, "estet徥�t\ue765쓥et"},                        // 2
		{"smash the state", 3212321, 1, "sah\ued6cesate"},                  // 3
		{"smash the state", 3212321, 3, "sah\ued6cesateᾞs 簾 ttems㿭Ḷh㛭t겤⍴"}, // 4
		{"smash the state", 3212321, 5, "sah\ued6cesateᾞs 簾 ttems㿭Ḷh㛭t겤⍴䕎shte\uf05e蕛ems罸h绉\ue908t流"}, // 5
		{"Hello world!", 1234567890, 1, "\u05cdlo咢㵲d!"},                                              // 6
		{"Hello world!", 13579, 3, "el 삅l⺭!elo웼rd!悆彷 ol!"},                                           // 7
		{"Hello world!", 13579, 5, "el 삅l⺭!elo웼rd!悆彷 ol!el셑譋o㡌ꊖel艴o꒻!"},                              // 8
		{"bye", 246810, 2, "yeye"}, // 9 YEYE!
	}
	for i, test := range tests {
		this, err := f.MutateString(test.inp, test.seed, test.iters)
		t.Logf("%q", this)
		if err != nil {
			t.Errorf("%d: error: %v", i, err)
		}
		if strings.Compare(this, test.out) != 0 {
			t.Errorf("%d: expected %q, got %q", i, test.out, this)
		}
	}
}

func TestMutateStringASCII(t *testing.T) {
	t.Parallel()
	tests := []struct {
		inp   string
		seed  int64
		iters int
		out   string
	}{
		{"test", 482914785, 1, "est"},                                                                                       // 0
		{"test", 111111111, 1, "est"},                                                                                       // 1
		{"test", 111111111, 5, "estet%\x1cteeet"},                                                                           // 2
		{"no gods no masters", 123456789, 1, "o&\x02sn {atLs"},                                                              // 3
		{"no gods no masters", 123456789, 3, "o&\x02sn {atLsog\\1 -n\x1bezogd omues"},                                       // 4
		{"no gods no masters", 123456789, 5, "o&\x02sn {atLsog\\1 -n\x1bezogd omueso,d n atrFogm Imyes"},                    // 5
		{"Haymarket affair", 2222222, 1, "am7^ fa;"},                                                                        // 6
		{"Haymarket affair", 2222222, 3, "am7^ fa;amre faramre \x1car"},                                                     // 7
		{"Haymarket affair", 2222222, 5, "am7^ fa;amre faramre \x1carHyaNtafirHyakt]f\x04&r"},                               // 8
		{"emma goldman", 321321321321321, 10, "mahlA<ma\x18lmnmBglmnmagomnm\x13gl4nmaglmnw*a\t\x11lnmag(p\x1am0g[enoaglmn"}, // 9
	}
	for i, test := range tests {
		this, err := f.MutateStringASCII(test.inp, test.seed, test.iters)
		t.Logf("%q", this)
		if err != nil {
			t.Errorf("%d: error: %v", i, err)
		}
		if strings.Compare(this, test.out) != 0 {
			t.Errorf("%d: expected %q, got %q", i, test.out, this)
		}
	}
}

func TestMutateStringErrors(t *testing.T) {
	t.Parallel()
	_, err := f.MutateString("", 12345, 1)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}

	_, err = f.MutateString("hello", 12345, 0)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}

	_, err = f.MutateString("hello", 12345, -1)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestMutateStringASCIIErrors(t *testing.T) {
	t.Parallel()
	_, err := f.MutateStringASCII("", 12345, 1)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}

	_, err = f.MutateStringASCII("hello", 12345, 0)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}

	_, err = f.MutateStringASCII("hello", 12345, -1)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestRandomString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		length int
		seed   int64
		iters  int
		output string
	}{
		{8, 23452345, 1, "셧账\xe6\xac"},                  // 0
		{8, 23452345, 3, "\u0085è¦¬"},                   // 1
		{8, 23452345, 5, "\u0085è¦¬"},                   // 2
		{8, 123123, 1, "㷠᧺\xe0\xaf"},                    // 3
		{8, 123123, 3, "᪱\ue288\xef\xbf"},               // 4
		{8, 123123, 5, "᪱\ue288\xef\xbf"},               // 5
		{16, 999999999, 1, "씽퐵ᦶ쭮\uf12a\xe8"},            // 6
		{16, 999999999, 3, "\u0094½쭮á¶\u00ad®\xc2"},     // 7
		{16, 999999999, 5, "\u0094½쭮á¶\u00ad®\xc2"},     // 8
		{32, 123123, 5, "᪱\ue288�¯뉝㷱\u0095ï¿í¥\u008c᳗"}, // 9
	}

	for i, test := range tests {
		this, err := f.RandomString(test.length, test.seed, test.iters)
		t.Logf("%q", this)
		if err != nil {
			t.Errorf("%d: error: %v", i, err)
		}
		if len(this) != test.length {
			t.Errorf("%d: len should be %d, got %d", i, test.length, len(this))
		}
		if strings.Compare(this, test.output) != 0 {
			t.Errorf("%d: expected %q, got %q", i, test.output, this)
		}
	}
}

func TestRandomStringErrors(t *testing.T) {
	t.Parallel()
	_, err := f.RandomString(0, 999, 1)
	if err == nil {
		t.Errorf("0: expected error, got nothing")
	}
	_, err = f.RandomString(1, 999, 0)
	if err == nil {
		t.Errorf("1: expected error, got nothing")
	}
	_, err = f.RandomString(1, 999, -1)
	if err == nil {
		t.Errorf("1: expected error, got nothing")
	}
}

func TestRandomStringASCII(t *testing.T) {
	t.Parallel()
	tests := []struct {
		length int
		seed   int64
		iters  int
		output string
	}{
		{8, 123456789, 1, "\x1cV&\r\x02M!p"},                                                                 // 0
		{8, 123456789, 3, "V&\x02!pV1\r"},                                                                    // 1
		{8, 123456789, 5, "V&\x02!pV1\r"},                                                                    // 2
		{16, 987654321, 1, "e\x17Bc\x1f\x13\"{B}8\x1dz\x13Tg"},                                               // 3
		{16, 987654321, 3, "\x17c\x13{}\x1d\x13g\x1dc\x13{B8zT"},                                             // 4
		{16, 987654321, 5, "\x17c\x13{}\x1d\x13g\x1dc\x13{B8zT"},                                             // 5
		{32, 987654321, 1, "e\x17Bc\x1f\x13\"{B}8\x1dz\x13Tgn\arlU\x1eL\x19\x045vWy`0\x14"},                  // 6
		{32, 987654321, 3, "\x17c\x13{}\x1d\x13g\x1dl\x1e\x19\x04vy0U\x17\x19\x13{W\x1d\x13g\al.\x1e\x195W"}, // 7
		{32, 987654321, 5, "\x17c\x13{}\x1d\x13g\x1dl\x1e\x19\x04vy0U\x17\x19\x13{W\x1d\x13g\al.\x1e\x195W"}, // 8
		{64, 987654321, 5, "\x17c\x13{}\x1d\x13g\x1dl\x1e\x19\x04vy0UB\x19\x02\x17W!%\x0e'\"d\x18\x1bV\x18m$\x17c\x13\"B/zTnrZL\x1bvy0<BM\x02\x17wP\x1f/0\ngHP"}, // 9
	}

	for i, test := range tests {
		this, err := f.RandomStringASCII(test.length, test.seed, test.iters)
		t.Logf("%q", this)
		if err != nil {
			t.Errorf("%d: error: %v", i, err)
		}
		if len(this) != test.length {
			t.Errorf("%d: len should be %d, got %d", i, test.length, len(this))
		}
		if strings.Compare(this, test.output) != 0 {
			t.Errorf("%d: expected %q got %q", i, test.output, this)
		}
	}
}

func TestRandomStringASCIIErrors(t *testing.T) {
	t.Parallel()
	_, err := f.RandomStringASCII(0, 999, 1)
	if err == nil {
		t.Errorf("0: expected error, got nothing")
	}
	_, err = f.RandomStringASCII(1, 999, 0)
	if err == nil {
		t.Errorf("1: expected error, got nothing")
	}
	_, err = f.RandomStringASCII(1, 999, -1)
	if err == nil {
		t.Errorf("1: expected error, got nothing")
	}
}

func TestRandomInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		limit  int
		seed   int64
		output int
	}{
		{10, 1, 1},       // 0
		{100, 10, 54},    // 0
		{1000, 100, 183}, // 0
	}
	for i, test := range tests {
		this, err := f.RandomInt(test.limit, test.seed)
		t.Logf("%q", this)
		if err != nil {
			t.Errorf("%d: error: %v", i, err)
		}
		if test.output != this {
			t.Errorf("%d: expected %d got %d", i, test.output, this)
		}
	}
}

func TestRandomIntErrors(t *testing.T) {
	t.Parallel()
	_, err := f.RandomInt(0, 0)
	if err == nil {
		t.Errorf("0: got nothing, expected error")
	}
	_, err = f.RandomInt(-10, 0)
	if err == nil {
		t.Errorf("1: got nothing, expected error")
	}
}

func TestMutateSelection(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		sentinel string
		seed     int64
		iters    int
		output   string
	}{
		{"hello !world!", "!", 1234, 1, "hello o⟭d"},
		{"hello !!!world!!!", "!!!", 1234, 1, "hello o⟭d"},
		{"hello !world!", "!", 1234, 3, "hello o⟭doldo갶d"},
		{"hello !!!world!!!", "!!!", 1234, 3, "hello o⟭doldo갶d"},
		{"hello !world!", "!", 1234, 5, "hello o⟭doldo갶d∻l푌ᢪ睁d"},
		{"hello !!!!world!!!!", "!!!!", 1234, 5, "hello o⟭doldo갶d∻l푌ᢪ睁d"},
		{"destroy!er!", "!", 9876543210, 1, "destroyr"},
		{"destroy!!er!!", "!!", 9876543210, 1, "destroyr"},
		{"destroy!!!er!!!", "!!!", 9876543210, 1, "destroyr"},
		{"destroy!!!!er!!!!", "!!!!", 9876543210, 1, "destroyr"},
		{"destroy!!!er!!!!", "!!!", 9876543210, 1, "destroyr!"},
	}
	for i, test := range tests {
		out, err := f.MutateSelection(test.input, test.sentinel, test.seed, test.iters)
		t.Logf("got %q", out)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		if strings.Compare(test.output, out) != 0 {
			t.Errorf("%d: got %q, expected %q", i, out, test.output)
		}
	}
}

func TestMutateSelectionASCII(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		sentinel string
		seed     int64
		iters    int
		output   string
	}{
		{"hello !world!", "!", 1234, 1, "hello omd"},
		{"hello !!!world!!!", "!!!", 1234, 1, "hello omd"},
		{"hello !world!", "!", 1234, 3, "hello omdoldo6d"},
		{"hello !!!world!!!", "!!!", 1234, 3, "hello omdoldo6d"},
		{"hello !world!", "!", 1234, 5, "hello omdoldo6d;lL*Ad"},
		{"hello !!!!world!!!!", "!!!!", 1234, 5, "hello omdoldo6d;lL*Ad"},
		{"destroy!er!", "!", 9876543210, 1, "destroyr"},
		{"destroy!!er!!", "!!", 9876543210, 1, "destroyr"},
		{"destroy!!!er!!!", "!!!", 9876543210, 1, "destroyr"},
		{"destroy!!!!er!!!!", "!!!!", 9876543210, 1, "destroyr"},
		{"destroy!!!er!!!!", "!!!", 9876543210, 1, "destroyr!"},
	}
	for i, test := range tests {
		out, err := f.MutateSelectionASCII(test.input, test.sentinel, test.seed, test.iters)
		t.Logf("got %q", out)
		if err != nil {
			t.Errorf("error: %v", err)
		}
		if strings.Compare(test.output, out) != 0 {
			t.Errorf("%d: got %q, expected %q", i, out, test.output)
		}
	}
}
