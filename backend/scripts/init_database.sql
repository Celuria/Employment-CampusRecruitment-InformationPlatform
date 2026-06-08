-- =============================================================================
-- 就业中心校招信息平台 — MySQL 数据库初始化脚本
-- 数据库名: campus_recruitment
-- 应用账号: employ_user / 123456
-- 适用版本: MySQL 8.0+
-- 执行方式: mysql -u root -p < scripts/init_database.sql
-- =============================================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -----------------------------------------------------------------------------
-- 1. 创建数据库
-- -----------------------------------------------------------------------------
CREATE DATABASE IF NOT EXISTS `campus_recruitment`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `campus_recruitment`;

-- -----------------------------------------------------------------------------
-- 2. 创建应用数据库用户并授权
-- -----------------------------------------------------------------------------
CREATE USER IF NOT EXISTS 'employ_user'@'localhost' IDENTIFIED BY '123456';
CREATE USER IF NOT EXISTS 'employ_user'@'127.0.0.1' IDENTIFIED BY '123456';
CREATE USER IF NOT EXISTS 'employ_user'@'%' IDENTIFIED BY '123456';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, REFERENCES
  ON `campus_recruitment`.* TO 'employ_user'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, REFERENCES
  ON `campus_recruitment`.* TO 'employ_user'@'127.0.0.1';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, INDEX, REFERENCES
  ON `campus_recruitment`.* TO 'employ_user'@'%';

FLUSH PRIVILEGES;

-- -----------------------------------------------------------------------------
-- 3. 删除旧表（如需重建，取消注释）
-- -----------------------------------------------------------------------------
-- DROP TABLE IF EXISTS `reminder_logs`;
-- DROP TABLE IF EXISTS `calendar_events`;
-- DROP TABLE IF EXISTS `job_fairs`;
-- DROP TABLE IF EXISTS `career_talks`;
-- DROP TABLE IF EXISTS `user_preferences`;
-- DROP TABLE IF EXISTS `users`;

