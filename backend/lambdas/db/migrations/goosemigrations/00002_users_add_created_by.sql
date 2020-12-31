-- +goose Up

ALTER TABLE users
    ADD COLUMN created_by TEXT;

-- +goose Down

ALTER TABLE users
    DROP COLUMN created_by;
