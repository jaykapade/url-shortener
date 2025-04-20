package shortener

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaykapade/url-shortener/internal/auth"
	"github.com/jaykapade/url-shortener/internal/db"
)

type ShortenerRequest struct {
	FullURL string `json:"full_url"`
}

type ShortenerHandler struct {
	DB *db.Queries
}

func (h *ShortenerHandler) CreateShortCodeHandler(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)

	var req ShortenerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a unique ID using Sonyflake
	id, err := GenerateFlakeID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Encode the ID into Base62 to create the short code
	shortCode := EncodeBase62(id)

	// Convert the UserID into pgtype.UUID
	userIDPg := pgtype.UUID{}
	err = userIDPg.Scan(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.DB.CreateLink(r.Context(), db.CreateLinkParams{
		UserID:    userIDPg,
		FullUrl:   req.FullURL,
		ShortCode: pgtype.Text{String: shortCode, Valid: true},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"short_code": shortCode})
}

func (h *ShortenerHandler) UpdateShortCodeHandler(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)

	var req ShortenerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the UserID into pgtype.UUID
	userIDPg := pgtype.UUID{}
	err = userIDPg.Scan(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Generate a unique ID using Sonyflake
	id, err := GenerateFlakeID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Encode the ID into Base62 to create the short code
	shortCode := EncodeBase62(id)

	err = h.DB.UpdateLink(r.Context(), db.UpdateLinkParams{
		FullUrl:   req.FullURL,
		ShortCode: pgtype.Text{String: shortCode, Valid: true},
		ID:        userIDPg,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"short_code": shortCode})
}

func (h *ShortenerHandler) RedirectLinkHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "short_code")
	shortCodePg := pgtype.Text{String: shortCode, Valid: true}

	link, err := h.DB.GetLinkByShortCode(r.Context(), shortCodePg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	h.DB.UpdateClickCount(r.Context(), link.ID)

	http.Redirect(w, r, link.FullUrl, http.StatusFound) // 302 redirect

}
