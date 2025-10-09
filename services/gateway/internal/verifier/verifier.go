package verifier

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type Verifier interface {
	ValidateJwt(tokenStr string) (string, string, error)
}

type jwtVerifier struct {
	mtx      sync.RWMutex
	jwkCache *jwk.Cache
	issuer   string
	audience string
	jwkUrl   string
}

var verifierInstance *jwtVerifier

func NewVerifier(url, issuer, audience string) (Verifier, error) {
	if verifierInstance != nil {
		return verifierInstance, nil
	}

	verifier := &jwtVerifier{
		mtx:      sync.RWMutex{},
		issuer:   issuer,
		audience: audience,
		jwkUrl:   url,
	}

	wl := httprc.NewMapWhitelist().Add(url)
	cliOpt := httprc.WithWhitelist(wl)
	jwkCli := httprc.NewClient(cliOpt)

	jwkCache, err := jwk.NewCache(context.Background(), jwkCli)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := jwkCache.Register(
		ctx,
		url,
		jwk.WithMinInterval(time.Hour*1),
		jwk.WithMaxInterval(time.Hour*24*7),
	); err != nil {
		return nil, err
	}

	verifier.jwkCache = jwkCache

	verifierInstance = verifier

	return verifierInstance, nil
}

func (v *jwtVerifier) ValidateJwt(tokenStr string) (string, string, error) {
	jwkSet, err := v.lookup()
	if err != nil {
		return "", "", err
	}

	parsedJwt, err := jwt.ParseString(
		tokenStr,
		jwt.WithVerify(true),
		jwt.WithValidate(true),
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.audience),
		jwt.WithKeySet(jwkSet),
		jwt.WithAcceptableSkew(time.Second*30),
	)
	if err != nil {
		return "", "", err
	}

	sub, ok := parsedJwt.Subject()
	if !ok {
		return "", "", fmt.Errorf(`token have no "sub" claim`)
	}
	if sub == "" {
		return "", "", fmt.Errorf(`token have no "sub" claim`)
	}

	var scope string
	if err := parsedJwt.Get("scope", &scope); err != nil {
		return "", "", err
	}

	return sub, scope, nil
}

func (v *jwtVerifier) lookup() (jwk.Set, error) {
	v.mtx.RLock()
	defer v.mtx.RUnlock()
	return v.jwkCache.Lookup(context.Background(), v.jwkUrl)
}
