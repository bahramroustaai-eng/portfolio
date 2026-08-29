-- +goose Up
CREATE TYPE risk_level as ENUM('low', 'medium', 'high');

CREATE TABLE items(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type_id INT,
    total_cost INT NOT NULL,
    unit INT NOT NULL ,
    ticker TEXT,
    affect_profit BOOLEAN DEFAULT false,
    risk_level risk_level NOT NULL DEFAULT 'low',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_items_type_id
                  FOREIGN KEY (type_id)
                  REFERENCES item_types(id)
);

-- +goose Down
DROP TABLE items;
DROP TYPE risk_level;
