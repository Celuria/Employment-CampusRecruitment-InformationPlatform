-- =============================================================================
-- 就业中心校招信息平台 — 演示假数据脚本
-- 数据库名: campus_recruitment
-- 适用版本: MySQL 8.0+
-- 前置条件: 已执行 init_database.sql 建库建表
-- 执行方式（Linux / macOS）:
--   mysql -u employ_user -p123456 -h 127.0.0.1 --default-character-set=utf8mb4 campus_recruitment < scripts/seed_demo_data.sql
-- 执行方式（Windows，避免中文乱码）:
--   mysql -u employ_user -p123456 -h 127.0.0.1 --default-character-set=utf8mb4 campus_recruitment -e "source F:/Project/.../backend/scripts/seed_demo_data.sql"
--
-- 说明:
--   1. 本脚本可重复执行：会先清理 [演示] 前缀的活动及 demo_ 前缀的测试用户
--   2. 演示学生统一密码: Student@123456
--   3. 管理员账号 admin 不会被删除（见 init_database.sql）
-- =============================================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

USE `campus_recruitment`;

-- -----------------------------------------------------------------------------
-- 0. 清理旧演示数据（按标记删除，不影响正式数据）
-- -----------------------------------------------------------------------------
DELETE FROM `reminder_logs`
WHERE `user_id` IN (
  SELECT `id` FROM (
    SELECT `id` FROM `users` WHERE `username` LIKE 'demo\_%' ESCAPE '\\'
  ) AS `demo_users`
);

DELETE FROM `calendar_events`
WHERE `user_id` IN (
  SELECT `id` FROM (
    SELECT `id` FROM `users` WHERE `username` LIKE 'demo\_%' ESCAPE '\\'
  ) AS `demo_users`
);

DELETE FROM `user_preferences`
WHERE `user_id` IN (
  SELECT `id` FROM (
    SELECT `id` FROM `users` WHERE `username` LIKE 'demo\_%' ESCAPE '\\'
  ) AS `demo_users`
);

DELETE FROM `users` WHERE `username` LIKE 'demo\_%' ESCAPE '\\';

DELETE FROM `career_talks` WHERE `title` LIKE '[演示]%';
DELETE FROM `job_fairs` WHERE `title` LIKE '[演示]%';

DELETE FROM `audit_logs` WHERE `detail` LIKE '%"seed":"demo"%';
DELETE FROM `sync_logs` WHERE `task_id` LIKE 'demo-seed-%';

