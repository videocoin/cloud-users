-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE users ADD `uirole` int(11) DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE users DROP `uirole`;
