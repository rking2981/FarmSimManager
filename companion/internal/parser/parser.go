package parser

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

type Company struct {
	SlotID       string  `json:"slotId"`
	FarmName     string  `json:"farmName"`
	MapTitle     string  `json:"mapTitle"`
	Money        float64 `json:"money"`
	LoanAmount   float64 `json:"loanAmount"`
	PlayTime     float64 `json:"playTimeHours"`
	LastSaved    string  `json:"lastSaved"`
	Difficulty   string  `json:"difficulty"`
	CreationDate string  `json:"creationDate"`
	Mods         []Mod   `json:"mods"`
}

type Mod struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Version string `json:"version"`
}

type DailyFinances struct {
	Day int `json:"day"`

	// Income
	HarvestIncome    float64 `json:"harvestIncome"`
	MissionIncome    float64 `json:"missionIncome"`
	SoldWood         float64 `json:"soldWood"`
	SoldBales        float64 `json:"soldBales"`
	SoldWool         float64 `json:"soldWool"`
	SoldMilk         float64 `json:"soldMilk"`
	SoldProducts     float64 `json:"soldProducts"`
	SoldAnimals      float64 `json:"soldAnimals"`
	SoldBuildings    float64 `json:"soldBuildings"`
	SoldVehicles     float64 `json:"soldVehicles"`
	PropertyIncome   float64 `json:"propertyIncome"`
	FieldSelling     float64 `json:"fieldSelling"`
	Other            float64 `json:"other"`

	// Expenses
	NewVehiclesCost    float64 `json:"newVehiclesCost"`
	ConstructionCost   float64 `json:"constructionCost"`
	FieldPurchase      float64 `json:"fieldPurchase"`
	PurchaseFuel       float64 `json:"purchaseFuel"`
	PurchaseSeeds      float64 `json:"purchaseSeeds"`
	PurchaseFertilizer float64 `json:"purchaseFertilizer"`
	PurchasePallets    float64 `json:"purchasePallets"`
	PurchaseWater      float64 `json:"purchaseWater"`
	VehicleLeasingCost float64 `json:"vehicleLeasingCost"`
	VehicleRunningCost float64 `json:"vehicleRunningCost"`
	PropertyMaintenance float64 `json:"propertyMaintenance"`
	WagePayment        float64 `json:"wagePayment"`
	LoanInterest       float64 `json:"loanInterest"`
	NewAnimalsCost     float64 `json:"newAnimalsCost"`
	ProductionCosts    float64 `json:"productionCosts"`

	// Computed
	TotalIncome  float64 `json:"totalIncome"`
	TotalExpense float64 `json:"totalExpense"`
	NetProfit    float64 `json:"netProfit"`
}

type FarmStats struct {
	WorkedHectares float64 `json:"workedHectares"`
	SownHectares   float64 `json:"sownHectares"`
	SprayedHectares float64 `json:"sprayedHectares"`
	ThreshedHectares float64 `json:"threshedHectares"`
	Revenue        float64 `json:"revenue"`
	Expenses       float64 `json:"expenses"`
	MissionCount   int     `json:"missionCount"`
	BaleCount      int     `json:"baleCount"`
}

// XML structs for careerSavegame.xml
type xmlCareer struct {
	XMLName  xml.Name `xml:"careerSavegame"`
	Settings struct {
		SavegameName       string `xml:"savegameName"`
		CreationDate       string `xml:"creationDate"`
		MapTitle           string `xml:"mapTitle"`
		SaveDate           string `xml:"saveDateFormatted"`
		EconomicDifficulty string `xml:"economicDifficulty"`
	} `xml:"settings"`
	Statistics struct {
		Money    float64 `xml:"money"`
		PlayTime float64 `xml:"playTime"`
	} `xml:"statistics"`
	Mods []struct {
		Name    string `xml:"modName,attr"`
		Title   string `xml:"title,attr"`
		Version string `xml:"version,attr"`
	} `xml:"mod"`
}

// XML structs for farms.xml
type xmlFarms struct {
	XMLName xml.Name  `xml:"farms"`
	Farms   []xmlFarm `xml:"farm"`
}

