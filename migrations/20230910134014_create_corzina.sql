-- +goose Up
-- +goose StatementBegin
CREATE TABLE corzina (
    id SERIAL primary key,
    tg_id integer NOT NULL , --tg id
    article integer NOT NULL,
    quantity integer NOT NULL,
    CREATED_AT timestamp NOT NULL DEFAULT (NOW() at time zone 'UTC') --(UTC+03)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists corzina;
-- +goose StatementEnd
