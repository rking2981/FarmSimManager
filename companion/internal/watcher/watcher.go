package watcher

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/farmsimcompanymanager/companion/internal/cloud"
	"github.com/farmsimcompanymanager/companion/internal/config"
	"github.com/farmsimcompanymanager/companion/internal/parser"
	"github.com/farmsimcompanymanager/companion/internal/server"
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	cfg    *config.Config
	store  *server.Store
	pusher *cloud.Pusher
}

func New(cfg *config.Config, store *server.Store) *Watcher {
	var pusher *cloud.Pusher
	if cfg.CloudAPIURL != "" && cfg.CloudToken != "" {
		pusher = cloud.NewPusher(cfg.CloudAPIURL, cfg.CloudToken)
		log.Printf("cloud sync enabled: %s", cfg.CloudAPIURL)
	}
	return &Watcher{cfg: cfg, store: store, pusher: pusher}
}

func (w *Watcher) Start() error {
	for {
		if _, err := os.Stat(w.cfg.SaveFolder); err == nil {
			break
		}
		log.Printf("save folder not found, retrying in 30s: %s", w.cfg.SaveFolder)
		time.Sleep(30 * time.Second)
	}

	w.reload()

	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer fw.Close()

	if err := fw.Add(w.cfg.SaveFolder); err != nil {
		return err
	}

	log.Printf("watching %s", w.cfg.SaveFolder)

	for {
		select {
		case event, ok := <-fw.Events:
			if !ok {
				return nil
			}
			if isCareerFile(event.Name) {
				log.Printf("save change detected: %s", event.Name)
				w.reload()
			}
		case err, ok := <-fw.Errors:
			if !ok {
				return nil
			}
			log.Printf("watcher error: %v", err)
		}
	}
}

func (w *Watcher) reload() {
	companies, err := parser.ParseSaveFolder(w.cfg.SaveFolder)
	if err != nil {
		log.Printf("parse error: %v", err)
		return
	}
	w.store.SetCompanies(companies)
	log.Printf("loaded %d companies", len(companies))

	if w.pusher == nil {
		return
	}

	// Read mod data once (shared across all slots — it covers all farms)
	modData, merr := parser.ParseModData(w.cfg.SaveFolder)
	if merr != nil {
		log.Printf("cloud: mod data parse error: %v", merr)
	}
	if modData != nil {
		log.Printf("cloud: mod data available (exported %s)", modData.ExportedAt)
	}

	// Push each company to the cloud
	for _, company := range companies {
		finances, _, ferr := parser.ParseFinances(w.cfg.SaveFolder, company.SlotID)
		if ferr != nil {
			log.Printf("cloud: finance parse error slot %s: %v", company.SlotID, ferr)
			continue
		}
		fields, ferr := parser.ParseFields(w.cfg.SaveFolder, company.SlotID)
		if ferr != nil {
			log.Printf("cloud: fields parse error slot %s: %v", company.SlotID, ferr)
			continue
		}
		vehicles, ferr := parser.ParseVehicles(w.cfg.SaveFolder, w.cfg.GameFolder, company.SlotID)
		if ferr != nil {
			log.Printf("cloud: vehicles parse error slot %s: %v", company.SlotID, ferr)
			continue
		}
		animals, ferr := parser.ParseAnimals(w.cfg.SaveFolder, company.SlotID)
		if ferr != nil {
			log.Printf("cloud: animals parse error slot %s: %v", company.SlotID, ferr)
			animals = []parser.Animal{}
		}
		if err := w.pusher.Push(company, finances, fields, vehicles, animals, modData); err != nil {
			log.Printf("cloud: push error slot %s: %v", company.SlotID, err)
		}
	}
}

func isCareerFile(path string) bool {
	base := filepath.Base(path)
	return strings.EqualFold(base, "careerSavegame.xml")
}