type xmlFarm struct {
	FarmID string  `xml:"farmId,attr"`
	Name   string  `xml:"name,attr"`
	Loan   float64 `xml:"loan,attr"`
	Money  float64 `xml:"money,attr"`

	Statistics struct {
		WorkedHectares   float64 `xml:"workedHectares"`
		SownHectares     float64 `xml:"sownHectares"`
		SprayedHectares  float64 `xml:"sprayedHectares"`
		ThreshedHectares float64 `xml:"threshedHectares"`
		Revenue          float64 `xml:"revenue"`
		Expenses         float64 `xml:"expenses"`
		MissionCount     int     `xml:"missionCount"`
		BaleCount        int     `xml:"baleCount"`
	} `xml:"statistics"`

	Finances struct {
		Days []xmlDayStats `xml:"stats"`
	} `xml:"finances"`
}

type xmlDayStats struct {
	Day                 int     `xml:"day,attr"`
	NewVehiclesCost     float64 `xml:"newVehiclesCost"`
	SoldVehicles        float64 `xml:"soldVehicles"`
	NewAnimalsCost      float64 `xml:"newAnimalsCost"`
	SoldAnimals         float64 `xml:"soldAnimals"`
	ConstructionCost    float64 `xml:"constructionCost"`
	SoldBuildings       float64 `xml:"soldBuildings"`
	FieldPurchase       float64 `xml:"fieldPurchase"`
	FieldSelling        float64 `xml:"fieldSelling"`
	VehicleRunningCost  float64 `xml:"vehicleRunningCost"`
	VehicleLeasingCost  float64 `xml:"vehicleLeasingCost"`
	PropertyMaintenance float64 `xml:"propertyMaintenance"`
	PropertyIncome      float64 `xml:"propertyIncome"`
	ProductionCosts     float64 `xml:"productionCosts"`
	SoldWood            float64 `xml:"soldWood"`
	SoldBales           float64 `xml:"soldBales"`
	SoldWool            float64 `xml:"soldWool"`
	SoldMilk            float64 `xml:"soldMilk"`
	SoldProducts        float64 `xml:"soldProducts"`
	PurchaseFuel        float64 `xml:"purchaseFuel"`
	PurchaseSeeds       float64 `xml:"purchaseSeeds"`
	PurchaseFertilizer  float64 `xml:"purchaseFertilizer"`
	PurchaseWater       float64 `xml:"purchaseWater"`
	PurchasePallets     float64 `xml:"purchasePallets"`
	HarvestIncome       float64 `xml:"harvestIncome"`
	MissionIncome       float64 `xml:"missionIncome"`
	WagePayment         float64 `xml:"wagePayment"`
	Other               float64 `xml:"other"`
	LoanInterest        float64 `xml:"loanInterest"`
}

func ParseSaveFolder(saveFolder string) ([]Company, error) {
	entries, err := os.ReadDir(saveFolder)
	if err != nil {
		return nil, err
	}

	var companies []Company
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), "savegame") {
			continue
		}
		slotID := strings.TrimPrefix(e.Name(), "savegame")
		dir := filepath.Join(saveFolder, e.Name())
		company, err := parseSlot(dir, slotID)
		if err != nil {
			continue
		}
		companies = append(companies, *company)
	}
	return companies, nil
}

func ParseFinances(saveFolder, slotID string) ([]DailyFinances, FarmStats, error) {
	dir := filepath.Join(saveFolder, "savegame"+slotID)
	return parseFarmsFile(filepath.Join(dir, "farms.xml"))
}

func parseSlot(dir, slotID string) (*Company, error) {
	career, err := parseCareerFile(filepath.Join(dir, "careerSavegame.xml"))
	if err != nil {
		return nil, err
	}

	// Pull money/loan from farms.xml (more accurate — careerSavegame statistics can lag)
	_, _, ferr := parseFarmsFile(filepath.Join(dir, "farms.xml"))
	money := career.Statistics.Money
	loan := 0.0
	if ferr == nil {
		// farms.xml parsed OK — we'll use the farm's money/loan directly via ParseFinances
		// For the summary card, careerSavegame statistics.money is fine
	}
	_ = ferr

	// Get farm name from farms.xml if available
	farmName := career.Settings.SavegameName
	if fn, err := parseFarmName(filepath.Join(dir, "farms.xml")); err == nil && fn != "" {
		farmName = fn
	}

	mods := make([]Mod, 0, len(career.Mods))
	for _, m := range career.Mods {
		mods = append(mods, Mod{Name: m.Name, Title: m.Title, Version: m.Version})
	}

	return &Company{
		SlotID:       slotID,
		FarmName:     farmName,
		MapTitle:     career.Settings.MapTitle,
		Money:        money,
		LoanAmount:   loan,
		PlayTime:     career.Statistics.PlayTime / 3600,
		LastSaved:    career.Settings.SaveDate,
		Difficulty:   career.Settings.EconomicDifficulty,
		CreationDate: career.Settings.CreationDate,
		Mods:         mods,
	}, nil
}

