-- +goose Up
CREATE TABLE debt_payments
(
    id          SERIAL PRIMARY KEY,
    debt_id     BIGINT      NOT NULL,
    payer_id    BIGINT      NOT NULL,
    receiver_id BIGINT      NOT NULL,
    amount      BIGINT      NOT NULL,
    note        TEXT,
    paid_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT check_different_users
        CHECK ( payer_id <> receiver_id ),
    CONSTRAINT check_positive_amount
        CHECK ( amount > 0 ),
    CONSTRAINT fk_debts_payments_debt
        FOREIGN KEY (debt_id)
            REFERENCES debts,
    CONSTRAINT fk_debts_payments_payer_id
        FOREIGN KEY (payer_id)
            REFERENCES users,
    CONSTRAINT fk_debts_payments_receiver_id
        FOREIGN KEY (receiver_id)
            REFERENCES users
);

ALTER TABLE debts
    ADD COLUMN description TEXT;

-- +goose Down
DROP TABLE IF EXISTS debt_payments;

ALTER TABLE debts
    DROP COLUMN IF EXISTS description;