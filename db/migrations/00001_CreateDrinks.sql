-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE drinks (
  id UUID PRIMARY KEY NOT NULL,
  name VARCHAR(40) NOT NULL,
  price VARCHAR(10) NOT NULL,
  start TIMESTAMP,
  "end" TIMESTAMP,
  created TIMESTAMP DEFAULT now()
);

CREATE TABLE drinks_ingredients (
  drink_id UUID NOT NULL REFERENCES drinks (id),
  ingredient TEXT
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE drinks_ingredients;
DROP TABLE drinks;