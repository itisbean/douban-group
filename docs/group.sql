/*
Navicat MySQL Data Transfer

Source Database       : spiders

Target Server Type    : MYSQL
Target Server Version : 50639
File Encoding         : 65001

Date: 2019-07-08 21:08:00
*/

-- ----------------------------
-- Table structure for sp_douban_group_dbhyz
-- ----------------------------
CREATE TABLE `sp_douban_group_dbhyz` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `topic_id` int(10) DEFAULT '0' COMMENT '标题ID',
  `topic` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '标题',
  `author_id` int(10) DEFAULT '0' COMMENT '发帖人ID',
  `author` varchar(20) DEFAULT '' COMMENT '发帖人',
  `create_time` datetime COMMENT '创建时间',
  `new_reply_time` datetime COMMENT '最后回复时间',
  `reply` int(10) unsigned DEFAULT '0' COMMENT '回复数',
  `liked` int(10) unsigned DEFAULT '0' COMMENT '点赞数量',
  `collect` int(10) unsigned DEFAULT '0' COMMENT '收藏数量',
  `sharing` int(10) unsigned DEFAULT '0' COMMENT '转发',
	`url` varchar(120) DEFAULT '' COMMENT '链接',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '内容',
  `version` int(10) unsigned DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `IX_topic_id`(`topic_id`),
  INDEX `version`(`version`),
  INDEX `IX_new_time`(`new_reply_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT='豆瓣火研组';


ALTER TABLE `spiders`.`sp_douban_group_dbhyz` 
ADD COLUMN `is_del` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否删除' AFTER `version`;