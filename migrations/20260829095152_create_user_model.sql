-- +goose Up
CREATE TABLE users (
    id serial PRIMARY KEY,
    user_name TEXT NOT NULL UNIQUE,
    password  TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE users;
