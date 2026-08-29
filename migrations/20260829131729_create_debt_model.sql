-- +goose Up
CREATE TABLE debts (
    id SERIAL PRIMARY KEY,
    lender_id  INT NOT NULL,
    borrower_id INT NOT NULL ,
    amount INT NOT NULL ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_debts_lender
                   FOREIGN KEY (lender_id)
                   REFERENCES users(id),

    CONSTRAINT fk_debts_borrower
                   FOREIGN KEY (borrower_id)
                   REFERENCES users(id),

    CONSTRAINT check_different_users
                   CHECK ( lender_id <> borrower_id ),
    CONSTRAINT check_positive_amount
                   CHECK ( amount > 0  )
                   );

-- +goose Down
DROP TABLE debts;
