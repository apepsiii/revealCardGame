package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"revealCard/internal/auth"
	"revealCard/internal/database"
	"revealCard/internal/handler"
	"revealCard/internal/repository"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 1. Inisialisasi Database SQLite
	dbPath := filepath.Join("data", "revealcard.db")
	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Gagal inisialisasi database: %v", err)
	}
	defer db.Close()
	log.Printf("✓ Database SQLite terhubung di: %s", dbPath)

	// 2. Parsing Template HTML
	tmpl, err := template.ParseFiles(
		filepath.Join("web", "templates", "game.html"),
		filepath.Join("web", "templates", "login.html"),
		filepath.Join("web", "templates", "admin", "index.html"),
	)
	if err != nil {
		log.Fatalf("Gagal membaca template HTML: %v", err)
	}
	log.Printf("✓ Template HTML berhasil dimuat.")

	// 3. Inisialisasi Layer Repository & Handler
	cardRepo := repository.NewCardRepository(db)
	gameH := handler.NewGameHandler(cardRepo, tmpl)
	adminH := handler.NewAdminHandler(cardRepo, tmpl)
	authH := handler.NewAuthHandler(tmpl)

	// 4. Setup Router HTTP (Go 1.22+ Standard Mux)
	mux := http.NewServeMux()

	// File Statis (CSS, JS, Aset)
	staticDir := filepath.Join("web", "static")
	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))
	mux.Handle("GET /static/", fileServer)

	// Rute Permainan Siswa (Terbuka untuk kelas)
	mux.HandleFunc("GET /{$}", gameH.Index)
	mux.HandleFunc("POST /api/cards/{id}/reveal", gameH.RevealCard)
	mux.HandleFunc("POST /api/cards/reset", gameH.ResetSession)

	// Rute Autentikasi Admin
	mux.HandleFunc("GET /login", authH.ShowLogin)
	mux.HandleFunc("POST /login", authH.Login)
	mux.HandleFunc("GET /logout", authH.Logout)
	mux.HandleFunc("POST /logout", authH.Logout)

	// Rute Panel Admin (Dilindungi Auth Middleware)
	mux.HandleFunc("GET /admin", auth.MiddlewareAuth(adminH.AdminIndex))
	mux.HandleFunc("POST /admin/cards", auth.MiddlewareAuth(adminH.CreateCard))
	mux.HandleFunc("POST /admin/cards/batch", auth.MiddlewareAuth(adminH.BatchCreate))
	mux.HandleFunc("GET /admin/template/csv", auth.MiddlewareAuth(adminH.DownloadTemplateCSV))
	mux.HandleFunc("POST /admin/cards/{id}", auth.MiddlewareAuth(adminH.UpdateCard))
	mux.HandleFunc("POST /admin/cards/{id}/delete", auth.MiddlewareAuth(adminH.DeleteCard))
	mux.HandleFunc("POST /admin/cards/bulk-seed", auth.MiddlewareAuth(adminH.ResetSeed))
	mux.HandleFunc("POST /admin/cards/bulk-delete", auth.MiddlewareAuth(adminH.BulkDelete))
	mux.HandleFunc("POST /admin/cards/delete-all", auth.MiddlewareAuth(adminH.DeleteAll))

	adminUser, _ := auth.GetCredentials()
	addr := fmt.Sprintf(":%s", port)
	log.Printf("==================================================")
	log.Printf("🃏 Card Reveal Challenge berjalan di http://localhost:%s", port)
	log.Printf("🔐 Admin Login: http://localhost:%s/login", port)
	log.Printf("👤 Default Admin Username: %s (Password: admin123)", adminUser)
	log.Printf("==================================================")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server berhenti dengan error: %v", err)
	}
}
