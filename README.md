# 🃏 Card Reveal Challenge

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/SQLite-CGO--Free-003B57?style=for-the-badge&logo=sqlite&logoColor=white" alt="SQLite">
  <img src="https://img.shields.io/badge/TailwindCSS-v3-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white" alt="Tailwind CSS">
  <img src="https://img.shields.io/badge/Web%20Audio-API%20Synth-F59E0B?style=for-the-badge" alt="Web Audio API">
  <img src="https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-4F46E5?style=for-the-badge" alt="Cross Platform">
</p>

<p align="center">
  <strong>Game Web Interaktif Kelas & Acara Perpisahan OSIS dengan Animasi 3D Flip Card, Efek Suara Sintetis, dan Panel Admin Terproteksi.</strong>
</p>

<p align="center">
  <a href="#-fitur-unggulan">Fitur Unggulan</a> •
  <a href="#-tampilan--estetika-kartu-as">Desain Kartu AS</a> •
  <a href="#-cara-bermain-di-kelas">Cara Bermain</a> •
  <a href="#-format-batch-import">Batch Import</a> •
  <a href="#-panduan-instalasi">Instalasi</a> •
  <a href="#-arsitektur--api">Arsitektur & API</a>
</p>

---

## 🌟 Sekilas Tentang Proyek

**Card Reveal Challenge** adalah aplikasi web ringan berperforma tinggi yang dirancang khusus untuk kegiatan kelas, *ice breaking*, maupun acara perpisahan seperti **Lengseran OSIS SMK NIBA**.

Setiap peserta memilih satu nomor kartu misterius dari grid interaktif. Saat kartu diklik, kartu akan berputar dengan animasi **3D Flip 360° yang realistis**, memutar efek suara *whoosh* dan *fanfare*, meluncurkan hujan konfeti, lalu memunculkan modal berukuran raksasa yang menampilkan kalimat tantangan/sifat siswa agar dapat dibaca bersama-sama oleh seluruh ruangan melalui proyektor.

---

## ✨ Fitur Unggulan

### 🎮 1. Papan Permainan Interaktif (Student Board)
- **Animasi 3D Card Flip**: Efek pembalikan kartu realistis berbasis CSS 3D Transforms (`perspective: 1200px` & kurva fisika `cubic-bezier`).
- **Web Audio API Sound Effects (Zero External Files)**: Suara *whoosh* pembalikan dan melodi perayaan (*chord fanfare*) disintesis langsung secara matematis di browser tanpa butuh file MP3/WAV eksternal.
- **Efek Konfeti Canvas**: Partikel konfeti warna-warni yang meledak otomatis saat kartu terbuka.
- **Focus Modal untuk Layar Proyektor**: Teks tantangan otomatis disajikan dalam jendela pop-up besar dengan kontras tinggi sehingga terbaca jelas dari baris belakang kelas.
- **Status Kartu Terpakai**: Kartu yang sudah terbuka ditandai dengan badge khusus dan sedikit meredup agar tidak dipilih ganda.
- **Smooth Staggered Reset**: Guru dapat mengulang sesi ronde baru dalam hitungan detik dengan animasi putar balik kartu bertahap.

---

### 🔐 2. Panel Admin Terproteksi & Manajemen Lengkap
- **Sistem Autentikasi Sesi**: Akses dashboard guru dilindungi verifikasi kredensial berbasis *HTTP-Only Cookie* dan *HMAC-SHA256 Token*.
- **CRUD Kartu Mandiri**: Tambah, edit teks, ganti nomor urut, pilih simbol kartu, atau nonaktifkan visibilitas kartu tertentu.
- **📥 Mesin Batch Import (Buat Puluhan Kartu Sekaligus)**:
  - Masukkan 20–30 kartu dalam satu detik tanpa perlu input manual satu per satu.
  - Mendukung input teks biasa (1 baris per tantangan), CSV, atau pemisah *pipe* (`|`).
  - Dilengkapi tombol download template CSV dan tombol contoh cepat (*quick autofill*).
  - Pilihan mode: **Append** (tambahkan ke kartu yang ada) atau **Replace All** (ganti total semua kartu).
