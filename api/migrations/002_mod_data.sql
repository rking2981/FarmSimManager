CREATE TABLE mod_snapshots (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id  UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    exported_at TIMESTAMPTZ,
    pushed_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    game_time   JSONB NOT NULL DEFAULT '{}',
    farms       JSONB NOT NULL DEFAULT '[]',
    crop_prices JSONB NOT NULL DEFAULT '[]',
    contracts   JSONB NOT NULL DEFAULT '[]',
    animals     JSONB NOT NULL DEFAULT '[]',
    workers     JSONB NOT NULL DEFAULT '[]'
);

CREATE INDEX idx_mod_snapshots_company_id ON mod_snapshots(company_id);