-- -----------------------------------------------------------------------------
-- 1. 演示学生账号（密码均为 Student@123456）
--    bcrypt: $2a$10$H6uXAx5QEeKJWx/3QDearuoOcK5XsiCjRq/Gv1OHJP7/j7GwUmfF2
-- -----------------------------------------------------------------------------
INSERT INTO `users` (
  `username`, `password_hash`, `email`, `name`, `college`, `major`, `grade`,
  `target_positions`, `phone`, `role`, `status`, `last_login_at`
) VALUES
(
  'demo_cs01', '$2a$10$H6uXAx5QEeKJWx/3QDearuoOcK5XsiCjRq/Gv1OHJP7/j7GwUmfF2',
  'demo.cs01@student.edu.cn', '张明', '计算机学院', '软件工程', '2026',
  JSON_ARRAY('前端开发', '后端开发'), '13800001001', 'student', 'active', '2026-06-07 20:30:00'
),
(
  'demo_cs02', '$2a$10$H6uXAx5QEeKJWx/3QDearuoOcK5XsiCjRq/Gv1OHJP7/j7GwUmfF2',
  'demo.cs02@student.edu.cn', '李雪', '计算机学院', '计算机科学与技术', '2025',
  JSON_ARRAY('算法工程师', '数据分析'), '13800001002', 'student', 'active', '2026-06-06 18:15:00'
),
(
  'demo_em01', '$2a$10$H6uXAx5QEeKJWx/3QDearuoOcK5XsiCjRq/Gv1OHJP7/j7GwUmfF2',
  'demo.em01@student.edu.cn', '王芳', '经济管理学院', '会计学', '2026',
  JSON_ARRAY('会计', '财务', '审计'), '13800002001', 'student', 'active', '2026-06-05 09:00:00'
),
(
  'demo_me01', '$2a$10$H6uXAx5QEeKJWx/3QDearuoOcK5XsiCjRq/Gv1OHJP7/j7GwUmfF2',
  'demo.me01@student.edu.cn', '赵强', '机械工程学院', '机械设计制造及其自动化', '2025',
  JSON_ARRAY('机械设计', '工艺工程师'), '13800003001', 'student', 'active', '2026-06-04 14:20:00'
),
(
  'demo_ee01', '$2a$10$H6uXAx5QEeKJWx/3QDearuoOcK5XsiCjRq/Gv1OHJP7/j7GwUmfF2',
  'demo.ee01@student.edu.cn', '陈琳', '电子信息学院', '通信工程', '2026',
  JSON_ARRAY('嵌入式开发', '硬件工程师'), '13800004001', 'student', 'active', '2026-06-03 11:45:00'
),
(
  'demo_ma01', '$2a$10$H6uXAx5QEeKJWx/3QDearuoOcK5XsiCjRq/Gv1OHJP7/j7GwUmfF2',
  'demo.ma01@student.edu.cn', '刘洋', '数学学院', '统计学', '2025',
  JSON_ARRAY('数据分析', '量化分析师'), '13800005001', 'student', 'active', '2026-06-02 16:00:00'
);

-- -----------------------------------------------------------------------------
-- 2. 用户偏好
-- -----------------------------------------------------------------------------
INSERT INTO `user_preferences` (
  `user_id`, `target_positions`, `preferred_cities`, `preferred_companies`,
  `focus_companies`, `remind_before`
)
SELECT `id`,
  `target_positions`,
  JSON_ARRAY('武汉', '北京', '上海'),
  JSON_ARRAY('字节跳动', '华为', '腾讯'),
  JSON_ARRAY('字节跳动'),
  JSON_ARRAY('1d', '3d')
FROM `users` WHERE `username` = 'demo_cs01';

INSERT INTO `user_preferences` (
  `user_id`, `target_positions`, `preferred_cities`, `preferred_companies`,
  `focus_companies`, `remind_before`
)
SELECT `id`,
  `target_positions`,
  JSON_ARRAY('深圳', '杭州', '成都'),
  JSON_ARRAY('阿里巴巴', '美团', '京东'),
  JSON_ARRAY('阿里巴巴', '美团'),
  JSON_ARRAY('1h', '1d')
FROM `users` WHERE `username` = 'demo_cs02';

INSERT INTO `user_preferences` (
  `user_id`, `target_positions`, `preferred_cities`, `preferred_companies`,
  `focus_companies`, `remind_before`
)
SELECT `id`,
  `target_positions`,
  JSON_ARRAY('武汉', '广州', '上海'),
  JSON_ARRAY('招商银行', '普华永道', '德勤'),
  JSON_ARRAY('招商银行'),
  JSON_ARRAY('1d')
FROM `users` WHERE `username` = 'demo_em01';

INSERT INTO `user_preferences` (
  `user_id`, `target_positions`, `preferred_cities`, `preferred_companies`,
  `focus_companies`, `remind_before`
)
SELECT `id`,
  `target_positions`,
  JSON_ARRAY('武汉', '苏州', '重庆'),
  JSON_ARRAY('比亚迪', '宁德时代', '三一重工'),
  JSON_ARRAY('比亚迪'),
  JSON_ARRAY('3d')
FROM `users` WHERE `username` = 'demo_me01';

