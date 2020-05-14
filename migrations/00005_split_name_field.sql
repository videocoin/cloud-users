-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE users DROP `name`, ADD `first_name` varchar(100) DEFAULT NULL, ADD `last_name` varchar(100) DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE users ADD `name` varchar(100) DEFAULT NULL, DROP `first_name`, DROP `last_name`;
