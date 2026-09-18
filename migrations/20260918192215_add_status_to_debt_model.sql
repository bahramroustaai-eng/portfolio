-- +goose Up

CREATE TYPE debt_status as ENUM (
    'pending',
    'paid',
    'canceled'
    );

ALTER TABLE debts
    ADD COLUMN status debt_status NOT NULL DEFAULT 'pending';

-- +goose Down
ALTER TABLE debts
    DROP COLUMN IF EXISTS status;

DROP TYPE IF EXISTS debt_status;