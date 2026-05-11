package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/farmsimcompanymanager/companion/internal/config"
	"github.com/farmsimcompanymanager/companion/internal/parser"
	"github.com/rs/cors"
)

type Server struct {
	cfg   *config.Config
	store *Store
	http  *http.Server
}

func New(cfg *config.Config, store *Store) *Server {
	return &Server{cfg: cfg, store: store}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/companies", s.requireToken(s.handleCompanies))
	mux.HandleFunc("/api/companies/", s.requireToken(s.handleCompanyRoute))
	mux.HandleFunc("/api/image", s.requireTokenOrQuery(s.handleImage))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://farmsimcompanymanager.vercel.app", "http://localhost:3000"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	s.http = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", s.cfg.Port),
		Handler: c.Handler(mux),
	}

	return s.http.ListenAndServe()
}

func (s *Server) requireToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+s.cfg.Token {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// requireTokenOrQuery allows token via Authorization header OR ?token= query param.
// Used for image endpoints where <img> tags can't set headers.
func (s *Server) requireTokenOrQuery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization") == "Bearer "+s.cfg.Token
		query := r.URL.Query().Get("token") == s.cfg.Token
		if !header && !query {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// handleCompanyRoute dispatches /api/companies/{slotId}/{resource}
func (s *Server) handleCompanyRoute(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/companies/"), "/")
	if len(parts) == 2 {
		switch parts[1] {
		case "finances":
			s.handleFinances(w, r, parts[0])
			return
		case "fields":
			s.handleFields(w, r, parts[0])
			return
		case "vehicles":
			s.handleVehicles(w, r, parts[0])
			return
		}
	}
	http.NotFound(w, r)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleCompanies(w http.ResponseWriter, r *http.Request) {
	companies := s.store.GetCompanies()
	if companies == nil {
		companies = []parser.Company{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}

func (s *Server) handleVehicles(w http.ResponseWriter, r *http.Request, slotID string) {
	vehicles, err := parser.ParseVehicles(s.cfg.SaveFolder, s.cfg.GameFolder, slotID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse vehicles: %v", err), http.StatusInternalServerError)
		return
	}
	if vehicles == nil {
		vehicles = []parser.Vehicle{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}

func (s *Server) handleFields(w http.ResponseWriter, r *http.Request, slotID string) {
	fields, err := parser.ParseFields(s.cfg.SaveFolder, slotID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse fields: %v", err), http.StatusInternalServerError)
		return
	}
	if fields == nil {
		fields = []parser.Field{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fields)
}

func (s *Server) handleFinances(w http.ResponseWriter, r *http.Request, slotID string) {
	days, stats, err := parser.ParseFinances(s.cfg.SaveFolder, slotID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse finances: %v", err), http.StatusInternalServerError)
		return
	}
	if days == nil {
		days = []parser.DailyFinances{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"days":  days,
		"stats": stats,
	})
}