- **🗑️ Bulk Delete (Hapus Massal)**:
  - Dilengkapi kotak centang per baris dan tombol **Select All**.
  - Tombol aksi dinamis *Hapus Terpilih (X Kartu)* yang otomatis aktif saat ada pilihan.
  - Tombol darurat *Hapus Semua Kartu* dengan dialog konfirmasi aman.
- **🔄 Reset Default 20 Soal**: Tombol pemulihan instan untuk mengembalikan seluruh kartu ke 20 soal resmi Lengseran OSIS SMK NIBA.

---

## 🎨 Tampilan & Estetika Kartu AS

Desain mengusung filosofi **Clean, Modern, Minimalis, & Ceria**:

| Simbol | Nama Motif | Aksen Warna | Nuansa / Fase |
| :---: | :---: | :---: | :--- |
| **♠** | **Spades (Sekop)** | Midnight Slate (`#1E293B`) | Tegas, Jenaka, Problem Solver |
| **♥** | **Hearts (Hati)** | Ruby Rose (`#E11D48`) | Penuh Kasih, Hangat, Curahan Hati |
| **♦** | **Diamonds (Wajik)** | Crimson Amber (`#F43F5E`) | Menarik, Cerdas, Fashionable |
| **♣** | **Clubs (Keriting)** | Charcoal Forest (`#0F172A`) | Unik, Seru, Penuh Perjuangan |
| **★** | **Star (Bintang Emas)** | Warm Yellow (`#EAB308`) | Momen Spesial & Mengesankan |

- **Cover Belakang**: Motif kartu remi geometris minimalis warna Royal Blue dengan efek *hover lift*.
- **Muka Depan**: Latar putih bersih dengan border aksen sesuai warna lambang kartu.

---

## 🎲 Cara Bermain di Kelas / Acara

```text
[Layar Utama Proyektor] ➔ http://localhost:8080
```

1. **Pembukaan**: Guru atau MC membuka halaman utama game di depan kelas.
2. **Pilih Nomor**: Siswa yang mendapat giliran memilih satu nomor kartu yang masih tertutup.
3. **Pembalikan Kartu**: Klik kartu yang dipilih. Kartu akan membalik 3D disertai efek suara dan hujan konfeti.
4. **Membaca Tantangan**: Teks muncul di layar besar. Contoh:  
   > *"Salaman dengan si paling susah di-chat di grup, tapi kalau diajak jajan langsung ada di depan!"*
5. **Eksekusi Aksi**: Siswa mencari teman yang paling sesuai dengan sifat tersebut dan bersalaman di depan kelas!
6. **Lanjutkan Ronde**: Tutup pop-up modal, kartu akan ditandai selesai, dan giliran berpindah ke siswa berikutnya.
7. **Ronde Baru**: Klik tombol **"Reset Sesi"** untuk memainkan kembali kartu dari awal.

---

## 📥 Format Batch Import Kartu

Admin dapat mengimpor data sekaligus melalui menu **Batch Import (Banyak)**. Tiga format berikut didukung secara otomatis:

### Opsi 1: Format 1 Baris per Tantangan *(Paling Direkomendasikan)*
> Sistem otomatis memberikan nomor urut 1, 2, 3... dan mengacak simbol kartu (Sekop, Hati, Wajik, Keriting, Bintang).

```text
Salaman dengan si paling susah di-chat di grup tapi kalau diajak ngopi langsung nongol.
Salaman dengan si paling jago ngeles kalau lagi ditagih uang kas kelas.
Salaman dengan si paling receh yang gampang banget ketawa.
Salaman dengan teman yang paling bakal kamu kangenin pas lulus nanti.
```

