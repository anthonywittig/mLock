-- +goose Up

ALTER TABLE units
    ADD COLUMN calendar_url TEXT NOT NULL DEFAULT '';

-- +goose Down

ALTER TABLE units
    DROP COLUMN calendar_url;
