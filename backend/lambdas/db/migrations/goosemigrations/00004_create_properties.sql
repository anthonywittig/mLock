-- +goose Up

CREATE TABLE properties (
	id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
		
-- +goose Down

DROP TABLE IF EXISTS properties;
