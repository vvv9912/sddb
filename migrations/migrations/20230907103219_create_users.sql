-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL primary key,
    tg_id integer NOT NULL , --tg id
    status_user integer NOT NULL ,
    state_user integer NOT NULL ,
    CREATED_AT timestamp NOT NULL DEFAULT (NOW() at time zone 'UTC')
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
