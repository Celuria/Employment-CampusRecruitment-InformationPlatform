-- 校区枚举迁移：main/shahe -> nanhu/mafangshan/yujiato，并新增 venue 楼栋字段
-- 执行: mysql -u employ_user -p campus_recruitment < migrations/002_campus_venue.sql

USE `campus_recruitment`;

-- 若已开启 auto_migrate 或列已存在，可跳过 ALTER 语句
ALTER TABLE `career_talks`
  ADD COLUMN `venue` VARCHAR(128) DEFAULT NULL COMMENT '具体楼栋地点' AFTER `campus`;

ALTER TABLE `job_fairs`
  ADD COLUMN `venue` VARCHAR(128) DEFAULT NULL COMMENT '具体楼栋地点' AFTER `campus`;

UPDATE `career_talks`
SET `venue` = TRIM(SUBSTRING_INDEX(`location`, '·', -1)),
    `campus` = 'nanhu',
    `location` = REPLACE(`location`, '本部校区', '南湖校区')
WHERE `campus` = 'main';

UPDATE `career_talks`
SET `venue` = TRIM(SUBSTRING_INDEX(`location`, '·', -1)),
    `campus` = 'mafangshan',
    `location` = REPLACE(`location`, '沙河校区', '马房山校区')
WHERE `campus` = 'shahe';

UPDATE `job_fairs`
SET `venue` = TRIM(SUBSTRING_INDEX(`location`, '·', -1)),
    `campus` = 'nanhu',
    `location` = REPLACE(`location`, '本部校区', '南湖校区')
WHERE `campus` = 'main';

UPDATE `job_fairs`
SET `venue` = TRIM(SUBSTRING_INDEX(`location`, '·', -1)),
    `campus` = 'mafangshan',
    `location` = REPLACE(`location`, '沙河校区', '马房山校区')
WHERE `campus` = 'shahe';
