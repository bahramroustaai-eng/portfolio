-- +goose Up
CREATE TABLE item_types (
    id Serial Primary Key,
    name Text not null unique
);


-- +goose Down
DROP TABLE item_types;
