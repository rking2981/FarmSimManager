package parser

import (
	"encoding/xml"
	"os"
	"path/filepath"
)

type Field struct {
	ID               int     `json:"id"`
	FruitType        string  `json:"fruitType"`
	PlannedFruit     string  `json:"plannedFruit"`
	GrowthState      int     `json:"growthState"`
	LastGrowthState  int     `json:"lastGrowthState"`
	GroundType       string  `json:"groundType"`
	SprayType        string  `json:"sprayType"`
	SprayLevel       int     `json:"sprayLevel"`
	LimeLevel        int     `json:"limeLevel"`
	PlowLevel        int     `json:"plowLevel"`
	WeedState        int     `json:"weedState"`
	StoneLevel       int     `json:"stoneLevel"`
	Owned            bool    `json:"owned"`
	SoilHealth       int     `json:"soilHealth"` // 0-100 computed score
	Status           string  `json:"status"`     // "ready", "growing", "empty", "needs_attention"
	NeedsLime        bool    `json:"needsLime"`
	NeedsPlow        bool    `json:"needsPlow"`
	NeedsFertilizer  bool    `json:"needsFertilizer"`
	HasWeeds         bool    `json:"hasWeeds"`
	HasStones        bool    `json:"hasStones"`
}

type xmlFields struct {
	XMLName xml.Name   `xml:"fields"`
	Fields  []xmlField `xml:"field"`
}

type xmlField struct {
	ID              int    `xml:"id,attr"`
	FruitType       string `xml:"fruitType,attr"`
	PlannedFruit    string `xml:"plannedFruit,attr"`
	GrowthState     int    `xml:"growthState,attr"`
	LastGrowthState int    `xml:"lastGrowthState,attr"`
	GroundType      string `xml:"groundType,attr"`
	SprayType       string `xml:"sprayType,attr"`
	SprayLevel      int    `xml:"sprayLevel,attr"`
	LimeLevel       int    `xml:"limeLevel,attr"`
	PlowLevel       int    `xml:"plowLevel,attr"`
	WeedState       int    `xml:"weedState,attr"`
	StoneLevel      int    `xml:"stoneLevel,attr"`
}

type xmlFarmlands struct {
	XMLName   xml.Name      `xml:"farmlands"`
	Farmlands []xmlFarmland `xml:"farmland"`
}

type xmlFarmland struct {
	ID     int `xml:"id,attr"`
	FarmID int `xml:"farmId,attr"`
}

func ParseFields(saveFolder, slotID string) ([]Field, error) {
	dir := filepath.Join(saveFolder, "savegame"+slotID)

	// Parse ownership map from farmland.xml
	owned, err := parseFarmlandOwnership(filepath.Join(dir, "farmland.xml"))
	if err != nil {
		owned = map[int]bool{} // non-fatal — show all fields
	}

	data, err := os.ReadFile(filepath.Join(dir, "fields.xml"))
	if err != nil {
		return nil, err
	}

	var raw xmlFields
	if err := xml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	fields := make([]Field, 0, len(raw.Fields))
	for _, f := range raw.Fields {
		field := buildField(f, owned[f.ID])
		fields = append(fields, field)
	}
	return fields, nil
}

func parseFarmlandOwnership(path string) (map[int]bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var fl xmlFarmlands
	if err := xml.Unmarshal(data, &fl); err != nil {
		return nil, err
	}
	owned := make(map[int]bool, len(fl.Farmlands))
	for _, f := range fl.Farmlands {
		if f.FarmID != 0 {
			owned[f.ID] = true
		}
	}
	return owned, nil
}

func buildField(f xmlField, owned bool) Field {
	needsLime := f.LimeLevel == 0
	needsPlow := f.PlowLevel == 0
	needsFertilizer := f.SprayLevel == 0 && f.FruitType != "UNKNOWN" && f.FruitType != ""
	hasWeeds := f.WeedState > 0
	hasStones := f.StoneLevel > 0

	// Soil health score: 100 = perfect, deduct for each issue
	soilHealth := 100
	if needsLime {
		soilHealth -= 25
	}
	if needsPlow {
		soilHealth -= 20
	}
	if needsFertilizer {
		soilHealth -= 25
	}
	if hasWeeds {
		soilHealth -= 20
	}
	if hasStones {
		soilHealth -= 10
	}
	if soilHealth < 0 {
		soilHealth = 0
	}

	status := classifyStatus(f)

	return Field{
		ID:              f.ID,
		FruitType:       f.FruitType,
		PlannedFruit:    f.PlannedFruit,
		GrowthState:     f.GrowthState,
		LastGrowthState: f.LastGrowthState,
		GroundType:      f.GroundType,
		SprayType:       f.SprayType,
		SprayLevel:      f.SprayLevel,
		LimeLevel:       f.LimeLevel,
		PlowLevel:       f.PlowLevel,
		WeedState:       f.WeedState,
		StoneLevel:      f.StoneLevel,
		Owned:           owned,
		SoilHealth:      soilHealth,
		Status:          status,
		NeedsLime:       needsLime,
		NeedsPlow:       needsPlow,
		NeedsFertilizer: needsFertilizer,
		HasWeeds:        hasWeeds,
		HasStones:       hasStones,
	}
}

func classifyStatus(f xmlField) string {
	if f.GroundType == "HARVEST_READY" {
		return "ready"
	}
	if f.FruitType != "UNKNOWN" && f.FruitType != "" && f.GrowthState > 0 {
		return "growing"
	}
	if f.WeedState > 0 || f.StoneLevel > 0 {
		return "needs_attention"
	}
	return "empty"
}
