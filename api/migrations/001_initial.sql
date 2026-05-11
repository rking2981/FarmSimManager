CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE companions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT UNIQUE NOT NULL,
    name       TEXT NOT NULL DEFAULT 'My PC',
    last_seen  TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE companies (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slot_id       TEXT NOT NULL,
    farm_name     TEXT NOT NULL,
    map_title     TEXT,
    difficulty    TEXT,
    creation_date TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, slot_id)
);

CREATE TABLE company_snapshots (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id       UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    money            NUMERIC NOT NULL DEFAULT 0,
    loan_amount      NUMERIC NOT NULL DEFAULT 0,
    play_time_hours  NUMERIC NOT NULL DEFAULT 0,
    last_saved       TEXT,
    pushed_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE daily_finances (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id           UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    day                  INT NOT NULL,
    harvest_income       NUMERIC NOT NULL DEFAULT 0,
    mission_income       NUMERIC NOT NULL DEFAULT 0,
    sold_wood            NUMERIC NOT NULL DEFAULT 0,
    sold_bales           NUMERIC NOT NULL DEFAULT 0,
    sold_wool            NUMERIC NOT NULL DEFAULT 0,
    sold_milk            NUMERIC NOT NULL DEFAULT 0,
    sold_products        NUMERIC NOT NULL DEFAULT 0,
    sold_animals         NUMERIC NOT NULL DEFAULT 0,
    sold_buildings       NUMERIC NOT NULL DEFAULT 0,
    sold_vehicles        NUMERIC NOT NULL DEFAULT 0,
    property_income      NUMERIC NOT NULL DEFAULT 0,
    field_selling        NUMERIC NOT NULL DEFAULT 0,
    other                NUMERIC NOT NULL DEFAULT 0,
    new_vehicles_cost    NUMERIC NOT NULL DEFAULT 0,
    construction_cost    NUMERIC NOT NULL DEFAULT 0,
    field_purchase       NUMERIC NOT NULL DEFAULT 0,
    purchase_fuel        NUMERIC NOT NULL DEFAULT 0,
    purchase_seeds       NUMERIC NOT NULL DEFAULT 0,
    purchase_fertilizer  NUMERIC NOT NULL DEFAULT 0,
    purchase_pallets     NUMERIC NOT NULL DEFAULT 0,
    purchase_water       NUMERIC NOT NULL DEFAULT 0,
    vehicle_leasing_cost NUMERIC NOT NULL DEFAULT 0,
    vehicle_running_cost NUMERIC NOT NULL DEFAULT 0,
    property_maintenance NUMERIC NOT NULL DEFAULT 0,
    wage_payment         NUMERIC NOT NULL DEFAULT 0,
    loan_interest        NUMERIC NOT NULL DEFAULT 0,
    new_animals_cost     NUMERIC NOT NULL DEFAULT 0,
    production_costs     NUMERIC NOT NULL DEFAULT 0,
    total_income         NUMERIC NOT NULL DEFAULT 0,
    total_expense        NUMERIC NOT NULL DEFAULT 0,
    net_profit           NUMERIC NOT NULL DEFAULT 0,
    pushed_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, day)
);

CREATE TABLE fields (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id       UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    field_id         INT NOT NULL,
    fruit_type       TEXT,
    planned_fruit    TEXT,
    growth_state     INT NOT NULL DEFAULT 0,
    last_growth_state INT NOT NULL DEFAULT 0,
    ground_type      TEXT,
    spray_type       TEXT,
    spray_level      INT NOT NULL DEFAULT 0,
    lime_level       INT NOT NULL DEFAULT 0,
    plow_level       INT NOT NULL DEFAULT 0,
    weed_state       INT NOT NULL DEFAULT 0,
    stone_level      INT NOT NULL DEFAULT 0,
    owned            BOOLEAN NOT NULL DEFAULT FALSE,
    soil_health      INT NOT NULL DEFAULT 100,
    status           TEXT NOT NULL DEFAULT 'empty',
    needs_lime       BOOLEAN NOT NULL DEFAULT FALSE,
    needs_plow       BOOLEAN NOT NULL DEFAULT FALSE,
    needs_fertilizer BOOLEAN NOT NULL DEFAULT FALSE,
    has_weeds        BOOLEAN NOT NULL DEFAULT FALSE,
    has_stones       BOOLEAN NOT NULL DEFAULT FALSE,
    pushed_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, field_id)
);

CREATE TABLE vehicles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id          UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    unique_id           TEXT NOT NULL,
    name                TEXT NOT NULL,
    category            TEXT,
    filename            TEXT,
    store_image_path    TEXT,
    is_mod              BOOLEAN NOT NULL DEFAULT FALSE,
    age_months          NUMERIC NOT NULL DEFAULT 0,
    price               NUMERIC NOT NULL DEFAULT 0,
    operating_time_hours NUMERIC NOT NULL DEFAULT 0,
    damage              NUMERIC NOT NULL DEFAULT 0,
    condition           TEXT NOT NULL DEFAULT 'good',
    pushed_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, unique_id)
);

CREATE INDEX idx_companies_user_id ON companies(user_id);
CREATE INDEX idx_daily_finances_company_id ON daily_finances(company_id);
CREATE INDEX idx_fields_company_id ON fields(company_id);
CREATE INDEX idx_vehicles_company_id ON vehicles(company_id);
CREATE INDEX idx_company_snapshots_company_id ON company_snapshots(company_id);
