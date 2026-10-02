package model

import "time"

// Card merepresentasikan entitas kartu tantangan
type Card struct {
	ID            int64     `json:"id"`
	CardNumber    int       `json:"card_number"`    // Nomor urut kartu (1 - 30)
	Suit          string    `json:"suit"`           // 'hearts', 'diamonds', 'spades', 'clubs', 'star'
	ChallengeText string    `json:"challenge_text"` // Teks kalimat tantangan
	IsRevealed    bool      `json:"is_revealed"`    // Apakah sudah dibuka dalam sesi saat ini
	IsActive      bool      `json:"is_active"`      // Apakah kartu aktif
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SuitSymbol mengembalikan simbol visual kartu
func (c Card) SuitSymbol() string {
	switch c.Suit {
	case "hearts":
		return "♥"
	case "diamonds":
		return "♦"
	case "spades":
		return "♠"
	case "clubs":
		return "♣"
	case "star":
		return "★"
	default:
		return "♠"
	}
}

// SuitColorClass mengembalikan kelas warna Tailwind untuk simbol kartu
func (c Card) SuitColorClass() string {
	switch c.Suit {
	case "hearts", "diamonds":
		return "text-rose-600"
	case "spades", "clubs":
		return "text-slate-800"
	case "star":
		return "text-amber-500"
	default:
		return "text-indigo-600"
	}
}

// SuitBorderClass mengembalikan kelas warna border sesuai suit
func (c Card) SuitBorderClass() string {
	switch c.Suit {
	case "hearts", "diamonds":
		return "border-rose-400"
	case "spades", "clubs":
		return "border-slate-400"
	case "star":
		return "border-amber-400"
	default:
		return "border-indigo-400"
	}
}

