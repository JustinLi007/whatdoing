package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Envelope map[string]any

func WriteJson(w http.ResponseWriter, statusCode int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(js)
	return err
}

func RequireNoError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

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

func SetCookie(w http.ResponseWriter, name, value string) {
	// TODO: add expiration
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   "localhost",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   false,
	})
}

func DeleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Domain:   "localhost",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	})
}