-- -----------------------------------------------------------------------------
-- 4. 用户表 users
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `users` (
  `id`                BIGINT UNSIGNED   NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username`          VARCHAR(64)       NOT NULL COMMENT '登录账号',
  `password_hash`     VARCHAR(255)      NOT NULL COMMENT '密码哈希(bcrypt)',
  `email`             VARCHAR(128)      NOT NULL COMMENT '邮箱',
  `name`              VARCHAR(64)       DEFAULT NULL COMMENT '真实姓名',
  `college`           VARCHAR(64)       DEFAULT NULL COMMENT '学院',
  `major`             VARCHAR(64)       DEFAULT NULL COMMENT '专业',
  `grade`             VARCHAR(16)       DEFAULT NULL COMMENT '年级',
  `target_positions`  JSON              DEFAULT NULL COMMENT '就业意向岗位',
  `phone`             VARCHAR(20)       DEFAULT NULL COMMENT '联系电话',
  `avatar`            VARCHAR(512)      DEFAULT NULL COMMENT '头像URL',
  `role`              VARCHAR(16)       NOT NULL DEFAULT 'student' COMMENT '角色: student/admin',
  `status`            VARCHAR(16)       NOT NULL DEFAULT 'active' COMMENT '状态: active/locked/disabled',
  `login_attempts`    INT               NOT NULL DEFAULT 0 COMMENT '连续登录失败次数',
  `locked_until`      DATETIME          DEFAULT NULL COMMENT '锁定截止时间',
  `last_login_at`     DATETIME          DEFAULT NULL COMMENT '最后登录时间',
  `created_at`        DATETIME          NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`        DATETIME          NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_username` (`username`),
  UNIQUE KEY `uk_users_email` (`email`),
  KEY `idx_users_role` (`role`),
  KEY `idx_users_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- -----------------------------------------------------------------------------
-- 5. 用户偏好表 user_preferences
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_preferences` (
  `id`                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '偏好ID',
  `user_id`             BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `target_positions`    JSON            DEFAULT NULL COMMENT '意向岗位',
  `preferred_cities`    JSON            DEFAULT NULL COMMENT '偏好城市',
  `preferred_companies` JSON            DEFAULT NULL COMMENT '偏好公司',
  `focus_companies`     JSON            DEFAULT NULL COMMENT '特别关注公司',
  `remind_before`       JSON            DEFAULT NULL COMMENT '提醒提前量: 1h/1d/3d',
  `created_at`          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_preferences_user_id` (`user_id`),
  CONSTRAINT `fk_user_preferences_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户偏好表';

-- -----------------------------------------------------------------------------
-- 6. 宣讲会表 career_talks
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `career_talks` (
  `id`               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '宣讲会ID',
  `title`            VARCHAR(200)    NOT NULL COMMENT '标题',
  `company`          VARCHAR(128)    NOT NULL COMMENT '举办公司',
  `industry_code`    VARCHAR(32)     DEFAULT NULL COMMENT '行业代码',
  `company_size`     VARCHAR(64)     DEFAULT NULL COMMENT '公司规模',
  `start_time`       DATETIME        NOT NULL COMMENT '开始时间',
  `end_time`         DATETIME        DEFAULT NULL COMMENT '结束时间',
  `location`         VARCHAR(256)    NOT NULL COMMENT '举办地点',
  `campus`           VARCHAR(32)     DEFAULT NULL COMMENT '校区: main/shahe/online',
  `format`           VARCHAR(16)     NOT NULL COMMENT '形式: online/offline/hybrid',
  `positions`        JSON            DEFAULT NULL COMMENT '面向岗位',
  `target_majors`    JSON            DEFAULT NULL COMMENT '面向专业',
  `registration_url` VARCHAR(512)    DEFAULT NULL COMMENT '报名链接',
  `source_url`       VARCHAR(512)    DEFAULT NULL COMMENT '信息来源URL',
  `logo_url`         VARCHAR(512)    DEFAULT NULL COMMENT 'Logo URL',
  `description`      TEXT            DEFAULT NULL COMMENT '详细描述',
  `publish_status`   VARCHAR(16)     NOT NULL DEFAULT 'draft' COMMENT '发布状态: draft/published/archived',
  `source_type`      VARCHAR(16)     NOT NULL DEFAULT 'manual' COMMENT '来源: manual/sync',
  `synced_at`        DATETIME        DEFAULT NULL COMMENT '最近同步时间',
  `created_by`       BIGINT UNSIGNED DEFAULT NULL COMMENT '创建人',
  `updated_by`       BIGINT UNSIGNED DEFAULT NULL COMMENT '更新人',
  `created_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_career_talks_start_time` (`start_time`),
  KEY `idx_career_talks_publish_status` (`publish_status`),
  KEY `idx_career_talks_company` (`company`),
  KEY `idx_career_talks_campus` (`campus`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='宣讲会表';

-- -----------------------------------------------------------------------------
-- 7. 双选会表 job_fairs
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `job_fairs` (
  `id`               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '双选会ID',
  `title`            VARCHAR(200)    NOT NULL COMMENT '名称',
  `start_date`       DATE            NOT NULL COMMENT '开始日期',
  `end_date`         DATE            DEFAULT NULL COMMENT '结束日期',
  `start_time`       DATETIME        DEFAULT NULL COMMENT '开始时刻',
  `location`         VARCHAR(256)    NOT NULL COMMENT '举办地点',
  `campus`           VARCHAR(32)     DEFAULT NULL COMMENT '校区',
  `company_count`    INT             DEFAULT NULL COMMENT '参与企业数量',
  `target_audience`  VARCHAR(256)    DEFAULT NULL COMMENT '面向对象',
  `target_majors`    JSON            DEFAULT NULL COMMENT '面向专业',
  `deadline`         DATETIME        DEFAULT NULL COMMENT '报名截止时间',
  `detail_url`       VARCHAR(512)    DEFAULT NULL COMMENT '详情链接',
  `source_url`       VARCHAR(512)    DEFAULT NULL COMMENT '信息来源URL',
  `description`      TEXT            DEFAULT NULL COMMENT '详细描述',
  `publish_status`   VARCHAR(16)     NOT NULL DEFAULT 'draft' COMMENT '发布状态',
  `source_type`      VARCHAR(16)     NOT NULL DEFAULT 'manual' COMMENT '来源',
  `synced_at`        DATETIME        DEFAULT NULL COMMENT '最近同步时间',
  `created_by`       BIGINT UNSIGNED DEFAULT NULL COMMENT '创建人',
  `updated_by`       BIGINT UNSIGNED DEFAULT NULL COMMENT '更新人',
  `created_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_job_fairs_start_date` (`start_date`),
  KEY `idx_job_fairs_publish_status` (`publish_status`),
  KEY `idx_job_fairs_deadline` (`deadline`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='双选会表';

-- -----------------------------------------------------------------------------
-- 8. 日历事件表 calendar_events
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `calendar_events` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日历事件ID',
  `user_id`       BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `event_type`    VARCHAR(16)     NOT NULL COMMENT '类型: career_talk/job_fair',
  `ref_id`        BIGINT UNSIGNED NOT NULL COMMENT '关联活动ID',
  `title`         VARCHAR(200)    NOT NULL COMMENT '活动标题',
  `start_time`    DATETIME        NOT NULL COMMENT '开始时间',
  `end_time`      DATETIME        DEFAULT NULL COMMENT '结束时间',
  `location`      VARCHAR(256)    DEFAULT NULL COMMENT '地点',
  `custom_note`   VARCHAR(500)    DEFAULT NULL COMMENT '个人备注',
  `remind_before` JSON            DEFAULT NULL COMMENT '提醒提前量',
  `status`        VARCHAR(16)     NOT NULL DEFAULT 'active' COMMENT '状态',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_calendar_user_event` (`user_id`, `event_type`, `ref_id`),
  KEY `idx_calendar_user_start` (`user_id`, `start_time`),
  CONSTRAINT `fk_calendar_events_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='日历事件表';

-- -----------------------------------------------------------------------------
-- 9. 提醒记录表 reminder_logs
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `reminder_logs` (
  `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '提醒记录ID',
  `calendar_event_id` BIGINT UNSIGNED NOT NULL COMMENT '日历事件ID',
  `user_id`           BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `event_title`       VARCHAR(200)    DEFAULT NULL COMMENT '活动标题',
  `event_type`        VARCHAR(16)     DEFAULT NULL COMMENT '活动类型',
  `remind_before`     VARCHAR(8)      DEFAULT NULL COMMENT '提前量: 1h/1d/3d',
  `scheduled_time`    DATETIME        NOT NULL COMMENT '计划发送时间',
  `sent_time`         DATETIME        DEFAULT NULL COMMENT '实际发送时间',
  `status`            VARCHAR(16)     NOT NULL DEFAULT 'pending' COMMENT '状态: pending/sent/failed/cancelled',
  `retry_count`       INT             NOT NULL DEFAULT 0 COMMENT '重试次数',
  `fail_reason`       VARCHAR(512)    DEFAULT NULL COMMENT '失败原因',
  `created_at`        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_reminder_status_scheduled` (`status`, `scheduled_time`),
  KEY `idx_reminder_user_id` (`user_id`),
  KEY `idx_reminder_calendar_event_id` (`calendar_event_id`),
  CONSTRAINT `fk_reminder_logs_calendar_event_id` FOREIGN KEY (`calendar_event_id`) REFERENCES `calendar_events` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_reminder_logs_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='邮件提醒记录表';

-- -----------------------------------------------------------------------------
-- 10. 审计日志表 audit_logs
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `audit_logs` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `operator_id`   BIGINT UNSIGNED NOT NULL COMMENT '操作人ID',
  `action`        VARCHAR(16)     NOT NULL COMMENT '动作: CREATE/UPDATE/DELETE/SYNC',
  `resource_type` VARCHAR(32)     NOT NULL COMMENT '资源类型',
  `resource_id`   BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '资源ID',
  `detail`        TEXT            DEFAULT NULL COMMENT '操作详情JSON',
  `ip`            VARCHAR(64)     DEFAULT NULL COMMENT '操作IP',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_audit_operator` (`operator_id`),
  KEY `idx_audit_action` (`action`),
  KEY `idx_audit_resource` (`resource_type`),
  KEY `idx_audit_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理端审计日志';

-- -----------------------------------------------------------------------------
-- 11. 同步记录表 sync_logs
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `sync_logs` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `task_id`       VARCHAR(64)     NOT NULL COMMENT '任务ID',
  `source_type`   VARCHAR(32)     NOT NULL COMMENT '同步类型',
  `status`        VARCHAR(16)     NOT NULL COMMENT '状态: pending/running/success/failed',
  `added_count`   INT             NOT NULL DEFAULT 0 COMMENT '新增条数',
  `updated_count` INT             NOT NULL DEFAULT 0 COMMENT '更新条数',
  `failed_count`  INT             NOT NULL DEFAULT 0 COMMENT '失败条数',
  `started_at`    DATETIME        NOT NULL COMMENT '开始时间',
  `finished_at`   DATETIME        DEFAULT NULL COMMENT '结束时间',
  `operator_id`   BIGINT UNSIGNED NOT NULL COMMENT '操作人ID',
  `error_message` VARCHAR(512)    DEFAULT NULL COMMENT '错误信息',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_sync_task` (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='信息同步记录';

SET FOREIGN_KEY_CHECKS = 1;

-- -----------------------------------------------------------------------------
-- 12. 初始化演示数据（可选）
-- -----------------------------------------------------------------------------

-- 初始管理员账号（用户名 admin，密码 Admin@123456）
INSERT INTO `users` (`username`, `password_hash`, `email`, `name`, `role`, `status`)
VALUES (
  'admin',
  '$2a$10$3HQVFi1u4oXsG4AnK55gqeuSil2IVlGpDbZb/.wlLAk8.MR1yBFou',
  'admin@employment-center.edu.cn',
  '系统管理员',
  'admin',
  'active'
) ON DUPLICATE KEY UPDATE `username` = `username`;

-- 示例宣讲会（已发布）
INSERT INTO `career_talks` (
  `title`, `company`, `industry_code`, `company_size`,
  `start_time`, `location`, `campus`, `format`,
  `positions`, `target_majors`, `publish_status`, `source_type`
) VALUES
(
  '字节跳动2025校园招聘技术专场宣讲会',
  '字节跳动', 'internet', '10000人以上',
  '2025-01-15 14:00:00', '本部校区 · 学术报告厅A301', 'main', 'hybrid',
  JSON_ARRAY('研发工程师', '产品经理', '算法'),
  JSON_ARRAY('计算机', '软件工程'),
  'published', 'manual'
),
(
  '阿里巴巴2025校招「星耀计划」全球宣讲会',
  '阿里巴巴集团', 'internet', '10000人以上',
  '2025-01-16 19:00:00', '沙河校区 · 教学楼B205', 'shahe', 'offline',
  JSON_ARRAY('Java开发', '数据分析师', '运营'),
  JSON_ARRAY('计算机', '信息管理'),
  'published', 'manual'
);

-- 示例双选会（已发布）
INSERT INTO `job_fairs` (
  `title`, `start_date`, `end_date`, `location`, `campus`,
  `company_count`, `target_audience`, `deadline`, `publish_status`, `source_type`
) VALUES
(
  '2025届春季大型双选会',
  '2025-03-15', '2025-03-16', '本部校区 · 体育馆', 'main',
  120, '2025届全体毕业生',
  '2025-03-10 18:00:00', 'published', 'manual'
);

-- =============================================================================
-- 初始化完成
-- 连接示例:
--   mysql -u employ_user -p123456 -h 127.0.0.1 campus_recruitment
--
-- 后端 config/config.yaml 请配置:
--   database:
--     user: employ_user
--     password: "123456"
--     dbname: campus_recruitment
-- =============================================================================
