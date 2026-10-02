# 🤖 Panduan Pengembang & AI Agent (AGENTS.md)

Dokumen ini memuat arsitektur, panduan desain, standar kode, dan batasan teknis untuk siapa pun (pengembang manusia maupun AI Agent) yang bekerja pada proyek **Card Reveal Challenge**.

---

## 🎯 Ringkasan Proyek

- **Nama Proyek**: Card Reveal Challenge (`revealCard`)
- **Tujuan**: Aplikasi web edukatif & interaktif untuk kelas. Siswa memilih kartu (total ~20-30 kartu), kartu terbuka dengan animasi *flip* 3D, lalu menampilkan kalimat tantangan yang harus dibaca/dilakukan.
- **Role Pengguna**:
  1. **Siswa / Kelas**: Halaman utama permainan yang bersih, interaktif, dan responsif.
  2. **Admin / Guru**: Halaman admin untuk mengelola isi tantangan kartu (CRUD), mengatur status aktif, dan mereset kartu.

---

## 🛠️ Aturan Arsitektur & Teknologi

### 1. Backend (Go)
- **Versi**: Go 1.22+.
- **Driver SQLite**: Utamakan **`modernc.org/sqlite`** (CGO-free) agar dapat berjalan mulus di lingkungan Windows (seperti Laragon) tanpa memerlukan compiler C/GCC eksternal.
- **Routing / HTTP Server**: Gunakan router standar Go 1.22+ (`http.NewServeMux()` dengan method matching baru) atau framework ringan seperti `chi`. Hindari dependensi berlebih yang memperlambat kompilasi.
- **Templating**: Gunakan `html/template` bawaan Go untuk rendering server-side yang aman terhadap XSS.

### 2. Frontend & Styling
- **Styling**: Tailwind CSS (melalui CDN modern atau bundler ringan) dikombinasikan dengan custom CSS untuk animasi 3D.
- **Animasi Flip Kartu**:
  - Wajib menggunakan CSS 3D Transforms (`perspective`, `transform-style: preserve-3d`, `rotateY(180deg)`).
  - Waktu transisi ideal: `0.6s` hingga `0.8s` dengan kurva `cubic-bezier(0.4, 0, 0.2, 1)` untuk efek pembalikan kartu yang realistis dan memuaskan.
- **Interaktivitas**: Vanilla JavaScript murni (ringan, tanpa dependensi framework SPA berat). Menyimpan state kartu terbuka secara lokal / AJAX sinkronisasi ke server.

### 3. Database (SQLite)
- Lokasi file: `data/revealcard.db`.
- Database diinisialisasi otomatis (migrasi otomatis) saat aplikasi pertama kali dijalankan.
- Data awal (*seeding*): Sediakan default 20-30 kartu dengan tantangan edukatif dan ramah siswa saat database baru dibuat.

---

## 🎨 Pedoman Desain & Visual Identity

> **Konsep Utama**: Clean, Modern, Minimalis, Ceria, dengan Tema Kartu AS (*Ace / Playing Cards*).

### 1. Palet Warna
- **Latar Belakang Game**: Lembut dan cerah, misalnya *Soft Slate* atau *Cream Ivory* (`#F8FAFC` atau `#FDFBF7`) dengan ornamen pola halus.
- **Warna Kartu**:
  - **Tampak Belakang (Cover Card)**: Warna cerah yang menarik (misal Royal Blue `#2563EB`, Crimson Red `#DC2626`, atau Golden Amber `#F59E0B`) dengan pola kartu remi geometris minimalis dan lambang kartu besar di tengah.
  - **Tampak Depan (Isi Tantangan)**: Putih bersih (`#FFFFFF`) dengan border halus beraksen warna lambang kartu.
- **Warna Lambang Kartu (Suits)**:
  - Sekop (♠ Spades) & Keriting (♣ Clubs): Midnight Black / Charcoal (`#1E293B`)
  - Hati (♥ Hearts) & Wajik (♦ Diamonds): Ruby Red (`#E11D48`)
  - Bintang Emas (★ Special): Warm Yellow (`#EAB308`)

### 2. Tipografi
- Font modern, tegas, dan mudah dibaca dari jarak jauh (misal: *Inter*, *Plus Jakarta Sans*, atau *Outfit*).
- Teks tantangan dibuat berukuran besar (minimal `text-xl` hingga `text-2xl`) agar jelas saat ditampilkan di proyektor kelas.

### 3. Pengalaman Pengguna (UX)
- Kartu yang sudah terbuka diberi efek visual penanda (misal sedikit gelap / tanda centang hijau halus / tidak bisa diklik lagi) agar tidak dipilih dua kali dalam satu sesi.
- Tersedia tombol **"Reset Sesi"** dengan konfirmasi agar guru dapat memulai ronde permainan baru dalam hitungan detik.
- Modal tampilan penuh (*Focus Modal*) saat kartu dibuka agar seluruh kelas dapat fokus membaca teks tantangan bersama.

---

## 📋 Struktur Data & Kontrak API

### Entitas `Card`
```go
type Card struct {
    ID            int64     `json:"id"`
    CardNumber    int       `json:"card_number"`    // Nomor kartu 1 - 30
    Suit          string    `json:"suit"`           // 'hearts', 'diamonds', 'spades', 'clubs', 'star'
    ChallengeText string    `json:"challenge_text"` // Isi kalimat tantangan
    IsRevealed    bool      `json:"is_revealed"`    // Status terbuka di sesi saat ini
    IsActive      bool      `json:"is_active"`      // Status kartu aktif
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
```

### Rute HTTP yang Direncanakan
- `GET  /` : Halaman permainan siswa (Game Board).
- `POST /api/cards/{id}/reveal` : Menandai kartu sebagai terbuka (sinkronisasi state).
- `POST /api/cards/reset` : Mereset semua kartu ke status tertutup.
- `GET  /admin` : Dashboard daftar kartu dan tombol aksi.
- `GET  /admin/cards/new` : Form input kartu baru.
- `POST /admin/cards` : Simpan kartu baru.
- `GET  /admin/cards/{id}/edit` : Form edit kartu.
- `POST /admin/cards/{id}` : Update kartu.
- `POST /admin/cards/{id}/delete` : Hapus kartu.
- `POST /admin/cards/bulk-seed` : Regenerasi / reset ke 25-30 tantangan default.

---

## 🚦 Rencana Langkah Eksekusi (Implementation Roadmap)

1. **Inisialisasi Project**:
   - `go.mod` dengan module path `revealCard`.
   - Setup konfigurasi direktori (`cmd/`, `internal/`, `web/`, `data/`).
2. **Database & Model Layer**:
   - Setup koneksi SQLite (`modernc.org/sqlite`).
   - Eksekusi skema tabel `cards` dan seeder default 25 kartu tantangan seru.
3. **Repository & Service Layer**:
   - Fungsi CRUD kartu dan manipulasi status `is_revealed`.
4. **Admin UI & Handlers**:
   - Halaman daftar kartu admin dengan tabel bersih, toggle aktif, modal edit, dan form tambah.
5. **Game Board Frontend**:
   - Grid kartu 3D flip interaktif dengan CSS presisi.
   - Pop-up modal tantangan siswa.
   - Fitur reset kartu.
6. **Finishing & Testing**:
   - Verifikasi performa dan estetika kartu AS yang ceria dan modern.
