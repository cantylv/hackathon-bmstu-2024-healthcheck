// Copyright © ivanlobanov. All rights reserved.
package functions

import (
	"net/http"
	"time"

	mc "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myconstants"
)

// SetCookieAndHeaders Sets up cookie header and csrf header.
func SetCookieAndHeaders(w http.ResponseWriter, username string) (http.ResponseWriter, error) {
	dateExp := time.Now().Add(14 * 24 * time.Hour)
	jwt, _ := NewJwtToken(NewJwtTokenProps{
		Username: username,
	}, dateExp)
	cookie := http.Cookie{
		Name:     mc.JwtToken,
		Value:    jwt,
		Expires:  dateExp,
		HttpOnly: false,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	return w, nil
}

func FlashCookie(w http.ResponseWriter, r *http.Request) {
	sessionCookie := &http.Cookie{
		Name:     mc.JwtToken,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: false,
		Path:     "/",
	}
	http.SetCookie(w, sessionCookie)
}
