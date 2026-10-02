package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"revealCard/internal/repository"
	"strconv"
)

// GameHandler menangani request pada sisi siswa / papan permainan
type GameHandler struct {
	repo  *repository.CardRepository
	tmpl  *template.Template
}

// NewGameHandler membuat instance GameHandler baru
func NewGameHandler(repo *repository.CardRepository, tmpl *template.Template) *GameHandler {
	return &GameHandler{
		repo: repo,
		tmpl: tmpl,
	}
}

// Index merender halaman utama permainan siswa
func (h *GameHandler) Index(w http.ResponseWriter, r *http.Request) {
	cards, err := h.repo.GetAllCards(true)
	if err != nil {
		http.Error(w, "Gagal mengambil data kartu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Cards": cards,
		"Title": "Card Reveal Challenge - Kelas Interaktif",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "game.html", data); err != nil {
		http.Error(w, "Gagal merender template: "+err.Error(), http.StatusInternalServerError)
	}
}

// RevealCard API untuk menandai kartu telah dibuka
func (h *GameHandler) RevealCard(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID kartu tidak valid", http.StatusBadRequest)
		return
	}

	card, err := h.repo.GetCardByID(id)
	if err != nil || card == nil {
		http.Error(w, "Kartu tidak ditemukan", http.StatusNotFound)
		return
	}

	if err := h.repo.RevealCard(id); err != nil {
		http.Error(w, "Gagal mengupdate status kartu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"card": map[string]interface{}{
			"id":             card.ID,
			"card_number":    card.CardNumber,
			"suit":           card.Suit,
			"suit_symbol":    card.SuitSymbol(),
			"suit_color":     card.SuitColorClass(),
			"challenge_text": card.ChallengeText,
		},
	})
}

// ResetSession API untuk mereset semua kartu kembali ke posisi tertutup
func (h *GameHandler) ResetSession(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.ResetAllRevealed(); err != nil {
		http.Error(w, "Gagal mereset status kartu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Semua kartu berhasil direset ke status tertutup.",
	})
}

