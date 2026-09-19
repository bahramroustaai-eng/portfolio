-- +goose Up
ALTER TABLE items
    ADD COLUMN user_id INT NOT NULL REFERENCES users(id);

CREATE INDEX idx_items_user_id_created_at
    ON items(user_id, created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_items_user_id_created_at;
ALTER TABLE items DROP COLUMN user_id;
