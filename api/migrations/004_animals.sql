CREATE TABLE animals (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id       UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    sub_type         TEXT NOT NULL,
    name             TEXT NOT NULL,
    category         TEXT NOT NULL,
    num_animals      INT NOT NULL DEFAULT 0,
    age_months       INT NOT NULL DEFAULT 0,
    health_pct       NUMERIC NOT NULL DEFAULT 100,
    reproduction_pct NUMERIC NOT NULL DEFAULT 100,
    base_price       NUMERIC NOT NULL DEFAULT 0,
    estimated_value  NUMERIC NOT NULL DEFAULT 0,
    pushed_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, sub_type)
);

CREATE INDEX idx_animals_company_id ON animals(company_id);
