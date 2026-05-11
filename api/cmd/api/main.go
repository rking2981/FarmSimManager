package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/farmsimcompanymanager/api/internal/db"
	"github.com/farmsimcompanymanager/api/internal/handlers"
	"github.com/farmsimcompanymanager/api/internal/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	godotenv.Load() // no-op in production

	ctx := context.Background()
	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		log.Printf("migration warning: %v", err)
	}

	auth := handlers.NewAuthHandler(pool)
	sync := handlers.NewSyncHandler(pool)
	companies := handlers.NewCompaniesHandler(pool)

	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("/auth/register", method("POST", auth.Register))
	mux.HandleFunc("/auth/login", method("POST", auth.Login))
	mux.HandleFunc("/auth/me", method("GET", middleware.RequireAuth(auth.Me)))

	// Companion sync (uses companion token, not JWT)
	mux.HandleFunc("/sync", method("POST", sync.Sync))

	// Web app endpoints (JWT protected)
	mux.HandleFunc("/api/companies", middleware.RequireAuth(companies.List))
	mux.HandleFunc("/api/companies/", middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/finances"):
			companies.Finances(w, r)
		case strings.HasSuffix(path, "/fields"):
			companies.Fields(w, r)
		case strings.HasSuffix(path, "/vehicles"):
			companies.Vehicles(w, r)
		case strings.HasSuffix(path, "/live"):
			companies.Live(w, r)
		default:
			http.NotFound(w, r)
		}
	}))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*.vercel.app", "http://localhost:3000"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API listening on :%s", port)
	if err := http.ListenAndServe(":"+port, c.Handler(mux)); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func method(m string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != m && r.Method != http.MethodOptions {
			http.Error(w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}