### Opsi 2: Format CSV Lengkap
```csv
card_number,suit,challenge_text
1,spades,Salaman dengan si paling susah di-chat di grup.
2,diamonds,Salaman dengan si paling fashionable anak bisnis.
3,hearts,Salaman dengan si paling kamu sayang sebagai sahabat.
```

### Opsi 3: Format Simbol & Teks (Pemisah Pipe `|`)
```text
spades | Salaman dengan si paling sering panic attack pas proker.
hearts | Salaman dengan si paling pendengar yang baik tempat curhat.
star | Salaman dengan si paling berkembang pesat di OSIS.
```

> 💡 **Unduh Template CSV**: Anda dapat mengunduh file [template_kartu_tantangan.csv](template_kartu_tantangan.csv) yang dapat dibuka dan diedit langsung di Microsoft Excel atau Google Sheets.

---

## 🛠️ Tech Stack & Keunggulan Arsitektur

```
┌────────────────────────────────────────────────────────┐
│                      Client Browser                    │
│      Tailwind CSS • CSS 3D Preserves • Web Audio Synth │
└───────────────────────────▲────────────────────────────┘
                            │ HTTP / JSON
┌───────────────────────────▼────────────────────────────┐
│                    Go Standard HTTP Mux                │
│    Session Auth Middleware • Static File Serving • SSR │
└───────────────────────────▲────────────────────────────┘
                            │ SQL Queries
┌───────────────────────────▼────────────────────────────┐
│              SQLite Database (modernc.org/sqlite)      │
│           CGO-Free Pure Go • Zero Driver Dependency    │
└────────────────────────────────────────────────────────┘
```

- **Go (Golang 1.22+)**: Menggunakan `http.NewServeMux()` modern dengan penanganan method (`GET`, `POST`) dan path values (`{id}`) tanpa dependensi framework luar.
- **SQLite (CGO-Free `modernc.org/sqlite`)**: Berjalan mulus di Windows (Laragon/XAMPP) tanpa perlu instalasi GCC/MinGW.
- **Zero Heavy JS Frameworks**: Menggunakan JavaScript Vanilla murni sehingga waktu muat halaman mendekati 0 milidetik.

---

## 📁 Struktur Direktori

```text
revealCard/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point HTTP server & rute
├── internal/
│   ├── auth/
│   │   └── auth.go                 # Autentikasi sesi cookie & middleware
│   ├── database/
│   │   └── db.go                   # Setup SQLite, migrasi, & seeder 20 kartu OSIS
│   ├── handler/
│   │   ├── admin_handler.go        # CRUD, Batch Import, Bulk Delete
│   │   ├── auth_handler.go         # Login & Logout handler
│   │   └── game_handler.go         # Game board & API reveal/reset
│   ├── model/
│   │   └── card.go                 # Entitas Card & helper styling simbol
│   └── repository/
│       └── card_repo.go            # Operasi database transaksi & batch
├── web/
│   ├── static/
│   │   ├── css/
│   │   │   └── style.css           # Animasi 3D flip, cover remi, efek glow
│   │   └── js/
│   │       └── game.js             # Web Audio SFX sintetis & logika reveal
│   └── templates/
│       ├── admin/
│       │   └── index.html          # Dashboard kelola kartu, batch modal, bulk delete
│       ├── game.html               # Papan permainan kartu interaktif siswa
│       └── login.html              # Halaman login admin modern & ceria
├── data/
│   └── revealcard.db               # File database SQLite lokal (auto-created)
├── template_kartu_tantangan.csv     # Contoh template pengisian kartu Excel
├── soal.md                         # Naskah mentah 20 kartu tantangan OSIS
├── AGENTS.md                       # Panduan arsitektur & konvensi kode
├── README.md                       # Dokumentasi resmi proyek
├── go.mod                          # Definisi modul Go
└── revealCard.exe                  # Binary executable Windows siap pakai
```

---

