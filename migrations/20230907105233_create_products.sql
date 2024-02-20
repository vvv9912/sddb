-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
    article SERIAL primary key,
    catalog text,
    name text,
    description text, --varchar(255)
    photo_url bytea[],
    price FLOAT,
    length integer,
    width integer,
    heigth integer,
    weight integer,
    available boolean default true
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists products;
-- +goose StatementEnd
