package helpers

import (
	"net/http"
	"time"

	"github.com/boratanrikulu/s-lyrics/models"
)

func SetTokenCookies(r models.RefreshAndAccessTokens, w http.ResponseWriter) {
	oneMonth := time.Hour * 24 * 30
	cookies := []http.Cookie {
		http.Cookie {
			Name: "AccessToken",
			Value: r.Response.AccessToken,
			// MaxAge of access token is 1 hour.
			MaxAge: r.Response.ExpiresIn,
		},
		http.Cookie {
			Name: "RefreshToken",
			Value: r.Response.RefreshToken,
			// We will set max age for refresh token 1 month.
			MaxAge: int(oneMonth.Seconds()),
		},
	}
	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
}

func UpdateTokenCookies(u models.UpdateAccessToken, w http.ResponseWriter) {
	cookies := []http.Cookie {
		http.Cookie {
			Name: "AccessToken",
			Value: u.Response.AccessToken,
			// MaxAge of access token is 1 hour.
			MaxAge: u.Response.ExpiresIn,
		},
	}

	// If there is a new refresh token.
	// Update it too.
	if u.Response.RefreshToken != "" {
		oneMonth := time.Hour * 24 * 30
		cookie := http.Cookie{
			Name: "RefreshToken",
			Value: u.Response.RefreshToken,
			// We will set max age for refresh token 1 month.
			MaxAge: int(oneMonth.Seconds()),
		}
		cookies = append(cookies, cookie)
	}

	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
}

func SetStateCookie(r models.Authorization, w http.ResponseWriter) {
	cookie := http.Cookie {
		Name: "State",
		Value: r.Request.State,
	}
	http.SetCookie(w, &cookie)
}

func ClearCookies(w http.ResponseWriter, r *http.Request) {
	// Clears all cookies.
	for _, cookie := range r.Cookies() {
		cookie.Value = ""
		cookie.Expires = time.Unix(0, 0)
		http.SetCookie(w, cookie)
	}
}
