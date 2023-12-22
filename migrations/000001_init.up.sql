CREATE TABLE `w_user`
(
    `id`         int PRIMARY KEY AUTO_INCREMENT,
    `uid`        varchar(255) UNIQUE NOT NULL,
    `username`   varchar(255),
    `password`   varchar(255),
    `birthday`   date,
    `email`      varchar(255),
    `phone`      varchar(255),
    `country`    varchar(255),
    `created_at` datetime DEFAULT (now())
);

CREATE TABLE `workspace`
(
    `id`            int PRIMARY KEY AUTO_INCREMENT,
    `name`          varchar(255),
    `pre_workspace` int,
    `rank`          int          NOT NULL,
    `ancient`       varchar(255) NOT NULL,
    `enable`        boolean      NOT NULL,
    `owner_id`      int UNIQUE   NOT NULL,
    `expiry_date`   date         NOT NULL,
    `auth`          json     DEFAULT (json_object()),
    `user_auth`     json     DEFAULT (json_object()),
    `created_at`    datetime DEFAULT (now())
);

CREATE TABLE `default_auth`
(
    `id`   int PRIMARY KEY AUTO_INCREMENT,
    `type` ENUM ('license', 'workspace'),
    `auth` json
);


CREATE TABLE `user_group`
(
    `user_id`  int,
    `group_id` int,
    PRIMARY KEY (`user_id`, `group_id`)

);

CREATE TABLE `w_group`
(
    `id`           int PRIMARY KEY AUTO_INCREMENT,
    `name`         varchar(255) NOT NULL,
    `creator_id`   int          NOT NULL,
    `workspace_id` int          NOT NULL,
    `enable`       boolean      NOT NULL,
    `created_at`   datetime DEFAULT (now())
);

CREATE TABLE `user_workspace`
(
    `user_id`      int,
    `workspace_id` int,
    `enable`       boolean NOT NULL,
    `auth`         json     DEFAULT (json_object()),
    `created_at`   datetime DEFAULT (now()),
    PRIMARY KEY (`user_id`, `workspace_id`)
);

ALTER TABLE `workspace`
    ADD FOREIGN KEY (`pre_workspace`) REFERENCES `workspace` (`id`) ON DELETE CASCADE;

ALTER TABLE `workspace`
    ADD FOREIGN KEY (`owner_id`) REFERENCES `w_user` (`id`) ON DELETE CASCADE;

ALTER TABLE `user_group`
    ADD FOREIGN KEY (`user_id`) REFERENCES `w_user` (`id`) ON DELETE CASCADE;

ALTER TABLE `user_group`
    ADD FOREIGN KEY (`group_id`) REFERENCES `w_group` (`id`) ON DELETE CASCADE;

ALTER TABLE `w_group`
    ADD FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`) ON DELETE CASCADE;

ALTER TABLE `user_workspace`
    ADD FOREIGN KEY (`user_id`) REFERENCES `w_user` (`id`) ON DELETE CASCADE;

ALTER TABLE `user_workspace`
    ADD FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`id`) ON DELETE CASCADE;