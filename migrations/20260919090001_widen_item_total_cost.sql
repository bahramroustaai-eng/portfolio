-- +goose Up
ALTER TABLE items
    ALTER COLUMN total_cost TYPE BIGINT;

-- +goose Down
ALTER TABLE items
    ALTER COLUMN total_cost TYPE INT;
