CREATE TABLE `user`
(
    `id`         int PRIMARY KEY AUTO_INCREMENT,
    `username`   varchar(255) UNIQUE NOT NULL,
    `name`       varchar(255),
    `password`   varchar(255)        NOT NULL,
    `birthday`   date,
    `email`      varchar(255),
    `phone`      varchar(255),
    `country`    varchar(255),
    `login_at`   datetime,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_at` datetime DEFAULT (now())
);

CREATE TABLE `workspace`
(
    `id`                  int PRIMARY KEY AUTO_INCREMENT,
    `name`                varchar(255) NOT NULL,
    `pre_workspace_id`    int,
    `rank`                int          NOT NULL,
    `ancient`             varchar(255) NOT NULL,
    `enable`              boolean      NOT NULL,
    `owner_id`            int          NOT NULL,
    `expiry_date`         date         NOT NULL,
    `auth`                json     DEFAULT (json_object()),
    `user_auth_const`     json     DEFAULT (json_object()),
    `user_auth_pass_down` json     DEFAULT (json_object()),
    `user_auth_custom`    json     DEFAULT (json_object()),
    `updated_at`          datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_at`          datetime DEFAULT (now())
);

CREATE TABLE `default_auth`
(
    `id`   int PRIMARY KEY AUTO_INCREMENT,
    `type` ENUM ('const', 'workspace') NOT NULL,
    `auth` json DEFAULT (json_object())
);

CREATE TABLE `w_user_group`
(
    `w_user_id`  int,
    `w_group_id` int,
    PRIMARY KEY (`w_user_id`, `w_group_id`)
);

CREATE TABLE `w_group`
(
    `id`           int PRIMARY KEY AUTO_INCREMENT,
    `name`         varchar(255) NOT NULL,
    `creator_id`   int          NOT NULL,
    `workspace_id` int          NOT NULL,
    `enable`       boolean      NOT NULL,
    `default_auth` json     DEFAULT (json_object()),
    `updated_at`   datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_at`   datetime DEFAULT (now())
);

CREATE TABLE `w_user`
(
    `id`           int PRIMARY KEY AUTO_INCREMENT,
    `user_id`      int     NOT NULL,
    `workspace_id` int     NOT NUll,
    `enable`       boolean NOT NULL,
    `auth`         json     DEFAULT (json_object()),
    `updated_at`   datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_at`   datetime DEFAULT (now())
);

ALTER TABLE `workspace`
    ADD FOREIGN KEY (`pre_workspace_id`) REFERENCES `workspace` (`id`) ON DELETE CASCADE;

ALTER TABLE `workspace`
    ADD FOREIGN KEY (`owner_id`) REFERENCES `user` (`id`);

ALTER TABLE `w_user_group`
    ADD FOREIGN KEY (`w_user_id`) REFERENCES `w_user` (`id`) ON DELETE CASCADE;

ALTER TABLE `w_user_group`
    ADD FOREIGN KEY (`w_group_id`) REFERENCES `w_group` (`id`) ON DELETE CASCADE;

ALTER TABLE `w_group`
    ADD FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`) ON DELETE CASCADE;

ALTER TABLE `w_user`
    ADD FOREIGN KEY (`user_id`) REFERENCES `user` (`id`);

ALTER TABLE `w_user`
    ADD FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`) ON DELETE CASCADE;
