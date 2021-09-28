CREATE TABLE `version` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `platform` tinyint(4) DEFAULT '0' COMMENT '0:h5,1:iOS,2:android,3:wxma',
    `ver` varchar(20) DEFAULT '' COMMENT '版本号，如 1.9.1',
    `desc` varchar(200) DEFAULT '' COMMENT '版本更新内容',
    `url` varchar(200) DEFAULT '' COMMENT 'cdn更新地址',
    `force` tinyint(4) DEFAULT '0' COMMENT '这个版本是否强更，1为需要强更',
    `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_deleted` int(1) DEFAULT NULL COMMENT '0 未删除，1 已删除',
    PRIMARY KEY (`id`),
    KEY `i_plt` (`platform`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;