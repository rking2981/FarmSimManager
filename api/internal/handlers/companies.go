package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/farmsimcompanymanager/api/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompaniesHandler struct {
	db *pgxpool.Pool
}

func NewCompaniesHandler(db *pgxpool.Pool) *CompaniesHandler {
	return &CompaniesHandler{db: db}
}

func (h *CompaniesHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	log.Printf("companies/list: user_id=%s", userID)

	// Debug: count all companies in DB and for this user
	var total, mine int
	h.db.QueryRow(r.Context(), `SELECT COUNT(*) FROM companies`).Scan(&total)
	h.db.QueryRow(r.Context(), `SELECT COUNT(*) FROM companies WHERE user_id = $1`, userID).Scan(&mine)
	log.Printf("companies/list: total=%d mine=%d", total, mine)

	rows, err := h.db.Query(r.Context(), `
		SELECT
			c.id, c.slot_id, c.farm_name, c.map_title, c.difficulty, c.creation_date, c.updated_at,
			COALESCE(s.money, 0), COALESCE(s.loan_amount, 0),
			COALESCE(s.play_time_hours, 0), COALESCE(s.last_saved, '')
		FROM companies c
		LEFT JOIN LATERAL (
			SELECT money, loan_amount, play_time_hours, last_saved
			FROM company_snapshots
			WHERE company_id = c.id
			ORDER BY pushed_at DESC LIMIT 1
		) s ON TRUE
		WHERE c.user_id = $1
		ORDER BY c.updated_at DESC`, userID)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Company struct {
		ID           string  `json:"id"`
		SlotID       string  `json:"slotId"`
		FarmName     string  `json:"farmName"`
		MapTitle     string  `json:"mapTitle"`
		Difficulty   string  `json:"difficulty"`
		CreationDate string  `json:"creationDate"`
		UpdatedAt    string  `json:"updatedAt"`
		Money        float64 `json:"money"`
		LoanAmount   float64 `json:"loanAmount"`
		PlayTimeHours float64 `json:"playTimeHours"`
		LastSaved    string  `json:"lastSaved"`
	}

	companies := []Company{}
	for rows.Next() {
		var c Company
		if err := rows.Scan(
			&c.ID, &c.SlotID, &c.FarmName, &c.MapTitle, &c.Difficulty,
			&c.CreationDate, &c.UpdatedAt, &c.Money, &c.LoanAmount,
			&c.PlayTimeHours, &c.LastSaved,
		); err != nil {
			continue
		}
		companies = append(companies, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}

func (h *CompaniesHandler) Finances(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	companyID := extractCompanyID(r.URL.Path)

	// Verify ownership
	var exists bool
	h.db.QueryRow(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM companies WHERE id = $1 AND user_id = $2)`,
		companyID, userID,
	).Scan(&exists)
	if !exists {
		http.NotFound(w, r)
		return
	}

	rows, err := h.db.Query(r.Context(), `
		SELECT day, harvest_income, mission_income, sold_wood, sold_bales, sold_wool,
			sold_milk, sold_products, sold_animals, sold_buildings, sold_vehicles,
			property_income, field_selling, other, new_vehicles_cost, construction_cost,
			field_purchase, purchase_fuel, purchase_seeds, purchase_fertilizer, purchase_pallets,
			purchase_water, vehicle_leasing_cost, vehicle_running_cost, property_maintenance,
			wage_payment, loan_interest, new_animals_cost, production_costs,
			total_income, total_expense, net_profit
		FROM daily_finances WHERE company_id = $1 ORDER BY day`, companyID)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Day struct {
		Day                int     `json:"day"`
		HarvestIncome      float64 `json:"harvestIncome"`
		MissionIncome      float64 `json:"missionIncome"`
		SoldWood           float64 `json:"soldWood"`
		SoldBales          float64 `json:"soldBales"`
		SoldWool           float64 `json:"soldWool"`
		SoldMilk           float64 `json:"soldMilk"`
		SoldProducts       float64 `json:"soldProducts"`
		SoldAnimals        float64 `json:"soldAnimals"`
		SoldBuildings      float64 `json:"soldBuildings"`
		SoldVehicles       float64 `json:"soldVehicles"`
		PropertyIncome     float64 `json:"propertyIncome"`
		FieldSelling       float64 `json:"fieldSelling"`
		Other              float64 `json:"other"`
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
		TotalIncome        float64 `json:"totalIncome"`
		TotalExpense       float64 `json:"totalExpense"`
		NetProfit          float64 `json:"netProfit"`
	}

	days := []Day{}
	for rows.Next() {
		var d Day
		rows.Scan(&d.Day, &d.HarvestIncome, &d.MissionIncome, &d.SoldWood, &d.SoldBales,
			&d.SoldWool, &d.SoldMilk, &d.SoldProducts, &d.SoldAnimals, &d.SoldBuildings,
			&d.SoldVehicles, &d.PropertyIncome, &d.FieldSelling, &d.Other,
			&d.NewVehiclesCost, &d.ConstructionCost, &d.FieldPurchase, &d.PurchaseFuel,
			&d.PurchaseSeeds, &d.PurchaseFertilizer, &d.PurchasePallets, &d.PurchaseWater,
			&d.VehicleLeasingCost, &d.VehicleRunningCost, &d.PropertyMaintenance,
			&d.WagePayment, &d.LoanInterest, &d.NewAnimalsCost, &d.ProductionCosts,
			&d.TotalIncome, &d.TotalExpense, &d.NetProfit)
		days = append(days, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"days": days})
}

func (h *CompaniesHandler) Fields(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	companyID := extractCompanyID(r.URL.Path)

	var exists bool
	h.db.QueryRow(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM companies WHERE id = $1 AND user_id = $2)`,
		companyID, userID,
	).Scan(&exists)
	if !exists {
		http.NotFound(w, r)
		return
	}

	rows, err := h.db.Query(r.Context(), `
		SELECT field_id, fruit_type, planned_fruit, growth_state, last_growth_state,
			ground_type, spray_type, spray_level, lime_level, plow_level, weed_state,
			stone_level, owned, soil_health, status,
			needs_lime, needs_plow, needs_fertilizer, has_weeds, has_stones
		FROM fields WHERE company_id = $1 ORDER BY field_id`, companyID)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Field struct {
		ID              int    `json:"id"`
		FruitType       string `json:"fruitType"`
		PlannedFruit    string `json:"plannedFruit"`
		GrowthState     int    `json:"growthState"`
		LastGrowthState int    `json:"lastGrowthState"`
		GroundType      string `json:"groundType"`
		SprayType       string `json:"sprayType"`
		SprayLevel      int    `json:"sprayLevel"`
		LimeLevel       int    `json:"limeLevel"`
		PlowLevel       int    `json:"plowLevel"`
		WeedState       int    `json:"weedState"`
		StoneLevel      int    `json:"stoneLevel"`
		Owned           bool   `json:"owned"`
		SoilHealth      int    `json:"soilHealth"`
		Status          string `json:"status"`
		NeedsLime       bool   `json:"needsLime"`
		NeedsPlow       bool   `json:"needsPlow"`
		NeedsFertilizer bool   `json:"needsFertilizer"`
		HasWeeds        bool   `json:"hasWeeds"`
		HasStones       bool   `json:"hasStones"`
	}

	fields := []Field{}
	for rows.Next() {
		var f Field
		rows.Scan(&f.ID, &f.FruitType, &f.PlannedFruit, &f.GrowthState, &f.LastGrowthState,
			&f.GroundType, &f.SprayType, &f.SprayLevel, &f.LimeLevel, &f.PlowLevel,
			&f.WeedState, &f.StoneLevel, &f.Owned, &f.SoilHealth, &f.Status,
			&f.NeedsLime, &f.NeedsPlow, &f.NeedsFertilizer, &f.HasWeeds, &f.HasStones)
		fields = append(fields, f)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fields)
}

func (h *CompaniesHandler) Vehicles(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	companyID := extractCompanyID(r.URL.Path)

	var exists bool
	h.db.QueryRow(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM companies WHERE id = $1 AND user_id = $2)`,
		companyID, userID,
	).Scan(&exists)
	if !exists {
		http.NotFound(w, r)
		return
	}

	rows, err := h.db.Query(r.Context(), `
		SELECT unique_id, name, category, filename, store_image_path,
			is_mod, age_months, price, operating_time_hours, damage, condition
		FROM vehicles WHERE company_id = $1 ORDER BY name`, companyID)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Vehicle struct {
		UniqueID           string  `json:"uniqueId"`
		Name               string  `json:"name"`
		Category           string  `json:"category"`
		Filename           string  `json:"filename"`
		StoreImagePath     string  `json:"storeImagePath"`
		IsMod              bool    `json:"isMod"`
		AgeMonths          float64 `json:"ageMonths"`
		Price              float64 `json:"price"`
		OperatingTimeHours float64 `json:"operatingTimeHours"`
		Damage             float64 `json:"damage"`
		Condition          string  `json:"condition"`
	}

	vehicles := []Vehicle{}
	for rows.Next() {
		var v Vehicle
		rows.Scan(&v.UniqueID, &v.Name, &v.Category, &v.Filename, &v.StoreImagePath,
			&v.IsMod, &v.AgeMonths, &v.Price, &v.OperatingTimeHours, &v.Damage, &v.Condition)
		vehicles = append(vehicles, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicles)
}

func (h *CompaniesHandler) Live(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	companyID := extractCompanyID(r.URL.Path)

	var exists bool
	h.db.QueryRow(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM companies WHERE id = $1 AND user_id = $2)`,
		companyID, userID,
	).Scan(&exists)
	if !exists {
		http.NotFound(w, r)
		return
	}

	row := h.db.QueryRow(r.Context(), `
		SELECT exported_at, game_time, farms, crop_prices, contracts, animals, workers, pushed_at
		FROM mod_snapshots
		WHERE company_id = $1
		ORDER BY pushed_at DESC LIMIT 1`, companyID)

	var exportedAt, pushedAt string
	var gameTime, farms, cropPrices, contracts, animals, workers []byte
	err := row.Scan(&exportedAt, &gameTime, &farms, &cropPrices, &contracts, &animals, &workers, &pushedAt)
	if err != nil {
		// No mod data yet
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"available": false})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"available":true,"exportedAt":%q,"pushedAt":%q,"gameTime":%s,"farms":%s,"cropPrices":%s,"contracts":%s,"animals":%s,"workers":%s}`,
		exportedAt, pushedAt, gameTime, farms, cropPrices, contracts, animals, workers)
}

func extractCompanyID(path string) string {
	// /api/companies/{id}/finances → id
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}
