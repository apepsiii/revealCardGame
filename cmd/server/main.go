package main

import (
	"fmt"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"revealCard/internal/auth"
	"revealCard/internal/database"
	"revealCard/internal/handler"
	"revealCard/internal/repository"
	"revealCard/web"
	"strings"
)

func init() {
	// Pastikan MIME types terdaftar secara eksplisit agar browser tidak menolak stylesheet/javascript
	_ = mime.AddExtensionType(".css", "text/css; charset=utf-8")
	_ = mime.AddExtensionType(".js", "application/javascript; charset=utf-8")
	_ = mime.AddExtensionType(".json", "application/json; charset=utf-8")
	_ = mime.AddExtensionType(".svg", "image/svg+xml")
	_ = mime.AddExtensionType(".png", "image/png")
	_ = mime.AddExtensionType(".jpg", "image/jpeg")
	_ = mime.AddExtensionType(".jpeg", "image/jpeg")
	_ = mime.AddExtensionType(".ico", "image/x-icon")
}

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

	// 2. Parsing Template HTML (Mendukung Disk Lokal & Embedded Binary)
	tmpl, err := web.LoadTemplates()
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

	// 5. File Statis (CSS, JS, Aset) dengan Fallback Embedded FS & Jaminan Header MIME
	staticFS := web.GetStaticFS()
	fileServer := http.StripPrefix("/static/", http.FileServer(staticFS))

	mux.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, ".css") {
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		} else if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		}
		fileServer.ServeHTTP(w, r)
	})

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
	localIP := getLocalIP()
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("==================================================")
	log.Printf("🃏 Akses di Laptop ini   : http://localhost:%s", port)
	log.Printf("📱 Akses di HP / Teman   : http://%s:%s", localIP, port)
	log.Printf("🔐 Admin Login           : http://localhost:%s/login", port)
	log.Printf("👤 Default Kredensial    : %s (Password: admin123)", adminUser)
	log.Printf("==================================================")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server berhenti dengan error: %v", err)
	}
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipStr := ipnet.IP.String()
				if !strings.HasPrefix(ipStr, "169.254.") {
					return ipStr
				}
			}
		}
	}
	return "127.0.0.1"
}
