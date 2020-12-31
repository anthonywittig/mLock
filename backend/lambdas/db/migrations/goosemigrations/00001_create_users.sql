-- +goose Up

CREATE TABLE users (
	id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE
);
		
-- +goose Down

DROP TABLE IF EXISTS users;
