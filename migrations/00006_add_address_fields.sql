-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE users ADD `country` varchar(100) DEFAULT NULL, ADD `region` varchar(100) DEFAULT NULL,
    ADD `city` varchar(100) DEFAULT NULL, ADD `zip` varchar(100) DEFAULT NULL,
    ADD `address_1` varchar(100) DEFAULT NULL, ADD `address_2` varchar(100) DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE users DROP `country`, DROP `region`, DROP `city`, DROP `zip`, DROP `address_1`, DROP `address_2`;
