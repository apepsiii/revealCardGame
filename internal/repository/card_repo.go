package repository

import (
	"database/sql"
	"fmt"
	"revealCard/internal/model"
	"strings"
	"time"
)

// CardRepository menangani interaksi database untuk entitas Card
type CardRepository struct {
	db *sql.DB
}

// NewCardRepository membuat instance baru dari CardRepository
func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db: db}
}

// GetAllCards mengambil semua kartu, dapat difilter hanya yang aktif
func (r *CardRepository) GetAllCards(onlyActive bool) ([]model.Card, error) {
	query := `
		SELECT id, card_number, suit, challenge_text, is_revealed, is_active, created_at, updated_at
		FROM cards
	`
	if onlyActive {
		query += " WHERE is_active = 1"
	}
	query += " ORDER BY card_number ASC, id ASC"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("gagal query kartu: %w", err)
	}
	defer rows.Close()

	var cards []model.Card
	for rows.Next() {
		var c model.Card
		err := rows.Scan(
			&c.ID,
			&c.CardNumber,
			&c.Suit,
			&c.ChallengeText,
			&c.IsRevealed,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal scan kartu: %w", err)
		}
		cards = append(cards, c)
	}

	return cards, rows.Err()
}

// GetCardByID mengambil satu kartu berdasarkan ID
func (r *CardRepository) GetCardByID(id int64) (*model.Card, error) {
	query := `
		SELECT id, card_number, suit, challenge_text, is_revealed, is_active, created_at, updated_at
		FROM cards WHERE id = ?
	`
	var c model.Card
	err := r.db.QueryRow(query, id).Scan(
		&c.ID,
		&c.CardNumber,
		&c.Suit,
		&c.ChallengeText,
		&c.IsRevealed,
		&c.IsActive,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil kartu ID %d: %w", id, err)
	}

	return &c, nil
}

// RevealCard menandai kartu sebagai terbuka
func (r *CardRepository) RevealCard(id int64) error {
	query := `UPDATE cards SET is_revealed = 1, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

// ResetAllRevealed mereset status semua kartu menjadi tertutup
func (r *CardRepository) ResetAllRevealed() error {
	query := `UPDATE cards SET is_revealed = 0, updated_at = ?`
	_, err := r.db.Exec(query, time.Now())
	return err
}

// CreateCard menambahkan kartu baru
func (r *CardRepository) CreateCard(c *model.Card) error {
	now := time.Now()
	query := `
		INSERT INTO cards (card_number, suit, challenge_text, is_revealed, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.Exec(query, c.CardNumber, c.Suit, c.ChallengeText, c.IsRevealed, c.IsActive, now, now)
	if err != nil {
		return fmt.Errorf("gagal menambah kartu: %w", err)
	}

	id, err := res.LastInsertId()
	if err == nil {
		c.ID = id
	}
	return nil
}

// UpdateCard memperbarui data kartu
func (r *CardRepository) UpdateCard(c *model.Card) error {
	query := `
		UPDATE cards
		SET card_number = ?, suit = ?, challenge_text = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, c.CardNumber, c.Suit, c.ChallengeText, c.IsActive, time.Now(), c.ID)
	return err
}

// DeleteCard menghapus kartu dari database
func (r *CardRepository) DeleteCard(id int64) error {
	query := `DELETE FROM cards WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// BulkDeleteCards menghapus beberapa kartu sekaligus berdasarkan kumpulan ID
func (r *CardRepository) BulkDeleteCards(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM cards WHERE id IN (%s)", strings.Join(placeholders, ","))
	_, err := r.db.Exec(query, args...)
	return err
}

// DeleteAllCards menghapus seluruh kartu dan mereset auto increment urutan
func (r *CardRepository) DeleteAllCards() error {
	if _, err := r.db.Exec("DELETE FROM cards"); err != nil {
		return err
	}
	_, _ = r.db.Exec("DELETE FROM sqlite_sequence WHERE name='cards'")
	return nil
}

// BatchCreateCards menyimpan daftar banyak kartu sekaligus dalam satu transaksi
func (r *CardRepository) BatchCreateCards(cards []model.Card, replaceAll bool) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer tx.Rollback()

	if replaceAll {
		if _, err := tx.Exec("DELETE FROM cards"); err != nil {
			return fmt.Errorf("gagal menghapus kartu lama: %w", err)
		}
		_, _ = tx.Exec("DELETE FROM sqlite_sequence WHERE name='cards'")
	}

	stmt, err := tx.Prepare(`
		INSERT INTO cards (card_number, suit, challenge_text, is_revealed, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("gagal prepare insert statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	for _, c := range cards {
		if _, err := stmt.Exec(c.CardNumber, c.Suit, c.ChallengeText, 0, 1, now, now); err != nil {
			return fmt.Errorf("gagal insert kartu batch #%d: %w", c.CardNumber, err)
		}
	}

	return tx.Commit()
}

// GetNextCardNumber mengambil nomor urut berikutnya yang disarankan
func (r *CardRepository) GetNextCardNumber() int {
	var maxNum sql.NullInt64
	_ = r.db.QueryRow("SELECT MAX(card_number) FROM cards").Scan(&maxNum)
	if maxNum.Valid {
		return int(maxNum.Int64) + 1
	}
	return 1
}

// ResetToDefaultSeed mengosongkan dan memasukkan kembali 25 data awal
func (r *CardRepository) ResetToDefaultSeed() error {
	_, err := r.db.Exec("DELETE FROM cards")
	if err != nil {
		return err
	}
	_, err = r.db.Exec("DELETE FROM sqlite_sequence WHERE name='cards'")
	if err != nil {
		// Ignore if table doesn't have autoincrement sequence entry yet
	}

	seeds := []struct {
		CardNumber int
		Suit       string
		Challenge  string
	}{
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

	tx, err := r.db.Begin()
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

