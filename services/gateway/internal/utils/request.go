package utils

import (
	"net/http"
	"strings"
)

func ParseRequestUrl(r *http.Request) (prefix, pathValue string, ok bool) {
	if r == nil {
		return "", "", false
	}
	if r.URL == nil {
		return "", "", false
	}

	path := strings.TrimPrefix(strings.TrimSpace(r.URL.Path), "/")
	if path == "" {
		return "", "", false
	}
	before, after, _ := strings.Cut(path, "/")
	return strings.TrimSpace(before), strings.TrimSpace(after), true
}

func HasScope(expected, actual string) bool {
	scopes := strings.SplitSeq(actual, ",")
	for v := range scopes {
		if expected == strings.ToLower(strings.TrimSpace(v)) {
			return true
		}
	}
	return false
}
