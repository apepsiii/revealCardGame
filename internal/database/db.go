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
		{1, "spades", "Salaman dengan si paling susah di-chat di grup"},
		{2, "diamonds", "Salaman dengan si paling jago ngeles"},
		{3, "clubs", "Salaman dengan si paling receh"},
		{4, "hearts", "Salaman dengan si paling sering ketiduran"},
		{5, "star", "Salaman dengan si paling random isi kepalanya"},
		{6, "spades", "Salaman dengan si paling sering panic attack"},
		{7, "hearts", "Salaman dengan si paling good vibes"},
		{8, "diamonds", "Salaman dengan si paling fashionable"},
		{9, "clubs", "Salaman dengan si paling tukang makan / nyemil"},
		{10, "star", "Salaman dengan si paling sabar"},
		{11, "spades", "Salaman dengan si paling problem solver"},
		{12, "hearts", "Salaman dengan si paling pendengar yang baik"},
		{13, "diamonds", "Salaman dengan si paling diam-diam care"},
		{14, "clubs", "Salaman dengan si paling rela berkorban"},
		{15, "star", "Salaman dengan si paling berkembang pesat"},
		{16, "spades", "Salaman dengan si paling tegar"},
		{17, "diamonds", "Salaman dengan si paling berkesan"},
		{18, "hearts", "Salaman dengan si paling kamu sayang"},
		{19, "clubs", "Salaman dengan seseorang yang pernah kamu repotin"},
		{20, "star", "Salaman dengan teman yang bakal paling kamu kangenin"},
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

