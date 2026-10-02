package handler

import (
	"html/template"
	"net/http"
	"revealCard/internal/auth"
	"strings"
)

// AuthHandler menangani alur login dan logout admin
type AuthHandler struct {
	tmpl *template.Template
}

// NewAuthHandler membuat instance AuthHandler baru
func NewAuthHandler(tmpl *template.Template) *AuthHandler {
	return &AuthHandler{tmpl: tmpl}
}

// ShowLogin menampilkan halaman form login
func (h *AuthHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	if auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	hasError := r.URL.Query().Get("error") == "1"
	loggedOut := r.URL.Query().Get("logout") == "1"

	data := map[string]interface{}{
		"Title":     "Login Guru / Admin - Card Reveal Challenge",
		"HasError":  hasError,
		"LoggedOut": loggedOut,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "login.html", data); err != nil {
		http.Error(w, "Gagal merender halaman login: "+err.Error(), http.StatusInternalServerError)
	}
}

// Login memproses form verifikasi username & password
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	expectedUser, expectedPass := auth.GetCredentials()

	if username == expectedUser && password == expectedPass {
		auth.SetAuthCookie(w, username)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
}

// Logout menghapus session dan mengarahkan kembali ke halaman login
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	auth.ClearAuthCookie(w)
	http.Redirect(w, r, "/login?logout=1", http.StatusSeeOther)
}

