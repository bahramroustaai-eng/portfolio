-- +goose Up
CREATE TABLE debts (
    id SERIAL PRIMARY KEY,
    lender  INT NOT NULL,
    borrower INT NOT NULL ,
    amount INT NOT NULL ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_debts_lender
                   FOREIGN KEY (lender)
                   REFERENCES users(id),

    CONSTRAINT fk_debts_borrower
                   FOREIGN KEY (borrower)
                   REFERENCES users(id),

    CONSTRAINT check_different_users
                   CHECK ( lender <> borrower ),
    CONSTRAINT check_positive_amount
                   CHECK ( amount > 0  )
                   );

-- +goose Down
DROP TABLE debts;
