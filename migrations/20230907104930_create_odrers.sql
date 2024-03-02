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
    UPDATE_AT timestamp NOT NULL DEFAULT (NOW() at time zone 'UTC+03')
);
--update
CREATE OR REPLACE FUNCTION update_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.UPDATE_AT = (NOW() at time zone 'UTC+03');
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER update_orders_timestamp
    BEFORE UPDATE ON orders
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists orders;
-- +goose StatementEnd
