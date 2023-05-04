
CREATE TABLE `comment_subject` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键/2023-04-13',
  `obj_type` bigint(20) unsigned NOT NULL COMMENT '与评论区关联的系统的类型',
  `obj_id` bigint(20) unsigned NOT NULL COMMENT '与评论区关联的系统的id',
  `like` bigint(20) DEFAULT NULL COMMENT '赞/2023-04-13',
  `dislike` bigint(20) DEFAULT NULL COMMENT '踩/2023-04-13',
  `reply_cnt` bigint(20) DEFAULT NULL COMMENT '评论数/2023-04-13',
  `state` bigint(20) unsigned NOT NULL COMMENT '状态/0启用/1删除',
  `seq` bigint(20) unsigned DEFAULT NULL COMMENT '序列号, 每次更新行时+1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '数据库创建时间',
  `create_by` bigint(20) unsigned NOT NULL COMMENT '创建者',
  `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '数据库修改时间',
  `last_modify_by` bigint(20) unsigned DEFAULT NULL COMMENT '最后修改者',
  `last_modify_time` bigint(20) unsigned DEFAULT NULL COMMENT '最后修改时间',
  PRIMARY KEY (`id`)  
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论区表/2023-04-13';
CREATE INDEX type_objid
ON `reply_comments` (obj_type, obj_id);

CREATE TABLE `comment_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键/2023-04-13',
  `subject` bigint(20) unsigned NOT NULL COMMENT '评论区id',
  `parent` bigint(20) unsigned NOT NULL COMMENT '父评/0代表根评论/2023-04-13',
  `floor` bigint(20) DEFAULT NULL COMMENT '楼层/2023-04-13',
  `userid` bigint(20) DEFAULT NULL COMMENT '用户ID/2023-04-13',
  `replyto` bigint(20) DEFAULT NULL COMMENT '回复用户ID/2023-04-13',
  `like` bigint(20) DEFAULT NULL COMMENT '赞/2023-04-13',
  `dislike` bigint(20) DEFAULT NULL COMMENT '踩/2023-04-13',
  `reply_cnt` bigint(20) DEFAULT NULL COMMENT '回复数/2023-04-13',
  `state` bigint(20) unsigned NOT NULL COMMENT '状态/0启用/1删除',
  `seq` bigint(20) unsigned DEFAULT NULL COMMENT '序列号, 每次更新行时+1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '数据库创建时间',
  `create_by` bigint(20) unsigned NOT NULL COMMENT '创建者',
  `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '数据库修改时间',
  `last_modify_by` bigint(20) unsigned DEFAULT NULL COMMENT '最后修改者',
  `last_modify_time` bigint(20) unsigned DEFAULT NULL COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE (subject, parent, floor)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论表/2023-04-13';
CREATE INDEX subject_parent_floor_createdat
ON `comment_index` (subject, parent, floor, created_at);


CREATE TABLE `comment_content` (
  `id` bigint(20) unsigned NOT NULL COMMENT '评论 Index ID/2023-04-13',
  `content` varchar(512) DEFAULT NULL COMMENT '评论内容/2023-04-13',
  `ip` varchar(20) DEFAULT NULL COMMENT 'IP/2023-04-13',
  `platform` tinyint(8) DEFAULT NULL COMMENT '发布平台/2023-04-13',
  `device` varchar(20) DEFAULT NULL COMMENT '发布设备/2023-04-13',
  `state` bigint(20) unsigned NOT NULL COMMENT '状态/0启用/1删除',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '数据库创建时间',
  `create_by` bigint(20) unsigned NOT NULL COMMENT '创建者',
  `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '数据库修改时间',
  `last_modify_by` bigint(20) unsigned DEFAULT NULL COMMENT '最后修改者',
  `last_modify_time` bigint(20) unsigned DEFAULT NULL COMMENT '最后修改时间',
  PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论内容表/2023-04-13';