package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"golang.org/x/oauth2/jws"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type Certs struct {
	Keys   map[string]*rsa.PublicKey
	Expiry time.Time
}

type key struct {
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"Kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type ClaimSet struct {
	jws.ClaimSet
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}

type response struct {
	Keys []*key `json:"keys"`
}

var (
	certs *Certs

	// Google Sign on certificates.
	googleOAuth2CertsURL = "https://www.googleapis.com/oauth2/v3/certs"
)

func getGoogleCerts() (*Certs, error) {

	log.Print("get Certs from Google")
	cacheAge := int64(7200)

	res, _ := http.Get(googleOAuth2CertsURL)
	resp := &response{}
	json.NewDecoder(res.Body).Decode(&resp)

	keys := map[string]*rsa.PublicKey{}
	for _, key := range resp.Keys {
		if key.Use == "sig" && key.Kty == "RSA" {
			n, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				return nil, err
			}
			e, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				return nil, err
			}
			ei := big.NewInt(0).SetBytes(e).Int64()
			if err != nil {
				return nil, err
			}
			keys[key.Kid] = &rsa.PublicKey{
				N: big.NewInt(0).SetBytes(n),
				E: int(ei),
			}
		}
	}

	certs = &Certs{
		Keys:   keys,
		Expiry: time.Now().Add(time.Second * time.Duration(cacheAge)),
	}

	return certs, nil

}

func parseJWT(t string) (*jws.Header, error) {
	s := strings.Split(t, ".")
	dh, err := base64.RawURLEncoding.DecodeString(s[0])
	h := &jws.Header{}
	err = json.NewDecoder(bytes.NewBuffer(dh)).Decode(h)

	if err != nil {
		log.Print(err)
		return nil, err
	}
	log.Print("parsed header", h)
	return h, nil
}

func decodeJWT(t string) (*ClaimSet, error) {
	s := strings.Split(t, ".")

	decoded, err := base64.RawURLEncoding.DecodeString(s[1])
	if err != nil {
		return nil, err
	}
	c := &ClaimSet{}
	err = json.NewDecoder(bytes.NewBuffer(decoded)).Decode(c)
	return c, err
}

func verifyJWT(t string) bool {
	c, err := getGoogleCerts()
	log.Print(c.Keys)
	h, err := parseJWT(t)
	key := c.Keys[h.KeyID]
	cs, err := decodeJWT(t)
	log.Print(cs.Exp)

	if key == nil {
		return false
	}
	err = jws.Verify(t, key)
	if err != nil {
		return false
	}
	if time.Now().Unix() > cs.Exp+int64(time.Second*60) {
		return false
	}

	return true
}
