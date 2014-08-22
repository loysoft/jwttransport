package jwttransport

import (
	"sync"
	"errors"
	"net/http"
	"code.google.com/p/goauth2/oauth/jwt"
)

type Transport struct {
	*Config
	*Token

	// mu guards modifying the token.
	mu sync.Mutex

	// Transport is the HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	// (It should never be an oauth.Transport.)
	Transport http.RoundTripper
}

// Client returns an *http.Client that makes OAuth-authenticated requests.
func (t *Transport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

// RoundTrip executes a single HTTP transaction using the Transport's
// Token as authorization headers.
//
// This method will attempt to renew the Token if it has expired and may return
// an error related to that Token renewal before attempting the client request.
// If the Token cannot be renewed a non-nil os.Error value will be returned.
// If the Token is invalid callers should expect HTTP-level errors,
// as indicated by the Response's StatusCode.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	accessToken, err := t.getAccessToken()
	if err != nil {
		return nil, err
	}
	// To set the Authorization header, we must make a copy of the Request
	// so that we don't modify the Request we were given.
	// This is required by the specification of http.RoundTripper.
	req = cloneRequest(req)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make the HTTP request.
	return t.transport().RoundTrip(req)
}

func (t *Transport) getAccessToken() (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Token == nil {
		if t.Config == nil {
			return "", JWTAuthError{"RoundTrip", "no Config supplied"}
		}
		if t.TokenCache == nil {
			return "", JWTAuthError{"RoundTrip", "no Token supplied"}
		}
		var err error
		t.Token, err = t.TokenCache.Token()
		if err != nil {
			return "", err
		}
	}

	// Refresh the Token if it has expired.
	if t.Token == nil || t.Expired() {
		if err := t.Refresh(); err != nil {
			return "", err
		}
	}
	if t.AccessToken == "" {
		return "", errors.New("no access token obtained from refresh")
	}
	return t.AccessToken, nil
}

func (t *Transport) PrepareToken() error {
	if t.Token == nil {
		t.Token = &Token{}
	}
	return t.Refresh()
}

// Refresh renews the Transport's AccessToken using its RefreshToken.
func (t *Transport) Refresh() error {
	if t.Token == nil {
		return JWTAuthError{"Refresh", "no existing Token"}
	}
	if t.Config == nil {
		return JWTAuthError{"Refresh", "no Config supplied"}
	}

	err := t.updateToken(t.Token)
	if err != nil {
		return err
	}
	if t.TokenCache != nil {
		return t.TokenCache.PutToken(t.Token)
	}
	return nil
}

// updateToken mutates both tok and v.
func (t *Transport) updateToken(tok *Token) error {

	jwtTok := jwt.NewToken(t.Config.ClientEmail, t.Config.Scope, []byte(t.Config.PrivateKey))
	jwtTok.ClaimSet.Aud = "https://accounts.google.com/o/oauth2/token"

	c := &http.Client{}

	o, err := jwtTok.Assert(c)
	if err != nil {
		return err
	}

	tok.AccessToken = o.AccessToken
	tok.Expiry = o.Expiry

	return nil
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}
