package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ModData is the raw JSON exported by the Lua mod.
// We keep it as raw JSON to forward directly to the API without re-marshalling.
type ModData struct {
	ExportedAt string          `json:"exportedAt"`
	GameTime   json.RawMessage `json:"gameTime"`
	Farms      json.RawMessage `json:"farms"`
	CropPrices json.RawMessage `json:"cropPrices"`
	Contracts  json.RawMessage `json:"contracts"`
	Animals    json.RawMessage `json:"animals"`
	Workers    json.RawMessage `json:"workers"`
}

// ParseModData reads the Lua mod's JSON output from modSettings.
// Returns nil if the file doesn't exist (mod not installed).
func ParseModData(saveFolder string) (*ModData, error) {
	// modSettings lives inside the game data folder (same level as savegame1, savegame2, etc.)
	modSettingsDir := filepath.Join(saveFolder, "modSettings", "FS25_FarmSimManager")
	path := filepath.Join(modSettingsDir, "data.json")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // mod not installed, non-fatal
		}
		return nil, err
	}

	var md ModData
	if err := json.Unmarshal(data, &md); err != nil {
		return nil, err
	}
	return &md, nil
}
