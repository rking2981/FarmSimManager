package server

import (
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/xypwn/filediver/dds"
)

// handleImage serves a DDS file from the game or mods folder, converted to PNG.
// Query param: ?path=data/vehicles/johnDeere/seriesS7/store_seriesS7.dds
func (s *Server) handleImage(w http.ResponseWriter, r *http.Request) {
	rel := r.URL.Query().Get("path")
	if rel == "" {
		http.Error(w, "missing path", http.StatusBadRequest)
		return
	}

	// Sanitize: no traversal outside allowed roots
	rel = filepath.FromSlash(rel)
	if strings.Contains(rel, "..") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	absPath := s.resolveDDSPath(rel)
	if absPath == "" {
		log.Printf("image not found: %s (gameFolder=%s)", rel, s.cfg.GameFolder)
		http.NotFound(w, r)
		return
	}

	f, err := os.Open(absPath)
	if err != nil {
		log.Printf("image open error: %s: %v", absPath, err)
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	img, err := dds.Decode(f, false)
	if err != nil {
		log.Printf("image decode error: %s: %v", absPath, err)
		http.Error(w, "failed to decode DDS", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	png.Encode(w, img)
}

// resolveDDSPath looks up the DDS file in the game folder first, then the mods folder.
func (s *Server) resolveDDSPath(rel string) string {
	if s.cfg.GameFolder != "" {
		candidate := filepath.Join(s.cfg.GameFolder, rel)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// Mod vehicles: path looks like "FS25_modName/vehicles/.../store_x.dds"
	// Mods live in the save folder's parent mods directory
	if s.cfg.SaveFolder != "" {
		modsDir := filepath.Join(filepath.Dir(s.cfg.SaveFolder), "mods")
		candidate := filepath.Join(modsDir, rel)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ""
}
