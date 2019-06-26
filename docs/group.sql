CREATE TABLE `sp_douban_group_dbhyz` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `topic_id` int(10) DEFAULT '0' COMMENT '标题ID',
  `topic` varchar(200) DEFAULT '' COMMENT '标题',
  `author_id` int(10) DEFAULT '0' COMMENT '发帖人ID',
  `author` varchar(20) DEFAULT '' COMMENT '发帖人',
  `create_time` datetime COMMENT '创建时间',
  `new_reply_time` datetime COMMENT '最后回复时间',
  `reply` int(10) unsigned DEFAULT '0' COMMENT '回复数',
  `liked` int(10) unsigned DEFAULT '0' COMMENT '点赞数量',
  `collect` int(10) unsigned DEFAULT '0' COMMENT '收藏数量',
  `sharing` int(10) unsigned DEFAULT '0' COMMENT '转发',
	`url` varchar(120) DEFAULT '' COMMENT '链接',
  `content` text COMMENT '内容',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='豆瓣火研组';