## 🚀 Panduan Instalasi & Menjalankan

### Prasyarat
- [Go 1.22+](https://go.dev/dl/) terinstal di komputer.
- Git (opsional).

### 1. Jalankan Langsung dari Source Code
```bash
# 1. Pindah ke direktori proyek
cd c:\laragon\www\revealCard

# 2. Sinkronkan pustaka Go
go mod tidy

# 3. Jalankan server
go run cmd/server/main.go
```

### 2. Atau Jalankan File Binary (.exe)
```bash
.\revealCard.exe
```

### 3. Buka di Peramban Web
- 🎮 **Papan Permainan Siswa**: [http://localhost:8080](http://localhost:8080)
- 🔐 **Panel Admin / Guru**: [http://localhost:8080/admin](http://localhost:8080/admin)

---

## 🔑 Kredensial Default & Pengaturan Lingkungan (Env)

Aplikasi dapat dikonfigurasi melalui Environment Variables:

| Variabel | Nilai Default | Keterangan |
| :--- | :---: | :--- |
| `PORT` | `8080` | Port listening server HTTP |
| `ADMIN_USER` | `admin` | Username login panel guru |
| `ADMIN_PASS` | `admin123` | Password login panel guru |
| `SESSION_SECRET` | `revealcard-super-secret-key-2026` | Kunci enkripsi token sesi HMAC |

*Contoh mengubah password sebelum menjalankan:*
```powershell
$env:ADMIN_PASS="guruHebat2026"
go run cmd/server/main.go
```

---

## 📡 Spesifikasi Rute & Kontrak API

| Method | Endpoint | Hak Akses | Deskripsi |
| :--- | :--- | :---: | :--- |
| `GET` | `/` | Publik | Halaman utama papan permainan siswa |
| `POST` | `/api/cards/{id}/reveal` | Publik | Menandai kartu terbuka & mengembalikan detail tantangan |
| `POST` | `/api/cards/reset` | Publik | Menutup kembali seluruh kartu ke status awal |
| `GET` | `/login` | Publik | Halaman form autentikasi admin |
| `POST` | `/login` | Publik | Verifikasi username & password |
| `GET` | `/logout` | Publik | Mengakhiri sesi admin & hapus cookie |
| `GET` | `/admin` | **Admin** | Dashboard kelola kartu & ringkasan statistik |
| `POST` | `/admin/cards` | **Admin** | Menambahkan satu kartu baru |
| `POST` | `/admin/cards/batch` | **Admin** | Batch import puluhan kartu dari teks/CSV |
| `GET` | `/admin/template/csv` | **Admin** | Mengunduh file template CSV untuk Excel |
| `POST` | `/admin/cards/{id}` | **Admin** | Memperbarui isi atau status kartu |
| `POST` | `/admin/cards/{id}/delete` | **Admin** | Menghapus satu kartu |
| `POST` | `/admin/cards/bulk-delete` | **Admin** | Menghapus banyak kartu yang dicentang |
| `POST` | `/admin/cards/delete-all` | **Admin** | Mengosongkan seluruh kartu dari database |
| `POST` | `/admin/cards/bulk-seed` | **Admin** | Reset database ke 20 kartu default OSIS NIBA |

---

## 💡 Tips Penggunaan di Acara / Kelas

1. **Gunakan Mode Layar Penuh (F11)** pada browser proyektor agar visual kartu tampil maksimal tanpa bilah navigasi browser.
2. **Aktifkan Audio Komputer/Speaker** agar efek suara sintetis *whoosh* saat membalik kartu dan *fanfare* terdengar jelas oleh seisi ruangan.
3. **Gunakan Tombol Esc** untuk menutup modal tantangan dengan cepat lewat keyboard saat siswa selesai membaca.

---

<p align="center">
  Dibuat dengan ❤️ untuk menceriakan momen kebersamaan dan kekeluargaan kelas/OSIS.
</p>
