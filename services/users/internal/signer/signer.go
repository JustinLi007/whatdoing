package signer

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jws"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type Signer interface {
	NewJwt(sub, scope string, ttl time.Duration) (string, error)
	GetJwkSet() jwk.Set
}

type JwtSigner struct {
	mtx      sync.RWMutex
	priv     ed25519.PrivateKey
	pub      ed25519.PublicKey
	set      jwk.Set
	issuer   string
	audience string
	kid      string
}

var signerInstance *JwtSigner

func NewSigner(iss, aud string) (Signer, error) {
	if signerInstance != nil {
		return signerInstance, nil
	}

	base64Seed, err := base64.StdEncoding.DecodeString(os.Getenv("BASE64_SEED"))
	if err != nil {
		return nil, err
	}

	priv := ed25519.NewKeyFromSeed(base64Seed)
	pub := priv.Public().(ed25519.PublicKey)

	jwtSigner := &JwtSigner{
		mtx:      sync.RWMutex{},
		priv:     priv,
		pub:      pub,
		set:      jwk.NewSet(),
		issuer:   iss,
		audience: aud,
	}

	if _, err := jwtSigner.createJwk(pub); err != nil {
		return nil, err
	}

	signerInstance = jwtSigner

	return signerInstance, nil
}

func (s *JwtSigner) NewJwt(sub, scope string, ttl time.Duration) (string, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	now := time.Now()
	builder := jwt.NewBuilder().
		Issuer(s.issuer).
		Subject(sub).
		Audience([]string{s.audience}).
		IssuedAt(now).
		NotBefore(now.Add(time.Second*-5)).
		Expiration(now.Add(ttl)).
		Claim("scope", scope).
		JwtID(uuid.NewString())

	tok, err := builder.Build()
	if err != nil {
		return "", err
	}

	// signingKey, ok := s.set.LookupKeyID(s.kid)
	_, ok := s.set.LookupKeyID(s.kid)
	if !ok {
		return "", fmt.Errorf(`no signing key`)
	}

	jwsHeaders := jws.NewHeaders()
	jwsHeaders.Set(jws.AlgorithmKey, jwa.EdDSA())
	jwsHeaders.Set(jws.KeyIDKey, s.kid)
	signedJwtBytes, err := jwt.Sign(
		tok,
		jwt.WithKey(
			jwa.EdDSA(),
			s.priv,
			jws.WithProtectedHeaders(jwsHeaders),
		),
	)
	if err != nil {
		return "", err
	}

	return string(signedJwtBytes), nil
}

func (s *JwtSigner) GetJwkSet() jwk.Set {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.set
}

func (s *JwtSigner) createJwk(pub ed25519.PublicKey) (jwk.Key, error) {
	key, err := jwk.Import(pub)
	if err != nil {
		return nil, err
	}

	if err := jwk.AssignKeyID(key); err != nil {
		return nil, err
	}

	kid, ok := key.KeyID()
	if !ok {
		return nil, fmt.Errorf(`key no "kid"`)
	}

	if err := key.Set(jwk.AlgorithmKey, jwa.EdDSA()); err != nil {
		return nil, err
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.kid = kid
	if err := s.set.AddKey(key); err != nil {
		return nil, err
	}

	return key, nil
}
