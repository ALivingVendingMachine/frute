package fuzzer

import (
	"errors"
	"math/rand"
	"strings"

	r "github.com/alivingvendingmachine/frute/replacer"
)

// MutateString takes a string and a seed for a random number generator.  Then,
// it randomly mutates a string and returns it
func MutateString(input string, seed int64, iters int) (string, error) {
	if strings.Compare(input, "") == 0 {
		return "", errors.New("MutateString: input string cannot be empty")
	}
	if iters < 1 {
		return "", errors.New("MutateString: iters cannot be less than 1")
	}
	var ret []rune
	r := rand.New(rand.NewSource(seed))

	for i := 0; i < iters; i++ {
		for j := 0; j < len(input); j++ {
			randInt := r.Intn(100)
			if randInt < 5 && j != 0 { //occasionally step back
				j--
			} else if randInt < 95 && j != len(input)-1 { //occasionally step forward
				j++
			}
			if randInt < 10 { // new rune based on input
				ret = append(ret, rune((int(input[j])+r.Intn(65536))%65536)) // Using dumb UTF
			} else if randInt < 30 { // completely new rune
				ret = append(ret, rune(r.Intn(65536))) // Using dumb UTF
			} else { // old char
				ret = append(ret, rune(input[j]))
			}
		}
	}
	return string(ret), nil
}

// MutateStringASCII takes a string, a seed for a random number generator, and a
// number of iterations, and mutates a string, returning an ascii string
func MutateStringASCII(input string, seed int64, iters int) (string, error) {
	if strings.Compare(input, "") == 0 {
		return "", errors.New("MutateStringASCII: input string cannot be empty")
	}
	if iters < 1 {
		return "", errors.New("MutateStringASCII: iters cannot be less than 1")
	}
	var ret []rune
	r := rand.New(rand.NewSource(seed))

	for i := 0; i < iters; i++ {
		for j := 0; j < len(input); j++ {
			randInt := r.Intn(100)
			if randInt < 5 && j != 0 { //occasionally step back
				j--
			} else if randInt < 95 && j != len(input)-1 { //occasionally step forward
				j++
			}
			if randInt < 10 { // new rune based on input
				ret = append(ret, rune((int(input[j])+r.Intn(128))%128)) // Using dumb UTF
			} else if randInt < 30 { // completely new rune
				ret = append(ret, rune(r.Intn(128))) // Using dumb UTF
			} else { // old char
				ret = append(ret, rune(input[j]%128))
			}
		}
	}
	return string(ret), nil
}

// RandomString takes a length as an int, a seed (int64) and a number of iterations.
// Then, it generates a random string of length "length"
func RandomString(length int, seed int64, iters int) (string, error) {
	if length < 1 {
		return "", errors.New("RandomString: length cannot be less than 1")
	}
	if iters < 1 {
		return "", errors.New("RandomString: number of iterations cannot be less than 1")
	}
	// generate string
	var ret []rune
	r := rand.New(rand.NewSource(seed))
	for j := 0; j < length; j++ {
		ret = append(ret, rune(r.Intn(65536)))
	}

	// mutate string
	if iters-1 > 0 {
		hold, err := MutateString(string(ret), seed, iters-1)
		if err != nil {
			return "", err
		}
		ret = []rune(hold)
	}

	return string(ret)[:length], nil
}

// RandomStringASCII takes a length (int), a seed for a PRNG (int64), and a number
// of iterations, and returns a randomly generated string of length "length"
func RandomStringASCII(length int, seed int64, iters int) (string, error) {
	if length < 1 {
		return "", errors.New("RandomString: length cannot be less than 1")
	}
	if iters < 1 {
		return "", errors.New("RandomString: number of iterations cannot be less than 1")
	}
	// generate string
	var ret []rune
	r := rand.New(rand.NewSource(seed))
	for j := 0; j < length; j++ {
		ret = append(ret, rune(r.Intn(128)))
	}

	// mutate string
	if iters-1 > 0 {
		hold, err := MutateStringASCII(string(ret), seed, iters-1)
		if err != nil {
			return "", err
		}
		ret = []rune(hold)
	}

	return string(ret)[:length], nil
}

// RandomInt takes a limit and and a seed, and returns a random integer [0, limit)
func RandomInt(limit int, seed int64) (int, error) {
	if limit < 1 {
		return 0, errors.New("RandomInt: limit cannot be negative")
	}

	r := rand.New(rand.NewSource(seed))
	return r.Intn(limit), nil
}

// MutateSelection takes an input, sentinel (both strings), a seed (int64), and
// a number of iterations.  It then mutates the string, and returns the input
// with the selection (between sentinels) mutated.
func MutateSelection(input string, sentinel string, seed int64, iters int) (string, string, error) {
	start, stop, err := r.Search(input, sentinel)
	if err != nil {
		return "", "", err
	}

	fuzz := input[start+(len(sentinel)) : stop-(len(sentinel))]

	fuzzed, err := MutateString(fuzz, seed, iters)
	if err != nil {
		return "", "", err
	}

	ret, err := r.Replace(input, fuzzed, sentinel)
	if err != nil {
		return "", "", err
	}

	return ret, fuzzed, nil
}

// MutateSelectionASCII takes an input, sentinel (both strings), a seed (int64), and
// a number of iterations.  It then mutates the string, and returns the input
// with the selection (between sentinels) mutated in the ASCII range.
func MutateSelectionASCII(input string, sentinel string, seed int64, iters int) (string, error) {
	start, stop, err := r.Search(input, sentinel)
	if err != nil {
		return "", err
	}

	fuzz := input[start+(len(sentinel)) : stop-(len(sentinel))]

	fuzzed, err := MutateStringASCII(fuzz, seed, iters)
	if err != nil {
		return "", err
	}

	ret, err := r.Replace(input, fuzzed, sentinel)
	if err != nil {
		return "", err
	}

	return ret, nil
}
