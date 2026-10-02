package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// InitDB membuka koneksi SQLite, menjalankan migrasi skema tabel, dan mengisi seeder jika kosong
func InitDB(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("gagal membuat direktori database: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka sqlite: %w", err)
	}

	// Set connection pooling yang aman untuk SQLite
	db.SetMaxOpenConns(1)

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("gagal migrasi database: %w", err)
	}

	if err := seedDefaultCards(db); err != nil {
		return nil, fmt.Errorf("gagal seeding kartu awal: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS cards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		card_number INTEGER NOT NULL,
		suit TEXT NOT NULL DEFAULT 'hearts',
		challenge_text TEXT NOT NULL,
		is_revealed BOOLEAN NOT NULL DEFAULT 0,
		is_active BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_cards_number ON cards(card_number);
	CREATE INDEX IF NOT EXISTS idx_cards_active ON cards(is_active);
	`
	_, err := db.Exec(schema)
	return err
}

type defaultSeed struct {
	CardNumber int
	Suit       string
	Challenge  string
}

func seedDefaultCards(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM cards").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Sudah ada data
	}

	seeds := []defaultSeed{
		{1, "spades", "Salaman dengan si paling susah di-chat di grup, tapi kalau diajak jajan atau ngopi tiba-tiba udah di depan."},
		{2, "diamonds", "Salaman dengan si paling jago ngeles kalau lagi ditagih uang kas atau disuruh maju pas rapat."},
		{3, "clubs", "Salaman dengan si paling receh, yang gampang banget ngakak padahal jokes-nya garing banget."},
		{4, "hearts", "Salaman dengan si paling sering ketiduran di ruang OSIS atau di kelas dengan gaya andalannya."},
		{5, "star", "Salaman dengan si paling random isi kepalanya, tapi ide-idenya selalu sukses bikin kita mikir, \"Kok kepikiran aja sih?!\""},
		{6, "spades", "Salaman dengan si paling sering panic attack kalau proker udah deket, tapi pada akhirnya tugasnya selalu beres kok."},
		{7, "hearts", "Salaman dengan si paling good vibes, yang senyum dan cerianya nular banget ke satu ruangan."},
		{8, "diamonds", "Salaman dengan si paling fashionable, yang gayanya selalu on point dan estetik banget the real anak bisnis."},
		{9, "clubs", "Salaman dengan si paling tukang makan, yang selalu punya snack di tasnya dan pasrah aja kalau kita comot."},
		{10, "star", "Salaman dengan si paling sabar, yang tahan banget ngadepin drama, omelan guru, atau teman-teman yang susah diatur."},
		{11, "spades", "Salaman dengan si paling problem solver, yang kalau kita lagi deadlock atau ada masalah, dia selalu punya jalan keluar."},
		{12, "hearts", "Salaman dengan si paling pendengar yang baik, tempat curhat sejuta umat yang selalu nyimpen rahasia kita."},
		{13, "diamonds", "Salaman dengan si paling diam-diam care, yang nggak banyak omong tapi selalu mastiin temen-temennya udah makan atau belum."},
		{14, "clubs", "Salaman dengan si paling rela berkorban, yang sering pulang paling malam dan capek paling banyak demi mastiin acara kita sukses."},
		{15, "star", "Salaman dengan si paling berkembang pesat, yang pas awal masuk OSIS masih pemalu banget, tapi sekarang udah jadi sosok yang keren parah."},
		{16, "spades", "Salaman dengan si paling tegar, yang kita tahu bebannya berat dan sering capek, tapi nggak pernah ngeluh di depan kita."},
		{17, "diamonds", "Salaman dengan si paling berkesan, yang kehadirannya bikin masa-masa SMK dan OSIS kamu jauh lebih seru dan berwarna."},
		{18, "hearts", "Salaman dengan si paling kamu sayang (sebagai sahabat), sambil bilang: \"Maafin ya kalau selama ini aku sering ngerepotin dan banyak salah.\""},
		{19, "clubs", "Salaman dengan seseorang yang pernah kamu repotin atau sempat beda pendapat, dan bilang: \"Makasih ya udah mau kerja sama bareng sampai akhir.\""},
		{20, "star", "Salaman dengan teman yang paling bakal kamu kangenin, karena kamu sadar abis ini bakal pada sibuk kelas 12 dan lulus dari SMK NIBA."},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO cards (card_number, suit, challenge_text, is_revealed, is_active)
		VALUES (?, ?, ?, 0, 1)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, s := range seeds {
		if _, err := stmt.Exec(s.CardNumber, s.Suit, s.Challenge); err != nil {
			return err
		}
	}

	return tx.Commit()
}

