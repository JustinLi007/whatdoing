package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCorrectPassword(t *testing.T) {
	p := &Password{}
	err := p.Set("1234")
	require.NoError(t, err)
	match, err := p.Validate("1234")
	require.NoError(t, err)
	assert.True(t, match)

}

func TestWrongPassword(t *testing.T) {
	p := &Password{}
	err := p.Set("1234")
	require.NoError(t, err)
	match, err := p.Validate("123")
	require.NoError(t, err)
	assert.False(t, match)
}
