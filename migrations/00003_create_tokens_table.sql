-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `user_api_tokens` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `token` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `xxx_unrecognized` varbinary(255) DEFAULT NULL,
  `xxx_sizecache` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_api_tokens_token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE user_api_tokens;