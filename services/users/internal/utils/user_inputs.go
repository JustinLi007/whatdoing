package utils

import "strings"

func IsValidEmail(email string) bool {
	e := strings.TrimSpace(email)
	if e == "" {
		return false
	}

	atParts := strings.Split(e, "@")
	if len(atParts) != 2 {
		return false
	}
	leftAt := strings.TrimSpace(atParts[0])
	rightAt := strings.TrimSpace(atParts[1])
	if leftAt == "" || rightAt == "" {
		return false
	}

	dotParts := strings.Split(rightAt, ".")
	if len(dotParts) != 2 {
		return false
	}
	leftDot := strings.TrimSpace(dotParts[0])
	rightDot := strings.TrimSpace(dotParts[1])
	if leftDot == "" || rightDot == "" {
		return false
	}

	return true
}

func IsValidPassword(password string) bool {
	p := strings.TrimSpace(password)
	if p == "" {
		return false
	}

	return true
}
