package handler

import (
	"html/template"
	"net/http"
	"revealCard/internal/model"
	"revealCard/internal/repository"
	"strconv"
	"strings"
)

// AdminHandler menangani operasional dashboard admin / guru
type AdminHandler struct {
	repo *repository.CardRepository
	tmpl *template.Template
}

// NewAdminHandler membuat instance AdminHandler baru
func NewAdminHandler(repo *repository.CardRepository, tmpl *template.Template) *AdminHandler {
	return &AdminHandler{
		repo: repo,
		tmpl: tmpl,
	}
}

// AdminIndex menampilkan daftar seluruh kartu di admin panel
func (h *AdminHandler) AdminIndex(w http.ResponseWriter, r *http.Request) {
	cards, err := h.repo.GetAllCards(false)
	if err != nil {
		http.Error(w, "Gagal mengambil data kartu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalCards := len(cards)
	activeCards := 0
	revealedCards := 0
	for _, c := range cards {
		if c.IsActive {
			activeCards++
		}
		if c.IsRevealed {
			revealedCards++
		}
	}

	nextNumber := h.repo.GetNextCardNumber()

	data := map[string]interface{}{
		"Cards":         cards,
		"TotalCards":    totalCards,
		"ActiveCards":   activeCards,
		"RevealedCards": revealedCards,
		"NextNumber":    nextNumber,
		"Title":         "Admin - Manajemen Kartu Tantangan",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "Gagal merender template admin: "+err.Error(), http.StatusInternalServerError)
	}
}

// CreateCard memproses penambahan kartu baru
func (h *AdminHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid: "+err.Error(), http.StatusBadRequest)
		return
	}

	cardNum, _ := strconv.Atoi(r.FormValue("card_number"))
	if cardNum <= 0 {
		cardNum = h.repo.GetNextCardNumber()
	}

	suit := strings.TrimSpace(r.FormValue("suit"))
	if suit == "" {
		suit = "hearts"
	}

	challenge := strings.TrimSpace(r.FormValue("challenge_text"))
	if challenge == "" {
		http.Error(w, "Kalimat tantangan tidak boleh kosong", http.StatusBadRequest)
		return
	}

	isActive := r.FormValue("is_active") == "1" || r.FormValue("is_active") == "on" || r.FormValue("is_active") == "true"

	card := &model.Card{
		CardNumber:    cardNum,
		Suit:          suit,
		ChallengeText: challenge,
		IsActive:      isActive,
		IsRevealed:    false,
	}

	if err := h.repo.CreateCard(card); err != nil {
		http.Error(w, "Gagal menyimpan kartu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin?msg=created", http.StatusSeeOther)
}

// UpdateCard memproses pembaruan kartu
func (h *AdminHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid: "+err.Error(), http.StatusBadRequest)
		return
	}

	cardNum, _ := strconv.Atoi(r.FormValue("card_number"))
	suit := strings.TrimSpace(r.FormValue("suit"))
	challenge := strings.TrimSpace(r.FormValue("challenge_text"))
	isActive := r.FormValue("is_active") == "1" || r.FormValue("is_active") == "on" || r.FormValue("is_active") == "true"

	card := &model.Card{
		ID:            id,
		CardNumber:    cardNum,
		Suit:          suit,
		ChallengeText: challenge,
		IsActive:      isActive,
	}

	if err := h.repo.UpdateCard(card); err != nil {
		http.Error(w, "Gagal mengupdate kartu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin?msg=updated", http.StatusSeeOther)
}

// DeleteCard menghapus kartu
func (h *AdminHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteCard(id); err != nil {
		http.Error(w, "Gagal menghapus kartu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin?msg=deleted", http.StatusSeeOther)
}

// ResetSeed mengembalikan kartu ke daftar 25 tantangan default
func (h *AdminHandler) ResetSeed(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.ResetToDefaultSeed(); err != nil {
		http.Error(w, "Gagal mereset ke seeder awal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin?msg=seeded", http.StatusSeeOther)
}

// BatchCreate memproses pembuatan banyak kartu sekaligus dari input teks / CSV
func (h *AdminHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid: "+err.Error(), http.StatusBadRequest)
		return
	}

	rawText := strings.TrimSpace(r.FormValue("raw_data"))
	if rawText == "" {
		http.Error(w, "Data tantangan tidak boleh kosong", http.StatusBadRequest)
		return
	}

	replaceAll := r.FormValue("replace_all") == "1" || r.FormValue("replace_all") == "true" || r.FormValue("replace_all") == "on"

	suitsCycle := []string{"spades", "hearts", "diamonds", "clubs", "star"}
	lines := strings.Split(rawText, "\n")

	var cards []model.Card
	currentNum := 1
	if !replaceAll {
		currentNum = h.repo.GetNextCardNumber()
	}

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "nomor") || strings.HasPrefix(lower, "card_number") || strings.HasPrefix(lower, "no,") || strings.HasPrefix(lower, "suit,") {
			continue
		}

		var cardNum int
		var suit string
		var challenge string

		var parts []string
		if strings.Contains(line, "|") {
			parts = strings.Split(line, "|")
		} else if strings.Contains(line, ";") {
			parts = strings.Split(line, ";")
		} else if strings.Contains(line, "\t") {
			parts = strings.Split(line, "\t")
		} else if strings.Contains(line, ",") {
			parts = strings.Split(line, ",")
		}

		if len(parts) >= 3 {
			n, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err == nil && n > 0 {
				cardNum = n
			} else {
				cardNum = currentNum
			}
			suit = strings.ToLower(strings.TrimSpace(parts[1]))
			challenge = strings.TrimSpace(strings.Join(parts[2:], ","))
		} else if len(parts) == 2 {
			suitCandidate := strings.ToLower(strings.TrimSpace(parts[0]))
			if isSuitValid(suitCandidate) {
				suit = suitCandidate
				challenge = strings.TrimSpace(parts[1])
			} else {
				n, err := strconv.Atoi(suitCandidate)
				if err == nil && n > 0 {
					cardNum = n
				} else {
					cardNum = currentNum
				}
				challenge = strings.TrimSpace(parts[1])
			}
		} else {
			challenge = line
		}

		if challenge == "" {
			continue
		}

		if cardNum <= 0 {
			cardNum = currentNum
		}
		currentNum = cardNum + 1

		if !isSuitValid(suit) {
			suit = suitsCycle[(cardNum-1)%len(suitsCycle)]
		}

		cards = append(cards, model.Card{
			CardNumber:    cardNum,
			Suit:          suit,
			ChallengeText: challenge,
			IsActive:      true,
			IsRevealed:    false,
		})
	}

	if len(cards) == 0 {
		http.Error(w, "Tidak ada data tantangan yang valid untuk disimpan", http.StatusBadRequest)
		return
	}

	if err := h.repo.BatchCreateCards(cards, replaceAll); err != nil {
		http.Error(w, "Gagal menyimpan batch kartu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin?msg=batch_success&count="+strconv.Itoa(len(cards)), http.StatusSeeOther)
}

// DownloadTemplateCSV mengunduh file template CSV siap pakai untuk diisi di Excel
func (h *AdminHandler) DownloadTemplateCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=template_kartu_tantangan.csv")

	csvContent := `card_number,suit,challenge_text
1,spades,Sebutkan 3 hal yang paling kamu syukuri hari ini dengan senyum lebar!
2,hearts,Peragakan gerakan hewan favoritmu selama 10 detik tanpa mengeluarkan suara!
3,diamonds,Bacakan satu pantun jenaka buatanmu sendiri untuk teman di sebelah kananmu!
4,clubs,Sebutkan nama 5 ibukota negara di dunia dalam waktu 15 detik!
5,star,Berikan pujian tulus dan tepuk tangan meriah untuk teman yang duduk paling belakang!
`
	w.Write([]byte(csvContent))
}

