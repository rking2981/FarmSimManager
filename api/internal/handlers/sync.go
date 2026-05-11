package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SyncHandler struct {
	db *pgxpool.Pool
}

func NewSyncHandler(db *pgxpool.Pool) *SyncHandler {
	return &SyncHandler{db: db}
}

// SyncPayload is what the companion POSTs on each save change.
type SyncPayload struct {
	Company  CompanyPayload   `json:"company"`
	Finances []FinancePayload `json:"finances"`
	Fields   []FieldPayload   `json:"fields"`
	Vehicles []VehiclePayload `json:"vehicles"`
	ModData  *ModDataPayload  `json:"modData,omitempty"`
}

type ModDataPayload struct {
	ExportedAt string          `json:"exportedAt"`
	GameTime   json.RawMessage `json:"gameTime"`
	Farms      json.RawMessage `json:"farms"`
	CropPrices json.RawMessage `json:"cropPrices"`
	Contracts  json.RawMessage `json:"contracts"`
	Animals    json.RawMessage `json:"animals"`
	Workers    json.RawMessage `json:"workers"`
}

type CompanyPayload struct {
	SlotID       string  `json:"slotId"`
	FarmName     string  `json:"farmName"`
	MapTitle     string  `json:"mapTitle"`
	Difficulty   string  `json:"difficulty"`
	CreationDate string  `json:"creationDate"`
	Money        float64 `json:"money"`
	LoanAmount   float64 `json:"loanAmount"`
	PlayTime     float64 `json:"playTimeHours"`
	LastSaved    string  `json:"lastSaved"`
}

