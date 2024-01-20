CREATE TABLE
    `user_msg` (
        `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        `msg_id` VARCHAR(128) NOT NULL,
        `chat_id` VARCHAR(128) NOT NULL,
        `chat_name` VARCHAR(32) NOT NULL,
        `sender_id` VARCHAR(128) NOT NULL,
        `sender_name` VARCHAR(32) NOT NULL,
        `msg_type` VARCHAR(32) NOT NULL,
        `content` TEXT,
        `interact` BOOLEAN NOT NULL,
        `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`),
        UNIQUE INDEX `uniq_msg_id` (`msg_id`),
        INDEX `idx_chat_id` (`chat_id`, `created_at`),
        INDEX `idx_chat_interact` (`chat_id`, `interact`, `created_at`)
    ) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;

CREATE TABLE
    `replied_msg` (
        `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        `msg_id` VARCHAR(128) NOT NULL,
        `user_msg_id` VARCHAR(128) NOT NULL,
        `platform` VARCHAR(32) NOT NULL,
        `model` VARCHAR(32) NOT NULL,
        `content` TEXT,
        `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`),
        UNIQUE INDEX `uniq_msg_id` (`msg_id`),
        INDEX `idx_user_msg_id` (`user_msg_id`)
    ) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;