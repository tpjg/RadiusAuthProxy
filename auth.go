package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/kirves/goradius"
)

var (
	radiusServer string
	radiusPort   string
	radiusSecret string
	cookieName   string
	auth         *goradius.AuthenticatorT
	sk           *securecookie.SecureCookie
)

func init() {
	flag.StringVar(&radiusServer, "radserver", "127.0.0.1", "Radius server IP address or hostname")
	flag.StringVar(&radiusPort, "radport", "1812", "Radius server port")
	flag.StringVar(&radiusSecret, "radsecret", "testing123", "Radius shared secret")
	flag.StringVar(&cookieName, "cookie", "RADAUTH", "Name of cookie")
}

func initAuth() {
	auth = goradius.Authenticator(radiusServer, radiusPort, radiusSecret)
	sk = securecookie.New(securecookie.GenerateRandomKey(16), securecookie.GenerateRandomKey(16))
}

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	// Check if cookie is set
	if cookie, err := r.Cookie(cookieName); err == nil {
		value := make(map[string]string)
		if err = sk.Decode(cookieName, cookie.Value, &value); err == nil {
			//fmt.Printf("Forwarding request for user %v\n", value["user"])
			fmt.Printf("%v [%v] %s %s %s \"%s\"\n", time.Now().Format("2006/01/02 15:04:05"), value["user"], r.Host, r.RemoteAddr, r.Method, r.URL.Path)
			return true
		}
		fmt.Printf("Cannot decode %v - %v\n", cookieName, err)
	} else {
		fmt.Printf("Cookie %v not found - %v\n", cookieName, err)
	}

	// Otherwise authenticate user
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return false
	}
	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false
	}
	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false
	}
	fmt.Printf("Authenticating user %v to radius server %v:%v\n", pair[0], radiusServer, radiusPort)
	ok, err := auth.Authenticate(pair[0], pair[1], "")
	if err != nil {
		fmt.Printf("Error trying to authenticate - %v", err)
	} else if ok {
		fmt.Printf("Storing user %v in cookie %v for subsequent requests\n", pair[0], cookieName)
		// Set the username in a secure cookie so that can be used for subsequent
		// requests
		value := map[string]string{
			"user": pair[0],
		}
		if encoded, err := sk.Encode(cookieName, value); err == nil {
			cookie := &http.Cookie{
				Name:  cookieName,
				Value: encoded,
				Path:  "/",
			}
			http.SetCookie(w, cookie)
		} else {
			fmt.Printf("Error setting cookie - %v\n", err)
		}
		fmt.Printf("%v [%v] %s %s %s \"%s\"\n", time.Now().Format("2006/01/02 15:04:05"), pair[0], r.Host, r.RemoteAddr, r.Method, r.URL.Path)
	}
	return ok
}