// BulkDelete memproses penghapusan banyak kartu terpilih dari checkbox
func (h *AdminHandler) BulkDelete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form tidak valid: "+err.Error(), http.StatusBadRequest)
		return
	}

	idStrs := r.Form["card_ids"]
	if len(idStrs) == 0 {
		raw := r.FormValue("card_ids_csv")
		if raw != "" {
			idStrs = strings.Split(raw, ",")
		}
	}

	var ids []int64
	for _, s := range idStrs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		id, err := strconv.ParseInt(s, 10, 64)
		if err == nil && id > 0 {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		http.Redirect(w, r, "/admin?msg=no_selection", http.StatusSeeOther)
		return
	}

	if err := h.repo.BulkDeleteCards(ids); err != nil {
		http.Error(w, "Gagal menghapus kartu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin?msg=bulk_deleted&count="+strconv.Itoa(len(ids)), http.StatusSeeOther)
}

// DeleteAll menghapus seluruh kartu sekaligus dari database
func (h *AdminHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.DeleteAllCards(); err != nil {
		http.Error(w, "Gagal mengosongkan kartu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin?msg=all_deleted", http.StatusSeeOther)
}

func isSuitValid(s string) bool {
	switch s {
	case "hearts", "spades", "diamonds", "clubs", "star":
		return true
	default:
		return false
	}
}
