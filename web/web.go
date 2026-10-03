package web

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed static/* templates/*
var EmbeddedFiles embed.FS

// GetStaticFS mengembalikan http.FileSystem baik dari folder disk lokal (jika ada) atau dari embedded binary
func GetStaticFS() http.FileSystem {
	// Cek apakah folder web/static ada di disk lokal (mode development)
	if fi, err := os.Stat(filepath.Join("web", "static")); err == nil && fi.IsDir() {
		return http.Dir(filepath.Join("web", "static"))
	}

	// Gunakan embedded filesystem (mode production / standalone binary)
	sub, err := fs.Sub(EmbeddedFiles, "static")
	if err != nil {
		panic("Gagal mengakses embedded static filesystem: " + err.Error())
	}
	return http.FS(sub)
}

// LoadTemplates memuat template HTML dari disk jika tersedia, atau dari embedded binary
func LoadTemplates() (*template.Template, error) {
	// Cek apakah folder templates ada di disk lokal
	if fi, err := os.Stat(filepath.Join("web", "templates")); err == nil && fi.IsDir() {
		return template.ParseFiles(
			filepath.Join("web", "templates", "game.html"),
			filepath.Join("web", "templates", "login.html"),
			filepath.Join("web", "templates", "admin", "index.html"),
		)
	}

	// Fallback ke embedded templates
	return template.ParseFS(
		EmbeddedFiles,
		"templates/game.html",
		"templates/login.html",
		"templates/admin/index.html",
	)
}