type FinancePayload struct {
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

type FieldPayload struct {
	ID              int     `json:"id"`
	FruitType       string  `json:"fruitType"`
	PlannedFruit    string  `json:"plannedFruit"`
	GrowthState     int     `json:"growthState"`
	LastGrowthState int     `json:"lastGrowthState"`
	GroundType      string  `json:"groundType"`
	SprayType       string  `json:"sprayType"`
	SprayLevel      int     `json:"sprayLevel"`
	LimeLevel       int     `json:"limeLevel"`
	PlowLevel       int     `json:"plowLevel"`
	WeedState       int     `json:"weedState"`
	StoneLevel      int     `json:"stoneLevel"`
	Owned           bool    `json:"owned"`
	SoilHealth      int     `json:"soilHealth"`
	Status          string  `json:"status"`
	NeedsLime       bool    `json:"needsLime"`
	NeedsPlow       bool    `json:"needsPlow"`
	NeedsFertilizer bool    `json:"needsFertilizer"`
	HasWeeds        bool    `json:"hasWeeds"`
	HasStones       bool    `json:"hasStones"`
}

type VehiclePayload struct {
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

// Sync handles POST /sync from the companion app.
// Auth: companion token in Authorization header.
func (h *SyncHandler) Sync(w http.ResponseWriter, r *http.Request) {
	// Validate companion token
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var userID string
	err := h.db.QueryRow(r.Context(),
		`UPDATE companions SET last_seen = NOW() WHERE token = $1 RETURNING user_id`,
		token,
	).Scan(&userID)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("sync: user_id=%s token=%s", userID, token)

	var payload SyncPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	now := time.Now()

	// Upsert company
	var companyID string
	err = h.db.QueryRow(ctx, `
		INSERT INTO companies (user_id, slot_id, farm_name, map_title, difficulty, creation_date, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, slot_id) DO UPDATE SET
			farm_name     = EXCLUDED.farm_name,
			map_title     = EXCLUDED.map_title,
			difficulty    = EXCLUDED.difficulty,
			creation_date = EXCLUDED.creation_date,
			updated_at    = EXCLUDED.updated_at
		RETURNING id`,
		userID, payload.Company.SlotID, payload.Company.FarmName,
		payload.Company.MapTitle, payload.Company.Difficulty,
		payload.Company.CreationDate, now,
	).Scan(&companyID)
	if err != nil {
		log.Printf("sync: company upsert error: %v", err)
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Upsert snapshot
	_, err = h.db.Exec(ctx, `
		INSERT INTO company_snapshots (company_id, money, loan_amount, play_time_hours, last_saved, pushed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT DO NOTHING`,
		companyID, payload.Company.Money, payload.Company.LoanAmount,
		payload.Company.PlayTime, payload.Company.LastSaved, now,
	)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Upsert finances
	for _, f := range payload.Finances {
		_, err = h.db.Exec(ctx, `
			INSERT INTO daily_finances (
				company_id, day,
				harvest_income, mission_income, sold_wood, sold_bales, sold_wool,
				sold_milk, sold_products, sold_animals, sold_buildings, sold_vehicles,
				property_income, field_selling, other,
				new_vehicles_cost, construction_cost, field_purchase,
				purchase_fuel, purchase_seeds, purchase_fertilizer, purchase_pallets, purchase_water,
				vehicle_leasing_cost, vehicle_running_cost, property_maintenance,
				wage_payment, loan_interest, new_animals_cost, production_costs,
				total_income, total_expense, net_profit, pushed_at
			) VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,
				$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34
			)
			ON CONFLICT (company_id, day) DO UPDATE SET
				harvest_income=$3, mission_income=$4, sold_wood=$5, sold_bales=$6,
				sold_wool=$7, sold_milk=$8, sold_products=$9, sold_animals=$10,
				sold_buildings=$11, sold_vehicles=$12, property_income=$13, field_selling=$14,
				other=$15, new_vehicles_cost=$16, construction_cost=$17, field_purchase=$18,
				purchase_fuel=$19, purchase_seeds=$20, purchase_fertilizer=$21,
				purchase_pallets=$22, purchase_water=$23, vehicle_leasing_cost=$24,
				vehicle_running_cost=$25, property_maintenance=$26, wage_payment=$27,
				loan_interest=$28, new_animals_cost=$29, production_costs=$30,
				total_income=$31, total_expense=$32, net_profit=$33, pushed_at=$34`,
			companyID, f.Day,
			f.HarvestIncome, f.MissionIncome, f.SoldWood, f.SoldBales, f.SoldWool,
			f.SoldMilk, f.SoldProducts, f.SoldAnimals, f.SoldBuildings, f.SoldVehicles,
			f.PropertyIncome, f.FieldSelling, f.Other,
			f.NewVehiclesCost, f.ConstructionCost, f.FieldPurchase,
			f.PurchaseFuel, f.PurchaseSeeds, f.PurchaseFertilizer, f.PurchasePallets, f.PurchaseWater,
			f.VehicleLeasingCost, f.VehicleRunningCost, f.PropertyMaintenance,
			f.WagePayment, f.LoanInterest, f.NewAnimalsCost, f.ProductionCosts,
			f.TotalIncome, f.TotalExpense, f.NetProfit, now,
		)
		if err != nil {
			http.Error(w, "db error (finances): "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Upsert fields
	for _, f := range payload.Fields {
		_, err = h.db.Exec(ctx, `
			INSERT INTO fields (
				company_id, field_id, fruit_type, planned_fruit, growth_state, last_growth_state,
				ground_type, spray_type, spray_level, lime_level, plow_level, weed_state,
				stone_level, owned, soil_health, status,
				needs_lime, needs_plow, needs_fertilizer, has_weeds, has_stones, pushed_at
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)
			ON CONFLICT (company_id, field_id) DO UPDATE SET
				fruit_type=$3, planned_fruit=$4, growth_state=$5, last_growth_state=$6,
				ground_type=$7, spray_type=$8, spray_level=$9, lime_level=$10, plow_level=$11,
				weed_state=$12, stone_level=$13, owned=$14, soil_health=$15, status=$16,
				needs_lime=$17, needs_plow=$18, needs_fertilizer=$19, has_weeds=$20,
				has_stones=$21, pushed_at=$22`,
			companyID, f.ID, f.FruitType, f.PlannedFruit, f.GrowthState, f.LastGrowthState,
			f.GroundType, f.SprayType, f.SprayLevel, f.LimeLevel, f.PlowLevel, f.WeedState,
			f.StoneLevel, f.Owned, f.SoilHealth, f.Status,
			f.NeedsLime, f.NeedsPlow, f.NeedsFertilizer, f.HasWeeds, f.HasStones, now,
		)
		if err != nil {
			http.Error(w, "db error (fields): "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Upsert vehicles
	for _, v := range payload.Vehicles {
		_, err = h.db.Exec(ctx, `
			INSERT INTO vehicles (
				company_id, unique_id, name, category, filename, store_image_path,
				is_mod, age_months, price, operating_time_hours, damage, condition, pushed_at
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
			ON CONFLICT (company_id, unique_id) DO UPDATE SET
				name=$3, category=$4, filename=$5, store_image_path=$6,
				is_mod=$7, age_months=$8, price=$9, operating_time_hours=$10,
				damage=$11, condition=$12, pushed_at=$13`,
			companyID, v.UniqueID, v.Name, v.Category, v.Filename, v.StoreImagePath,
			v.IsMod, v.AgeMonths, v.Price, v.OperatingTimeHours, v.Damage, v.Condition, now,
		)
		if err != nil {
			http.Error(w, "db error (vehicles): "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Upsert mod data if present
	if payload.ModData != nil {
		md := payload.ModData
		_, err = h.db.Exec(ctx, `
			INSERT INTO mod_snapshots
				(company_id, exported_at, game_time, farms, crop_prices, contracts, animals, workers, pushed_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			companyID, md.ExportedAt,
			md.GameTime, md.Farms, md.CropPrices,
			md.Contracts, md.Animals, md.Workers, now,
		)
		if err != nil {
			// Non-fatal — log and continue
			http.Error(w, "db error (mod_data): "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	log.Printf("sync: stored company_id=%s slot=%s user=%s", companyID, payload.Company.SlotID, userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "companyId": companyID})
}
