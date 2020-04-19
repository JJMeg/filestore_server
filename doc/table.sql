CREATE TABLE `tbl_file`(
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
  `create_at` datetime DEFAULT NOW() COMMENT '创建时间',
  `update_at` datetime DEFAULT NOW() on update current_timestamp() COMMENT '更新日期',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态 可用禁用已删除',
  `ext1` int(11) DEFAULT '0' COMMENT '备用字段1',
  `ext2` text COMMENT '备用字段2',
  PRIMARY KEY(`id`),
  UNIQUE KEY `idx_file_hash` (`file_sha1`),
  KEY `idx_status` (`status`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `tbl_user`(
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` char(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_pwd` varchar(256) NOT NULL DEFAULT '' COMMENT '用户encoded密码',
  `email` char(64) DEFAULT '' COMMENT '邮箱',
  `phone` varchar(128) DEFAULT '' COMMENT '手机号',
  `email_validated` tinyint(1) DEFAULT 0 COMMENT '邮箱是否已验证',
  `phone_validated` tinyint(1) DEFAULT 0 COMMENT '手机是否已验证',
  `signup_at` datetime DEFAULT current_timestamp COMMENT '注册时间',
  `last_active` datetime DEFAULT current_timestamp on update current_timestamp() COMMENT '最后活跃时间戳',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '账户状态 可用禁用已删除',
  `profile` text COMMENT '用户属性',
  PRIMARY KEY(`id`),
  UNIQUE KEY `idx_phone` (`phone`),
  KEY `idx_status` (`status`)
  ) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `tbl_user_token`(
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` char(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_token` char(40) NOT NULL DEFAULT '' COMMENT '用户登陆token',
  PRIMARY KEY(`id`),
  UNIQUE KEY `idx_username` (`user_name`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
