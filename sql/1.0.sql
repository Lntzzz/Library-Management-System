-- MySQL 5.8
-- Database: `book_manage_system`

create database if not exists book_manage_system;
use book_manage_system;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`         VARCHAR(20) NOT NULL,
    `name`       VARCHAR(40) NOT NULL,
    `password`   VARCHAR(36) NOT NULL,
    `type`       INT(11)     NOT NULL,
    `created_at` DATETIME    NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `book`;
CREATE TABLE `book`
(
    `id`           VARCHAR(20) NOT NULL,
    `name`         VARCHAR(40) NOT NULL,
    `author`       VARCHAR(36) NOT NULL,
    `stock`        INT(11)     NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `borrow_record`;
CREATE TABLE `borrow_record`
(
    `id`          VARCHAR(20) NOT NULL,
    `user_id`     VARCHAR(40) NOT NULL,
    `book_id`     VARCHAR(36) NOT NULL,
    `status`      INT(11)     NOT NULL,
    `borrowed_at` DATETIME    NOT NULL,
    `returned_at` DATETIME    NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;