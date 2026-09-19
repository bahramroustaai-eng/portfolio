-- +goose Up
ALTER TABLE items
    ADD COLUMN price_per_unit BIGINT NOT NULL default 0 CHECK (price_per_unit > 0);

-- +goose Down
ALTER TABLE items
    DROP COLUMN IF EXISTS price_per_unit;