INSERT INTO `user_preferences` (
  `user_id`, `target_positions`, `preferred_cities`, `preferred_companies`,
  `focus_companies`, `remind_before`
)
SELECT `id`,
  `target_positions`,
  JSON_ARRAY('武汉', '南京', '西安'),
  JSON_ARRAY('华为', '中兴', '小米'),
  JSON_ARRAY('华为'),
  JSON_ARRAY('1h', '3d')
FROM `users` WHERE `username` = 'demo_ee01';

INSERT INTO `user_preferences` (
  `user_id`, `target_positions`, `preferred_cities`, `preferred_companies`,
  `focus_companies`, `remind_before`
)
SELECT `id`,
  `target_positions`,
  JSON_ARRAY('上海', '北京', '香港'),
  JSON_ARRAY('中金公司', '华泰证券', '蚂蚁集团'),
  JSON_ARRAY('中金公司'),
  JSON_ARRAY('1d', '3d')
FROM `users` WHERE `username` = 'demo_ma01';

-- -----------------------------------------------------------------------------
-- 3. 宣讲会（标题以 [演示] 开头，时间设在 2026 年 6–9 月）
-- -----------------------------------------------------------------------------
INSERT INTO `career_talks` (
  `title`, `company`, `industry_code`, `company_size`,
  `start_time`, `end_time`, `location`, `campus`, `venue`, `format`,
  `positions`, `target_majors`, `registration_url`, `description`,
  `publish_status`, `source_type`, `created_by`
) VALUES
(
  '[演示] 腾讯 2026 校园招聘技术专场',
  '腾讯', 'internet', '10000人以上',
  '2026-06-12 14:00:00', '2026-06-12 16:00:00',
  '南湖校区 · 图书馆报告厅A301', 'nanhu', '图书馆报告厅A301', 'offline',
  JSON_ARRAY('后端开发', '前端开发', '游戏开发', '产品经理'),
  JSON_ARRAY('软件工程', '计算机科学与技术', '网络工程'),
  'https://join.qq.com', '面向 2026 届计算机相关专业，含笔试经验分享与现场答疑。',
  'published', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
),
(
  '[演示] 华为 2026 研发岗宣讲会',
  '华为技术有限公司', 'internet', '10000人以上',
  '2026-06-15 19:00:00', '2026-06-15 21:00:00',
  '马房山校区 · 东院教学楼B205', 'mafangshan', '东院教学楼B205', 'hybrid',
  JSON_ARRAY('嵌入式开发', '硬件工程师', '通信工程师', '算法工程师'),
  JSON_ARRAY('电子信息工程', '通信工程', '自动化', '计算机科学与技术'),
  'https://career.huawei.com', '含线上直播链接，欢迎电子信息与计算机学院同学参加。',
  'published', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
),
(
  '[演示] 招商银行 管培生 & 金融科技专场',
  '招商银行', 'finance', '10000人以上',
  '2026-06-18 14:30:00', '2026-06-18 16:30:00',
  '南湖校区 · 经管学院报告厅', 'nanhu', '经管学院报告厅', 'offline',
  JSON_ARRAY('管培生', '金融科技', '数据分析', '风险管理'),
  JSON_ARRAY('金融学', '会计学', '工商管理', '统计学'),
  'https://career.cmbchina.com', '经管学院专场，含简历一对一辅导环节。',
  'published', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
),
(
  '[演示] 比亚迪 智能制造工程师宣讲',
  '比亚迪', 'manufacturing', '10000人以上',
  '2026-06-20 10:00:00', '2026-06-20 12:00:00',
  '余家头校区 · 机械楼201', 'yujiato', '机械楼201', 'offline',
  JSON_ARRAY('机械设计', '工艺工程师', '质量工程师', '智能制造工程师'),
  JSON_ARRAY('机械设计制造及其自动化', '车辆工程', '工业设计'),
  'https://hr.byd.com', '新能源汽车板块重点招聘，提供实习转正通道。',
  'published', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
),
(
  '[演示] 普华永道 审计 & 咨询校招宣讲',
  '普华永道', 'consulting', '10000人以上',
  '2026-06-22 19:00:00', '2026-06-22 20:30:00',
  '南湖校区 · 外国语学院多功能厅', 'nanhu', '外国语学院多功能厅', 'offline',
  JSON_ARRAY('审计', '税务', '管理咨询', '财务顾问'),
  JSON_ARRAY('会计学', '工商管理', '英语', '金融学'),
  'https://www.pwc.com/cn/careers', '四大会计师事务所专场，中英文简历均可投递。',
  'published', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
),
(
  '[演示] 美团 产品 & 运营校招线上宣讲',
  '美团', 'internet', '10000人以上',
  '2026-06-25 19:30:00', '2026-06-25 21:00:00',
  '线上 · 腾讯会议', 'online', '腾讯会议', 'online',
  JSON_ARRAY('产品经理', '运营', '数据分析', '商业分析'),
  JSON_ARRAY('计算机科学与技术', '信息管理', '市场营销', '统计学'),
  'https://zhaopin.meituan.com', '纯线上宣讲，报名后邮件发送会议号。',
  'published', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
),
(
  '[演示] 中金公司 量化研究校招专场',
  '中金公司', 'finance', '1000-9999人',
  '2026-07-03 15:00:00', '2026-07-03 17:00:00',
  '南湖校区 · 数学学院学术报告厅', 'nanhu', '数学学院学术报告厅', 'offline',
  JSON_ARRAY('量化分析师', '金融建模', '风险管理', '行业研究'),
  JSON_ARRAY('数学与应用数学', '统计学', '金融学', '信息与计算科学'),
  'https://www.cicc.com/careers', '数学与金融交叉学科优先，含笔试题型讲解。',
  'published', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
),
(
  '[演示] 小米 硬件 & 软件联合宣讲会',
  '小米', 'internet', '10000人以上',
  '2026-07-08 14:00:00', '2026-07-08 16:00:00',
  '马房山校区 · 电子信息楼报告厅', 'mafangshan', '电子信息楼报告厅', 'offline',
  JSON_ARRAY('硬件工程师', '嵌入式开发', 'Android开发', '测试工程师'),
  JSON_ARRAY('电子信息工程', '微电子科学与工程', '软件工程', '自动化'),
  'https://hr.xiaomi.com', 'IoT 与手机业务线联合招聘。',
  'published', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
),
(
  '[演示] 三一重工 结构工程师宣讲（草稿）',
  '三一重工', 'manufacturing', '10000人以上',
  '2026-07-15 10:00:00', '2026-07-15 11:30:00',
  '余家头校区 · 工程训练中心', 'yujiato', '工程训练中心', 'offline',
  JSON_ARRAY('结构工程师', '机械设计', '液压工程师'),
  JSON_ARRAY('机械设计制造及其自动化', '车辆工程'),
  NULL, '草稿状态，仅供管理后台测试。',
  'draft', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
),
(
  '[演示] 德勤 秋季补招宣讲会',
  '德勤', 'consulting', '10000人以上',
  '2026-09-10 19:00:00', '2026-09-10 20:30:00',
  '南湖校区 · 就业指导中心', 'nanhu', '就业指导中心', 'offline',
  JSON_ARRAY('审计', '咨询', '税务', '风险咨询'),
  JSON_ARRAY('会计学', '工商管理', '金融学'),
  'https://www2.deloitte.com/cn/careers', '秋季补招专场，面向 2025 届未就业毕业生。',
  'published', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
);