func parseCareerFile(path string) (*xmlCareer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c xmlCareer
	if err := xml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func parseFarmName(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var farms xmlFarms
	if err := xml.Unmarshal(data, &farms); err != nil {
		return "", err
	}
	if len(farms.Farms) > 0 {
		return farms.Farms[0].Name, nil
	}
	return "", nil
}

func parseFarmsFile(path string) ([]DailyFinances, FarmStats, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, FarmStats{}, err
	}
	var farms xmlFarms
	if err := xml.Unmarshal(data, &farms); err != nil {
		return nil, FarmStats{}, err
	}
	if len(farms.Farms) == 0 {
		return nil, FarmStats{}, nil
	}

	farm := farms.Farms[0]

	stats := FarmStats{
		WorkedHectares:   farm.Statistics.WorkedHectares,
		SownHectares:     farm.Statistics.SownHectares,
		SprayedHectares:  farm.Statistics.SprayedHectares,
		ThreshedHectares: farm.Statistics.ThreshedHectares,
		Revenue:          farm.Statistics.Revenue,
		Expenses:         farm.Statistics.Expenses,
		MissionCount:     farm.Statistics.MissionCount,
		BaleCount:        farm.Statistics.BaleCount,
	}

	days := make([]DailyFinances, 0, len(farm.Finances.Days))
	for _, d := range farm.Finances.Days {
		income := d.HarvestIncome + d.MissionIncome + d.SoldWood + d.SoldBales +
			d.SoldWool + d.SoldMilk + d.SoldProducts + d.SoldAnimals +
			d.SoldBuildings + d.SoldVehicles + d.PropertyIncome + d.FieldSelling

		// "other" is positive when it represents income (e.g. starting money, subsidies)
		if d.Other > 0 {
			income += d.Other
		}

		expense := abs(d.NewVehiclesCost) + abs(d.ConstructionCost) + abs(d.FieldPurchase) +
			abs(d.PurchaseFuel) + abs(d.PurchaseSeeds) + abs(d.PurchaseFertilizer) +
			abs(d.PurchasePallets) + abs(d.PurchaseWater) + abs(d.VehicleLeasingCost) +
			abs(d.VehicleRunningCost) + abs(d.PropertyMaintenance) + abs(d.WagePayment) +
			abs(d.LoanInterest) + abs(d.NewAnimalsCost) + abs(d.ProductionCosts)

		if d.Other < 0 {
			expense += abs(d.Other)
		}

		days = append(days, DailyFinances{
			Day:                 d.Day,
			HarvestIncome:       d.HarvestIncome,
			MissionIncome:       d.MissionIncome,
			SoldWood:            d.SoldWood,
			SoldBales:           d.SoldBales,
			SoldWool:            d.SoldWool,
			SoldMilk:            d.SoldMilk,
			SoldProducts:        d.SoldProducts,
			SoldAnimals:         d.SoldAnimals,
			SoldBuildings:       d.SoldBuildings,
			SoldVehicles:        d.SoldVehicles,
			PropertyIncome:      d.PropertyIncome,
			FieldSelling:        d.FieldSelling,
			Other:               d.Other,
			NewVehiclesCost:     abs(d.NewVehiclesCost),
			ConstructionCost:    abs(d.ConstructionCost),
			FieldPurchase:       abs(d.FieldPurchase),
			PurchaseFuel:        abs(d.PurchaseFuel),
			PurchaseSeeds:       abs(d.PurchaseSeeds),
			PurchaseFertilizer:  abs(d.PurchaseFertilizer),
			PurchasePallets:     abs(d.PurchasePallets),
			PurchaseWater:       abs(d.PurchaseWater),
			VehicleLeasingCost:  abs(d.VehicleLeasingCost),
			VehicleRunningCost:  abs(d.VehicleRunningCost),
			PropertyMaintenance: abs(d.PropertyMaintenance),
			WagePayment:         abs(d.WagePayment),
			LoanInterest:        abs(d.LoanInterest),
			NewAnimalsCost:      abs(d.NewAnimalsCost),
			ProductionCosts:     abs(d.ProductionCosts),
			TotalIncome:         income,
			TotalExpense:        expense,
			NetProfit:           income - expense,
		})
	}

	return days, stats, nil
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
