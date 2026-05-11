package main

import (
	"log"

	"github.com/farmsimcompanymanager/companion/internal/config"
	"github.com/farmsimcompanymanager/companion/internal/server"
	"github.com/farmsimcompanymanager/companion/internal/tray"
	"github.com/farmsimcompanymanager/companion/internal/watcher"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf("config:      %s", config.ConfigPath())
	log.Printf("save folder: %s", cfg.SaveFolder)
	log.Printf("game folder: %s", cfg.GameFolder)
	log.Printf("token:       %s", cfg.Token)
	log.Printf("listening:   http://127.0.0.1:%d", cfg.Port)

	store := server.NewStore()

	w := watcher.New(cfg, store)
	go func() {
		if err := w.Start(); err != nil {
			log.Printf("watcher error: %v", err)
		}
	}()

	srv := server.New(cfg, store)
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("server error: %v", err)
		}
	}()

	tray.Run(cfg, srv)
}