-- -----------------------------------------------------------------------------
-- 4. 双选会
-- -----------------------------------------------------------------------------
INSERT INTO `job_fairs` (
  `title`, `start_date`, `end_date`, `start_time`, `location`, `campus`, `venue`,
  `company_count`, `target_audience`, `target_majors`, `deadline`, `description`,
  `publish_status`, `source_type`, `created_by`
) VALUES
(
  '[演示] 2026 届夏季综合类双选会',
  '2026-06-28', '2026-06-29', '2026-06-28 09:00:00',
  '余家头校区 · 体育馆', 'yujiato', '体育馆',
  85, '2026 届全体毕业生',
  JSON_ARRAY('不限专业'),
  '2026-06-25 18:00:00',
  '涵盖互联网、制造、金融、咨询等行业，提供现场面试机会。',
  'published', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
),
(
  '[演示] 计算机 & 电子信息专场双选会',
  '2026-07-05', '2026-07-05', '2026-07-05 09:00:00',
  '马房山校区 · 大学生活动中心', 'mafangshan', '大学生活动中心',
  42, '计算机、电子信息、自动化等相关专业',
  JSON_ARRAY('软件工程', '计算机科学与技术', '电子信息工程', '通信工程', '自动化'),
  '2026-07-02 12:00:00',
  'IT 与硬科技专场，含华为、小米、中兴等知名企业。',
  'published', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
),
(
  '[演示] 经管类专场双选会',
  '2026-07-12', '2026-07-12', '2026-07-12 13:30:00',
  '南湖校区 · 经管学院广场', 'nanhu', '经管学院广场',
  36, '经济、管理、财会类专业',
  JSON_ARRAY('会计学', '金融学', '工商管理', '市场营销'),
  '2026-07-09 18:00:00',
  '银行、证券、会计师事务所及快消企业集中参展。',
  'published', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
),
(
  '[演示] 2027 届秋季大型双选会（预告）',
  '2026-09-20', '2026-09-21', '2026-09-20 08:30:00',
  '南湖校区 · 体育场', 'nanhu', '体育场',
  150, '2027 届全体毕业生',
  JSON_ARRAY('不限专业'),
  '2026-09-15 18:00:00',
  '秋季最大规模双选会，目前开放企业报名中。',
  'published', 'manual',
  (SELECT `id` FROM `users` WHERE `username` = 'admin' LIMIT 1)
);

