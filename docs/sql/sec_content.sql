/*
 Navicat Premium Data Transfer

 Source Server         : 192.168.123.96
 Source Server Type    : MySQL
 Source Server Version : 50651 (5.6.51)
 Source Host           : 192.168.123.96:3306
 Source Schema         : blade_ops

 Target Server Type    : MySQL
 Target Server Version : 50651 (5.6.51)
 File Encoding         : 65001

 Date: 08/03/2024 10:18:47
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sec_content
-- ----------------------------
DROP TABLE IF EXISTS `security_content`;
CREATE TABLE `security_content` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` varchar(64) NOT NULL COMMENT '实验唯一标识\n',
  `pid` int(11) DEFAULT NULL COMMENT '实验执行进程号\n',
  `result` text COMMENT '实验执行结果',
  `start_at` datetime DEFAULT NULL COMMENT '实验开始时间',
  `end_at` datetime DEFAULT NULL COMMENT '实验结束时间',
  `is_deleted` int(1) unsigned DEFAULT '0' COMMENT '记录是否删除',
  `is_end` tinyint(1) NOT NULL COMMENT '实验执行是否结束',
  `is_destroyed` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '实验是否被人为终止',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid` (`uid`) USING HASH
) ENGINE=InnoDB AUTO_INCREMENT=63 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
