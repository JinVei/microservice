CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键/2022-07-18',
  `username` varchar(64) NOT NULL COMMENT '用户名/2022-07-18',
  `password` varchar(64) NOT NULL COMMENT '密码/2022-07-18',
  `telnumber` varchar(64) DEFAULT NULL UNIQUE COMMENT '手机号码/2022-07-18',
  `email` varchar(64) DEFAULT NULL UNIQUE COMMENT '邮箱/2022-07-18',
  `salt` varchar(64) NOT NULL COMMENT '密码对应的盐值/2022-07-18',
  `gender` tinyint(3) unsigned NOT NULL COMMENT '性别/0未知/1男/2女/2022-07-18',
  `status` tinyint(3) unsigned NOT NULL COMMENT '用户状态/1启用/2永久冻结/2022-07-18',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '数据库创建时间/2022-07-18',
  `create_by` bigint(20) unsigned NOT NULL COMMENT '创建者/2022-07-18',
  `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间/2022-07-18',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '数据库修改时间',
  `last_modify_by` bigint(20) unsigned DEFAULT NULL COMMENT '最后修改者/2022-07-18',
  `last_modify_time` bigint(20) unsigned DEFAULT NULL COMMENT '最后修改时间/2022-07-18',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表/2022-07-18';

CREATE INDEX username_telnumber_email
ON `users` (username, telnumber, email);

