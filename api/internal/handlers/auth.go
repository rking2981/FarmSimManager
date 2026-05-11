package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/farmsimcompanymanager/api/internal/auth"
	"github.com/farmsimcompanymanager/api/internal/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthHandler struct {
	db *pgxpool.Pool
}

func NewAuthHandler(db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	body.Email = strings.ToLower(strings.TrimSpace(body.Email))
	if body.Email == "" || len(body.Password) < 8 {
		http.Error(w, "email and password (min 8 chars) required", http.StatusBadRequest)
		return
	}

	hash, err := auth.HashPassword(body.Password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var userID string
	err = h.db.QueryRow(r.Context(),
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`,
		body.Email, hash,
	).Scan(&userID)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			http.Error(w, "email already registered", http.StatusConflict)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Auto-create a companion token for this user
	companionToken := uuid.New().String()
	_, err = h.db.Exec(r.Context(),
		`INSERT INTO companions (user_id, token, name) VALUES ($1, $2, 'My PC')`,
		userID, companionToken,
	)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	jwtToken, err := auth.SignToken(userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"token":          jwtToken,
		"companionToken": companionToken,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	body.Email = strings.ToLower(strings.TrimSpace(body.Email))

	var userID, hash string
	err := h.db.QueryRow(r.Context(),
		`SELECT id, password_hash FROM users WHERE email = $1`,
		body.Email,
	).Scan(&userID, &hash)
	if err != nil || !auth.CheckPassword(hash, body.Password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Get companion token
	var companionToken string
	h.db.QueryRow(r.Context(),
		`SELECT token FROM companions WHERE user_id = $1 LIMIT 1`,
		userID,
	).Scan(&companionToken)

	jwtToken, err := auth.SignToken(userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":          jwtToken,
		"companionToken": companionToken,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	var email string
	h.db.QueryRow(r.Context(), `SELECT email FROM users WHERE id = $1`, userID).Scan(&email)

	var companionToken string
	var companionLastSeen *string
	h.db.QueryRow(r.Context(),
		`SELECT token, last_seen::text FROM companions WHERE user_id = $1 LIMIT 1`,
		userID,
	).Scan(&companionToken, &companionLastSeen)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":             userID,
		"email":          email,
		"companionToken": companionToken,
		"companionLastSeen": companionLastSeen,
	})
}
