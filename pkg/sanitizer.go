package omh

import "regexp"

var insane = regexp.MustCompile(`[^a-zA-Z0-9\-]`)

func Sanitize(in string) string {
	return insane.ReplaceAllString(in, "")
}
