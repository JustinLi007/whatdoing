package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidEmail(t *testing.T) {
	email := "sample@something.com"
	valid := IsValidEmail(email)
	assert.True(t, valid)
}

func TestValidEmail2(t *testing.T) {
	email := "a@b.c"
	valid := IsValidEmail(email)
	assert.True(t, valid)
}

func TestInvalidEmail(t *testing.T) {
	email := "@something.com"
	valid := IsValidEmail(email)
	assert.False(t, valid)
}

func TestInvalidEmail2(t *testing.T) {
	email := "sample@.com"
	valid := IsValidEmail(email)
	assert.False(t, valid)
}

func TestInvalidEmail3(t *testing.T) {
	email := "sample@"
	valid := IsValidEmail(email)
	assert.False(t, valid)
}

func TestInvalidEmail4(t *testing.T) {
	email := "sample@something"
	valid := IsValidEmail(email)
	assert.False(t, valid)
}

func TestInvalidEmail5(t *testing.T) {
	email := "sample@something."
	valid := IsValidEmail(email)
	assert.False(t, valid)
}

func TestInvalidEmail6(t *testing.T) {
	email := "something.com"
	valid := IsValidEmail(email)
	assert.False(t, valid)
}

func TestInvalidEmail7(t *testing.T) {
	email := "@.com"
	valid := IsValidEmail(email)
	assert.False(t, valid)
}
