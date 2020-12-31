-- +goose Up

UPDATE users
    SET created_by = email
    WHERE created_by IS NULL;

ALTER TABLE users ALTER COLUMN created_by DROP NOT NULL;

-- +goose Down
