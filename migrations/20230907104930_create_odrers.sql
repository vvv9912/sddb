-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id SERIAL primary key,
    tg_id integer,
    status_order integer,
    pvz jsonb,
    type_dostavka integer,
    orderr text,
    CREATED_AT timestamp NOT NULL DEFAULT (NOW() at time zone 'UTC+03'),
    READ_AT timestamp NOT NULL DEFAULT (NOW() at time zone 'UTC+03')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists orders;
-- +goose StatementEnd
