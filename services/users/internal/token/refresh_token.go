package token

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"time"
)

const (
	REFRESH_TOKEN_TTL = time.Hour * 24
)

type Token struct {
	PlainText string    `json:"plain_text"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
}

func NewToken(ttl time.Duration) *Token {
	t := &Token{}

	b := make([]byte, 32)
	rand.Read(b)

	hash := sha256.Sum256(b)

	t.Hash = hash[:]
	t.PlainText = base64.StdEncoding.EncodeToString(hash[:])
	t.Expiry = time.Now().Add(ttl).UTC()

	return t
}

func (t *Token) Valid() bool {
	return subtle.ConstantTimeCompare([]byte(t.PlainText), t.Hash) == 1
}

func (t *Token) GetPlainText() string {
	if t.PlainText != "" {
		return t.PlainText
	}
	if t.Hash == nil || len(t.Hash) != 32 {
		return ""
	}
	t.PlainText = base64.StdEncoding.EncodeToString(t.Hash)
	return t.PlainText
}

func (t *Token) GetHash() []byte {
	if t.Hash != nil && len(t.Hash) == 32 {
		return t.Hash
	}
	if t.PlainText == "" {
		return nil
	}
	hash, err := base64.StdEncoding.DecodeString(t.PlainText)
	if err != nil {
		return nil
	}
	t.Hash = hash
	return t.Hash
}
