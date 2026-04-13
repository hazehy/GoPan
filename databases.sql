CREATE TABLE `user_basic` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `identity` char(36) NOT NULL COMMENT '用户唯一标识',
    `name` varchar(60) NOT NULL,
    `password` varchar(255) NOT NULL,
    `email` varchar(100) NOT NULL,
    `status` tinyint NOT NULL DEFAULT 1 COMMENT '1-正常 2-禁用',
    `role` tinyint NOT NULL DEFAULT 1 COMMENT '1-普通用户 2-管理员',
    `upload_permission` tinyint NOT NULL DEFAULT 1 COMMENT '上传权限 1-允许 2-禁止',
    `download_permission` tinyint NOT NULL DEFAULT 1 COMMENT '下载权限 1-允许 2-禁止',
    `share_permission` tinyint NOT NULL DEFAULT 1 COMMENT '分享权限 1-允许 2-禁止',
    `last_login_at` datetime DEFAULT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_identity` (`identity`),
    UNIQUE KEY `uk_user_email` (`email`),
    KEY `idx_user_role` (`role`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `repository_pool` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `identity` char(36) NOT NULL COMMENT '文件唯一标识',
    `hash` char(64) NOT NULL COMMENT '文件内容哈希',
    `name` varchar(255) NOT NULL COMMENT '文件名称',
    `ext` varchar(30) DEFAULT NULL COMMENT '文件扩展名',
    `size` bigint NOT NULL COMMENT '文件大小(字节)',
    `path` varchar(255) NOT NULL COMMENT '文件存储路径',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_repo_identity` (`identity`),
    UNIQUE KEY `uk_repo_hash` (`hash`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `user_repository` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `identity` char(36) NOT NULL COMMENT '用户文件唯一标识',
    `parent_id` bigint unsigned DEFAULT NULL COMMENT '父级目录ID，NULL为根目录',
    `user_identity` char(36) NOT NULL COMMENT '对应用户的唯一标识',
    `repository_identity` char(36) DEFAULT NULL COMMENT '公共池中文件的唯一标识，文件夹为空',
    `ext` varchar(30) DEFAULT NULL COMMENT '文件扩展名',
    `name` varchar(255) NOT NULL COMMENT '用户定义的文件名',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_repo_identity` (`identity`),
    KEY `idx_user_parent` (`user_identity`, `parent_id`),
    UNIQUE KEY `uk_user_parent_name` (`user_identity`, `name`),
    CONSTRAINT `fk_user_repo_user` FOREIGN KEY (`user_identity`) REFERENCES `user_basic` (`identity`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `share_link` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `identity` char(36) NOT NULL COMMENT '分享唯一标识',
    `user_identity` char(36) NOT NULL COMMENT '对应用户的唯一标识',
    `repository_identity` char(36) NOT NULL COMMENT '用户池子中的唯一标识',
    `share_code` varchar(32) DEFAULT NULL COMMENT '分享码',
    `expires` int DEFAULT NULL COMMENT '失效时间，NULL为永不失效',
    `click_num` int unsigned NOT NULL DEFAULT 0 COMMENT '点击次数',
    `password_hash` varchar(255) DEFAULT NULL COMMENT '访问密码哈希',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_share_identity` (`identity`),
    UNIQUE KEY `uk_share_code` (`share_code`),
    KEY `idx_share_user` (`user_identity`),
    KEY `idx_share_repo` (`repository_identity`),
    CONSTRAINT `fk_share_user` FOREIGN KEY (`user_identity`) REFERENCES `user_basic` (`identity`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE `audit_log` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `identity` char(36) NOT NULL COMMENT '日志唯一标识',
    `actor_identity` char(36) DEFAULT NULL COMMENT '操作人标识',
    `actor_name` varchar(60) DEFAULT NULL COMMENT '操作人名称',
    `actor_role` tinyint NOT NULL DEFAULT 1 COMMENT '1-普通用户 2-管理员',
    `action` varchar(64) NOT NULL COMMENT '操作类型',
    `target_type` varchar(64) DEFAULT NULL COMMENT '目标类型',
    `target_identity` char(36) DEFAULT NULL COMMENT '目标标识',
    `detail` varchar(500) DEFAULT NULL COMMENT '操作描述',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_audit_identity` (`identity`),
    KEY `idx_audit_actor` (`actor_identity`),
    KEY `idx_audit_action` (`action`),
    KEY `idx_audit_created` (`created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;