-- +goose Up

CREATE TABLE units (
	id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    property_id UUID NOT NULL,
    calendar_url TEXT,
    updated_by TEXT NOT NULL,
    CONSTRAINT fk_property FOREIGN KEY(property_id) references properties(id)
);
		
-- +goose Down

DROP TABLE IF EXISTS units;
