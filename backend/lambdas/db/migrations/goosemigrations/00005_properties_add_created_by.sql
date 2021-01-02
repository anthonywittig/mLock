-- +goose Up

ALTER TABLE properties
    ADD COLUMN created_by TEXT NOT NULL;

-- +goose Down

ALTER TABLE properties
    DROP COLUMN created_by;
