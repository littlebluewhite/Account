INSERT INTO `user`
VALUES (1, 'wilson', 'wilson', "123456", '2023-12-15', NULL, NULL, NULL, NULL, '2023-12-15 02:38:44', '2023-12-15 02:38:44');

INSERT INTO `workspace`
VALUES (1, 'admin', NULL, 0, '', 1, 1, '9999-02-15',
        '{\"alarm\": {\"rule\": 200}, \"account\": {\"user\": 50, \"group\": 10, \"sub_workspace\": 4}, \"schedule\": {\"schedule\": 4, \"task_template\": 6, \"time_template\": 10, \"command_template\": 20}, \"node_object\": {\"node\": 100, \"object\": 1000}}',
        '{}', '{}', '{}', '2023-12-15 02:39:54', '2023-12-15 02:39:54');

INSERT INTO `w_user`
VALUES (1, 1, 1, 1,
        '{\"const\": {\"user\": 15, \"group\": 15, \"default_auth\": {\"custom\": 15, \"system\": 15}, \"sub_workspace\": {\"control\": 15, \"is_join\": true}}, \"custom\": {}, \"system\": {}}',
        '2023-12-15 02:40:25', '2023-12-15 02:40:25');

INSERT INTO `default_auth` VALUES (1,'const','{\"default_auth\": {\"pass_down\": 0,\"custom\": 0},\"user\": 0,\"group\":0,\"sub_workspace\":{\"is_join\": false,\"control\": 0}}');