-- -----------------------------------------------------------------------------
-- 5. 日历事件（部分演示学生已加入活动）
-- -----------------------------------------------------------------------------
INSERT INTO `calendar_events` (
  `user_id`, `event_type`, `ref_id`, `title`, `start_time`, `end_time`,
  `location`, `custom_note`, `remind_before`, `status`
)
SELECT u.`id`, 'career_talk', ct.`id`, ct.`title`, ct.`start_time`, ct.`end_time`,
  ct.`location`, '重点关注后端岗位', JSON_ARRAY('1d', '3d'), 'active'
FROM `users` u
JOIN `career_talks` ct ON ct.`title` = '[演示] 腾讯 2026 校园招聘技术专场'
WHERE u.`username` = 'demo_cs01';

INSERT INTO `calendar_events` (
  `user_id`, `event_type`, `ref_id`, `title`, `start_time`, `end_time`,
  `location`, `custom_note`, `remind_before`, `status`
)
SELECT u.`id`, 'career_talk', ct.`id`, ct.`title`, ct.`start_time`, ct.`end_time`,
  ct.`location`, '算法岗笔试准备', JSON_ARRAY('1h', '1d'), 'active'
FROM `users` u
JOIN `career_talks` ct ON ct.`title` = '[演示] 华为 2026 研发岗宣讲会'
WHERE u.`username` = 'demo_cs02';

INSERT INTO `calendar_events` (
  `user_id`, `event_type`, `ref_id`, `title`, `start_time`, `end_time`,
  `location`, `custom_note`, `remind_before`, `status`
)
SELECT u.`id`, 'career_talk', ct.`id`, ct.`title`, ct.`start_time`, ct.`end_time`,
  ct.`location`, NULL, JSON_ARRAY('1d'), 'active'
FROM `users` u
JOIN `career_talks` ct ON ct.`title` = '[演示] 招商银行 管培生 & 金融科技专场'
WHERE u.`username` = 'demo_em01';

