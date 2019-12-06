-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `users` (
  `id` varchar(36) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  `password` varchar(100) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `role` int(11) DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT NULL,
  `activated_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `token` varchar(255) DEFAULT NULL,
  `xxx_unrecognized` varbinary(255) DEFAULT NULL,
  `xxx_sizecache` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_users_email` (`email`),
  KEY `idx_users_token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE users;