package parser

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

type Animal struct {
	SubType        string  `json:"subType"`
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	NumAnimals     int     `json:"numAnimals"`
	AgeMonths      int     `json:"ageMonths"`
	HealthPct      float64 `json:"healthPct"`
	ReproductionPct float64 `json:"reproductionPct"`
	BasePrice      float64 `json:"basePrice"`
	EstimatedValue float64 `json:"estimatedValue"`
}

type xmlPlaceables struct {
	XMLName    xml.Name       `xml:"placeables"`
	Placeables []xmlPlaceable `xml:"placeable"`
}

type xmlPlaceable struct {
	FarmID           int              `xml:"farmId,attr"`
	HusbandryAnimals *xmlHusbandryAnimals `xml:"husbandryAnimals"`
}

type xmlHusbandryAnimals struct {
	Clusters xmlClusters `xml:"clusters"`
}

type xmlClusters struct {
	Animals []xmlAnimal `xml:"animal"`
}

type xmlAnimal struct {
	SubType      string  `xml:"subType,attr"`
	NumAnimals   int     `xml:"numAnimals,attr"`
	Age          int     `xml:"age,attr"`
	Health       float64 `xml:"health,attr"`
	Reproduction float64 `xml:"reproduction,attr"`
}

// Age-based price multipliers matching the game's dealer pricing
var ageTierMultiplier = func(ageMonths int) float64 {
	switch {
	case ageMonths < 3:
		return 0.04
	case ageMonths < 12:
		return 0.105
	default:
		return 0.235
	}
}

// Animal base prices from maps_fillTypes.xml
var animalBasePrices = map[string]float64{
	"COW_SWISS_BROWN":      5000,
	"COW_HOLSTEIN":         5000,
	"COW_LIMOUSIN":         5000,
	"COW_ANGUS":            5000,
	"COW_WATERBUFFALO":     5000,
	"COW_HIGHLAND_CATTLE":  5000,
	"SHEEP_LANDRACE":       4000,
	"SHEEP_SWISS_MOUNTAIN": 4000,
	"SHEEP_STEINSCHAF":     4000,
	"SHEEP_BLACK_WELSH":    4000,
	"PIG_LANDRACE":         3000,
	"PIG_BLACK_PIED":       3000,
	"PIG_BERKSHIRE":        3000,
}

// Animal display names
var animalNames = map[string]string{
	"COW_SWISS_BROWN":      "Brown-Swiss",
	"COW_HOLSTEIN":         "Holstein",
	"COW_LIMOUSIN":         "Limousin",
	"COW_ANGUS":            "Angus",
	"COW_WATERBUFFALO":     "Water Buffalo",
	"COW_HIGHLAND_CATTLE":  "Highland Cattle",
	"SHEEP_LANDRACE":       "Landrace of Bentheim",
	"SHEEP_SWISS_MOUNTAIN": "Swiss Black-Brown Mountain",
	"SHEEP_STEINSCHAF":     "Steinschaf",
	"SHEEP_BLACK_WELSH":    "Black Welsh Mountain",
	"PIG_LANDRACE":         "German Landrace",
	"PIG_BLACK_PIED":       "Bentheim Black Pied",
	"PIG_BERKSHIRE":        "Berkshire",
}

func animalCategory(subType string) string {
	switch {
	case strings.HasPrefix(subType, "COW_"):
		return "Cow"
	case strings.HasPrefix(subType, "SHEEP_"):
		return "Sheep"
	case strings.HasPrefix(subType, "PIG_"):
		return "Pig"
	case strings.HasPrefix(subType, "HORSE_"):
		return "Horse"
	case strings.HasPrefix(subType, "CHICKEN_"):
		return "Chicken"
	case strings.HasPrefix(subType, "GOAT_"):
		return "Goat"
	default:
		return "Animal"
	}
}

func ParseAnimals(saveFolder, slotID string) ([]Animal, error) {
	path := filepath.Join(saveFolder, "savegame"+slotID, "placeables.xml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw xmlPlaceables
	if err := xml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	// difficulty sell multiplier — default NORMAL
	diffMultiplier := 1.8

	var animals []Animal
	for _, p := range raw.Placeables {
		if p.FarmID == 0 || p.HusbandryAnimals == nil {
			continue
		}
		for _, a := range p.HusbandryAnimals.Clusters.Animals {
			if a.NumAnimals == 0 {
				continue
			}
			basePrice := animalBasePrices[a.SubType]
			mult := ageTierMultiplier(a.Age)
			estimatedValue := basePrice * mult * diffMultiplier * float64(a.NumAnimals)

			name := animalNames[a.SubType]
			if name == "" {
				name = a.SubType
			}

			// health and reproduction are stored as 0-based damage (0 = full health)
			healthPct := (1.0 - a.Health) * 100
			if healthPct < 0 {
				healthPct = 0
			}

			animals = append(animals, Animal{
				SubType:         a.SubType,
				Name:            name,
				Category:        animalCategory(a.SubType),
				NumAnimals:      a.NumAnimals,
				AgeMonths:       a.Age,
				HealthPct:       healthPct,
				ReproductionPct: (1.0 - a.Reproduction) * 100,
				BasePrice:       basePrice,
				EstimatedValue:  estimatedValue,
			})
		}
	}

	return animals, nil
}