INSERT INTO `calendar_events` (
  `user_id`, `event_type`, `ref_id`, `title`, `start_time`, `end_time`,
  `location`, `custom_note`, `remind_before`, `status`
)
SELECT u.`id`, 'job_fair', jf.`id`, jf.`title`, jf.`start_time`, NULL,
  jf.`location`, '带好简历原件', JSON_ARRAY('1d', '3d'), 'active'
FROM `users` u
JOIN `job_fairs` jf ON jf.`title` = '[演示] 2026 届夏季综合类双选会'
WHERE u.`username` IN ('demo_cs01', 'demo_me01');

INSERT INTO `calendar_events` (
  `user_id`, `event_type`, `ref_id`, `title`, `start_time`, `end_time`,
  `location`, `custom_note`, `remind_before`, `status`
)
SELECT u.`id`, 'job_fair', jf.`id`, jf.`title`, jf.`start_time`, NULL,
  jf.`location`, NULL, JSON_ARRAY('3d'), 'active'
FROM `users` u
JOIN `job_fairs` jf ON jf.`title` = '[演示] 经管类专场双选会'
WHERE u.`username` = 'demo_em01';

-- -----------------------------------------------------------------------------
-- 6. 提醒记录（模拟已发送与待发送）
-- -----------------------------------------------------------------------------
INSERT INTO `reminder_logs` (
  `calendar_event_id`, `user_id`, `event_title`, `event_type`,
  `remind_before`, `scheduled_time`, `sent_time`, `status`
)
SELECT ce.`id`, ce.`user_id`, ce.`title`, ce.`event_type`,
  '1d', DATE_SUB(ce.`start_time`, INTERVAL 1 DAY),
  DATE_SUB(ce.`start_time`, INTERVAL 1 DAY), 'sent'
FROM `calendar_events` ce
JOIN `users` u ON u.`id` = ce.`user_id`
WHERE u.`username` = 'demo_cs01'
  AND ce.`title` = '[演示] 腾讯 2026 校园招聘技术专场'
LIMIT 1;

INSERT INTO `reminder_logs` (
  `calendar_event_id`, `user_id`, `event_title`, `event_type`,
  `remind_before`, `scheduled_time`, `status`
)
SELECT ce.`id`, ce.`user_id`, ce.`title`, ce.`event_type`,
  '3d', DATE_SUB(ce.`start_time`, INTERVAL 3 DAY), 'pending'
FROM `calendar_events` ce
JOIN `users` u ON u.`id` = ce.`user_id`
WHERE u.`username` = 'demo_cs01'
  AND ce.`title` = '[演示] 腾讯 2026 校园招聘技术专场'
LIMIT 1;

-- -----------------------------------------------------------------------------
-- 7. 管理端审计 & 同步记录
-- -----------------------------------------------------------------------------
INSERT INTO `audit_logs` (`operator_id`, `action`, `resource_type`, `resource_id`, `detail`, `ip`)
SELECT `id`, 'CREATE', 'career_talk', 0,
  '{"seed":"demo","message":"批量导入演示宣讲会数据"}', '127.0.0.1'
FROM `users` WHERE `username` = 'admin' LIMIT 1;

INSERT INTO `sync_logs` (
  `task_id`, `source_type`, `status`, `added_count`, `updated_count`, `failed_count`,
  `started_at`, `finished_at`, `operator_id`, `error_message`
)
SELECT
  'demo-seed-20260608', 'employment_center', 'success', 10, 4, 0,
  '2026-06-08 08:00:00', '2026-06-08 08:02:15',
  `id`, NULL
FROM `users` WHERE `username` = 'admin' LIMIT 1;

SET FOREIGN_KEY_CHECKS = 1;

-- =============================================================================
-- 演示数据导入完成
--
-- 测试账号:
--   管理员  admin / Admin@123456
--   学生    demo_cs01 ~ demo_ma01 / Student@123456
--
-- 数据概览:
--   学生用户 6 名 | 宣讲会 10 场 | 双选会 4 场 | 日历事件若干 | 提醒记录 2 条
-- =============================================================================
