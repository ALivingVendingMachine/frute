package replacer

import (
	"errors"
	"regexp"
)

// Search takes a string of input, and a string containing which characters to
// search for (which must be a regex), and returns the very first and very last
// incidence of that string in the input.
func Search(input string, regexString string) (int, int, error) {
	regex, err := regexp.Compile(regexString)
	if err != nil {
		return -1, -1, err
	}

	ret := regex.FindAllIndex([]byte(input), -1)
	if ret == nil {
		return -1, -1, errors.New("search: no matches found")
	}

	if len(ret) < 2 {
		return -1, -1, errors.New("search: unpaired sentinel")
	}

	return ret[0][0], ret[len(ret)-1][1], nil
}

// Replace takes an input string, a string to replace it with, and a sentinel to
// search for. It then returns the input string with the input replaced
func Replace(input string, replace string, sentinel string) (string, error) {
	start, stop, err := Search(input, sentinel)

	if err != nil {
		return "", err
	}

	ret := input[:start] + replace + input[stop:]

	return ret, nil
}
