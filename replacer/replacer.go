package replacer

import (
	"errors"
	"regexp"
)

func Search(input string, regexString string) ([][]int, error) { //(error, int, int) {
	regex, err := regexp.Compile(regexString)
	if err != nil {
		return nil, err // -1, -1
	}

	ret := regex.FindAllIndex([]byte(input), -1)
	if ret == nil {
		return nil, errors.New("search: no matches found")
	}

	return ret, nil
}
