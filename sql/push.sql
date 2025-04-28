CREATE TABLE `user_basic` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `created_at` DATETIME,
  `updated_at` DATETIME,
  `deleted_at` DATETIME,
  `name` VARCHAR(100),
  `password` VARCHAR(255),
  `phone` VARCHAR(20),
  `email` VARCHAR(100),
  `identity` VARCHAR(255),
  `clent_ip` VARCHAR(50),
  `clent_port` VARCHAR(20),
  `login_time` DATETIME,
  `heartbeat_time` DATETIME,
  `login_out_time` DATETIME,
  `device_info` VARCHAR(255),
  `is_online` BOOLEAN DEFAULT FALSE,
  INDEX (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `friendship` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `friend_id` BIGINT UNSIGNED NOT NULL,
  `status` VARCHAR(20) NOT NULL, -- apply / accepted / blocked / deleted
  `created_at` DATETIME,
  INDEX (`user_id`),
  INDEX (`friend_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `group` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(100) NOT NULL,
  `owner_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME,
  INDEX (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



CREATE TABLE `group_member` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `group_id` BIGINT UNSIGNED NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `role` VARCHAR(20) NOT NULL, -- admin / member
  INDEX (`group_id`),
  INDEX (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `message` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `created_at` DATETIME,
  `updated_at` DATETIME,
  `deleted_at` DATETIME,
  `type` VARCHAR(50),
  `from` VARCHAR(100), -- 可以是用户 ID，也可以是设备 ID 等标识
  `to` VARCHAR(100),   -- 可以是用户 ID 或群 ID
  `content` TEXT,
  `timestamp` BIGINT,
  `extra` TEXT,
  INDEX (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
