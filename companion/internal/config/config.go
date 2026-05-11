package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

type Config struct {
	SaveFolder     string `json:"saveFolder"`
	GameFolder     string `json:"gameFolder"`
	Port           int    `json:"port"`
	Token          string `json:"token"`
	CloudAPIURL    string `json:"cloudApiUrl"`
	CloudToken     string `json:"cloudToken"`
}

// documentsDir returns the real My Documents path via the Windows Shell API,
// which correctly handles OneDrive folder redirection.
func documentsDir() string {
	path, err := windows.KnownFolderPath(windows.FOLDERID_Documents, 0)
	if err != nil {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Documents")
	}
	return path
}

func defaultSaveFolder() string {
	return filepath.Join(documentsDir(), "My Games", "FarmingSimulator2025")
}

// detectGameFolder finds the FS25 install directory via the Steam registry.
func detectGameFolder() string {
	const appID = "2530670"

	steamKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer steamKey.Close()

	steamPath, _, err := steamKey.GetStringValue("SteamPath")
	if err != nil {
		return ""
	}

	vdfPath := filepath.Join(filepath.FromSlash(steamPath), "steamapps", "libraryfolders.vdf")
	data, err := os.ReadFile(vdfPath)
	if err != nil {
		return ""
	}

	// Parse VDF block by block. Each library block contains a "path" key and
	// an "apps" sub-block listing app IDs. We collect (path, hasApp) per block
	// by tracking brace depth so path and app ID stay associated correctly.
	type libraryEntry struct {
		path   string
		hasApp bool
	}

	lines := strings.Split(string(data), "\n")
	var entries []libraryEntry
	var current libraryEntry
	depth := 0
	inApps := false

	for _, raw := range lines {
		line := strings.TrimSpace(raw)

		if line == "{" {
			depth++
			continue
		}
		if line == "}" {
			if depth == 2 {
				// Closing a top-level library block
				entries = append(entries, current)
				current = libraryEntry{}
				inApps = false
			}
			if depth == 3 {
				inApps = false
			}
			depth--
			continue
		}

		// Parse key-value pairs: "key"  "value"
		parts := strings.SplitN(line, `"`, 5)
		if len(parts) < 5 {
			// Check for sub-block opener like `"apps"`
			if len(parts) == 3 && strings.TrimSpace(parts[1]) == "apps" {
				inApps = true
			}
			continue
		}
		key := parts[1]
		val := parts[3]

		if depth == 2 && key == "path" {
			current.path = filepath.FromSlash(strings.ReplaceAll(val, `\\`, `\`))
		}
		if inApps && key == appID {
			current.hasApp = true
		}
	}

	for _, e := range entries {
		if e.hasApp && e.path != "" {
			candidate := filepath.Join(e.path, "steamapps", "common", "Farming Simulator 25")
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
	}
	return ""
}

func ConfigPath() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "FarmSimCompanyManager", "config.json")
}

func Load() (*Config, error) {
	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := &Config{
				SaveFolder: defaultSaveFolder(),
				GameFolder: detectGameFolder(),
				Port:       3847,
				Token:      generateToken(),
			}
			return cfg, cfg.Save()
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	path := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func generateToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "change-me-please"
	}
	return hex.EncodeToString(b)
}
