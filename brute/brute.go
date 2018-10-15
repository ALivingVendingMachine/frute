package brute

import (
	"github.com/alivingvendingmachine/frute/replacer"
)

// Forcer takes an input, an array of replacements, and an array of sentinels.
// Then, at each sentinel it uses the matching replacement, and returns the input
// with the replacements.
func Forcer(input string, replace []string, sentinels []string) (string, error) {
	for i := range replace {
		if i >= len(sentinels) {
			break
		}
		start, stop, err := replacer.Search(input, sentinels[i])
		if err != nil {
			if err.Error() != "search: no matches found" {
				return "", err
			}
		}
		input = input[:start] + replace[i] + input[stop:]
	}

	return input, nil
}
