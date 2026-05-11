package parser

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

type Vehicle struct {
	UniqueID       string  `json:"uniqueId"`
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	Filename       string  `json:"filename"`
	StoreImagePath string  `json:"storeImagePath"`
	IsMod          bool    `json:"isMod"`
	Age            float64 `json:"ageMonths"`
	Price          float64 `json:"price"`
	OperatingTime  float64 `json:"operatingTimeHours"`
	Damage         float64 `json:"damage"`
	Condition      string  `json:"condition"`
}

type xmlVehicles struct {
	XMLName  xml.Name     `xml:"vehicles"`
	Vehicles []xmlVehicle `xml:"vehicle"`
}

type xmlVehicle struct {
	Filename      string  `xml:"filename,attr"`
	ModName       string  `xml:"modName,attr"`
	UniqueID      string  `xml:"uniqueId,attr"`
	Age           float64 `xml:"age,attr"`
	Price         float64 `xml:"price,attr"`
	FarmID        int     `xml:"farmId,attr"`
	PropertyState string  `xml:"propertyState,attr"`
	OperatingTime float64 `xml:"operatingTime,attr"`
	Wearable      *struct {
		Damage float64 `xml:"damage,attr"`
	} `xml:"wearable"`
}

// xmlVehicleDef is the structure of a vehicle's own XML file (in the game/mod folder).
type xmlVehicleDef struct {
	XMLName  xml.Name `xml:"vehicle"`
	StoreData struct {
		Image string `xml:"image"`
	} `xml:"storeData"`
}

func ParseVehicles(saveFolder, gameFolder, slotID string) ([]Vehicle, error) {
	path := filepath.Join(saveFolder, "savegame"+slotID, "vehicles.xml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw xmlVehicles
	if err := xml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	modsDir := filepath.Join(filepath.Dir(saveFolder), "mods")

	vehicles := make([]Vehicle, 0)
	for _, v := range raw.Vehicles {
		if v.PropertyState != "OWNED" {
			continue
		}
		name, category := parseVehicleFilename(v.Filename)
		damage := 0.0
		if v.Wearable != nil {
			damage = v.Wearable.Damage
		}
		isMod := !strings.HasPrefix(v.Filename, "data/")
		storeImage := resolveStoreImage(v.Filename, v.ModName, gameFolder, modsDir)

		vehicles = append(vehicles, Vehicle{
			UniqueID:       v.UniqueID,
			Name:           name,
			Category:       category,
			Filename:       v.Filename,
			StoreImagePath: storeImage,
			IsMod:          isMod,
			Age:            v.Age,
			Price:          v.Price,
			OperatingTime:  v.OperatingTime / 3600,
			Damage:         damage,
			Condition:      conditionLabel(damage),
		})
	}
	return vehicles, nil
}

// resolveStoreImage reads the vehicle's own XML to extract the <image> path,
// then converts it to a relative DDS path the image endpoint can serve.
func resolveStoreImage(filename, modName, gameFolder, modsDir string) string {
	absXML := resolveVehicleXML(filename, modName, gameFolder, modsDir)
	if absXML == "" {
		return fallbackStoreImage(filename)
	}

	data, err := os.ReadFile(absXML)
	if err != nil {
		return fallbackStoreImage(filename)
	}

	var def xmlVehicleDef
	if err := xml.Unmarshal(data, &def); err != nil {
		return fallbackStoreImage(filename)
	}

	img := def.StoreData.Image
	if img == "" {
		return fallbackStoreImage(filename)
	}

	// Convert game path references to relative DDS paths:
	// "$data/vehicles/..." → "data/vehicles/..."
	// "$moddir$FS25_Mod/..." → "FS25_Mod/..."
	img = strings.ReplaceAll(img, "$data/", "data/")
	img = strings.ReplaceAll(img, "$moddir$", "")

	// Strip any leading slash
	img = strings.TrimPrefix(img, "/")

	// Replace .png extension with .dds (game XMLs reference .png but files are .dds)
	img = strings.TrimSuffix(img, ".png") + ".dds"

	return img
}

// resolveVehicleXML returns the absolute path to a vehicle's XML definition file.
func resolveVehicleXML(filename, modName, gameFolder, modsDir string) string {
	if strings.HasPrefix(filename, "data/") && gameFolder != "" {
		p := filepath.Join(gameFolder, filepath.FromSlash(filename))
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	// Mod vehicle: filename is like "$moddir$FS25_Mod/vehicles/truck.xml"
	// or just "FS25_Mod/vehicles/truck.xml" after $moddir$ stripping
	modFilename := filename
	modFilename = strings.ReplaceAll(modFilename, "$moddir$", "")
	modFilename = strings.TrimPrefix(modFilename, "/")

	if modsDir != "" {
		// Try as a path relative to mods dir
		p := filepath.Join(modsDir, filepath.FromSlash(modFilename))
		if _, err := os.Stat(p); err == nil {
			return p
		}

		// Also try inside a zip-extracted mod folder by mod name
		if modName != "" {
			p2 := filepath.Join(modsDir, modName, filepath.FromSlash(modFilename))
			if _, err := os.Stat(p2); err == nil {
				return p2
			}
		}
	}

	return ""
}

// fallbackStoreImage guesses the store image from the filename when XML reading fails.
func fallbackStoreImage(filename string) string {
	dir := filepath.ToSlash(filepath.Dir(filename))
	base := strings.TrimSuffix(filepath.Base(filename), ".xml")
	// Strip $moddir$ prefix if present
	dir = strings.ReplaceAll(dir, "$moddir$", "")
	dir = strings.TrimPrefix(dir, "/")
	return dir + "/store_" + base + ".dds"
}

func parseVehicleFilename(filename string) (string, string) {
	// Strip $moddir$ModName/ prefix for display
	f := strings.ReplaceAll(filename, "$moddir$", "")
	base := strings.TrimSuffix(filepath.Base(f), ".xml")
	name := splitCamelCase(base)
	lower := strings.ToLower(filename)
	return name, classifyVehicle(lower, name)
}

func classifyVehicle(lower, name string) string {
	switch {
	case contains(lower, "tractor"):
		return "Tractor"
	case contains(lower, "harvester", "combine", "header"):
		return "Harvester"
	case contains(lower, "trailer", "wagon"):
		return "Trailer"
	case contains(lower, "cultivator", "plow", "plough"):
		return "Tillage"
	case contains(lower, "seeder", "sower", "drill", "planter"):
		return "Seeder"
	case contains(lower, "sprayer", "spreader", "fertilizer"):
		return "Sprayer"
	case contains(lower, "loader", "frontloader", "telehandler", "forklift"):
		return "Loader"
	case contains(lower, "mower", "tedder", "rake", "baler", "wrapper"):
		return "Forage"
	case contains(lower, "truck", "pickup", "transtar", "superduty"):
		return "Truck"
	case contains(lower, "train"):
		return "Train"
	case contains(lower, "car", "suv"):
		return "Car"
	case contains(lower, "komatsu", "lumberjack", "woodchipper", "yarder", "harvester911"):
		return "Forestry"
	case contains(lower, "pump", "tank", "silo", "bunker"):
		return "Storage"
	default:
		return "Equipment"
	}
}

func contains(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func splitCamelCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, ' ')
		}
		if i == 0 && r >= 'a' && r <= 'z' {
			r -= 32
		}
		result = append(result, r)
	}
	return string(result)
}

func conditionLabel(damage float64) string {
	switch {
	case damage < 0.25:
		return "good"
	case damage < 0.60:
		return "fair"
	default:
		return "poor"
	}
}
