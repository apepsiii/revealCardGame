package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"time"
)

const (
	CookieName   = "revealcard_session"
	CookieMaxAge = 86400 * 7 // 7 hari
)

// GetCredentials mengambil username dan password admin dari environment atau default
func GetCredentials() (string, string) {
	user := os.Getenv("ADMIN_USER")
	if user == "" {
		user = "admin"
	}
	pass := os.Getenv("ADMIN_PASS")
	if pass == "" {
		pass = "admin123"
	}
	return user, pass
}

func getSecret() []byte {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "revealcard-super-secret-key-2026"
	}
	return []byte(secret)
}

// GenerateToken membuat token HMAC sederhana untuk session cookie
func GenerateToken(username string) string {
	h := hmac.New(sha256.New, getSecret())
	h.Write([]byte(username))
	return hex.EncodeToString(h.Sum(nil))
}

// SetAuthCookie memasang HTTP cookie terenkripsi/terverifikasi
func SetAuthCookie(w http.ResponseWriter, username string) {
	token := GenerateToken(username)
	cookieVal := username + ":" + token

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    cookieVal,
		Path:     "/",
		MaxAge:   CookieMaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearAuthCookie menghapus cookie saat logout
func ClearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
}

// IsAuthenticated memeriksa validitas sesi user dari cookie
func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie(CookieName)
	if err != nil || cookie.Value == "" {
		return false
	}

	expectedUser, _ := GetCredentials()
	expectedToken := GenerateToken(expectedUser)
	expectedVal := expectedUser + ":" + expectedToken

	return cookie.Value == expectedVal
}

// MiddlewareAuth memproteksi endpoint admin, redirect ke /login jika belum auth
func MiddlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

