package utils

import (
	"encoding/base64"
	"log"
	"math/rand"
	"net/http"
	"time"
)

/*
Value aus Cookie
*/
func GetInfoFromCookie(c *http.Cookie) string {
	if c != nil {
		log.Print("Token from Cookie: " + c.Value)
		return c.Value
	}
	log.Print("Cookie is null")
	return ""
}

/*
schreibt Token in Cookie
*/
func GenerateTokenCookie(w http.ResponseWriter, n string, value string, e time.Time) {
	cookie := http.Cookie{Name: n, Value: value, Expires: e, HttpOnly: true, Secure: true}
	http.SetCookie(w, &cookie)
	log.Print("set new cookie " + n)
}

/*
generiert Cookie mit random state
*/
func GenerateStateCookie(w http.ResponseWriter) string {
	var exp = time.Now().Add(365 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: exp, HttpOnly: true, Secure: true}
	http.SetCookie(w, &cookie)
	return state
}
