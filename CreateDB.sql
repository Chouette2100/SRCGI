-- MySQL dump 10.13  Distrib 8.0.33, for Linux (x86_64)
--
-- Host: localhost    Database: iapetus_showroom
-- ------------------------------------------------------
-- Server version	8.0.33-0ubuntu0.20.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `event`
--

DROP TABLE IF EXISTS `event`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `event` (
  `eventid` char(100) NOT NULL,
  `event_name` char(200) DEFAULT NULL,
  `period` char(100) DEFAULT NULL,
  `starttime` datetime DEFAULT NULL,
  `endtime` datetime DEFAULT NULL,
  `noentry` int DEFAULT '0',
  `intervalmin` int DEFAULT '5',
  `modmin` int DEFAULT '4',
  `modsec` int DEFAULT '10',
  `fromorder` int DEFAULT '1',
  `toorder` int DEFAULT '10',
  `resethh` int DEFAULT '4',
  `resetmm` int DEFAULT '0',
  `nobasis` int DEFAULT '164614',
  `maxdsp` int DEFAULT '10',
  `cmap` int DEFAULT '1',
  `target` int DEFAULT '0',
  `rstatus` varchar(20) DEFAULT '',
  `maxpoint` int DEFAULT '0',
  `achk` int DEFAULT '0',
  `aclr` int DEFAULT '0',
  PRIMARY KEY (`eventid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `eventrank`
--

DROP TABLE IF EXISTS `eventrank`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `eventrank` (
  `eventid` char(100) NOT NULL,
  `userid` int NOT NULL,
  `ts` datetime NOT NULL,
  `listner` varchar(200) NOT NULL,
  `lastname` varchar(200) DEFAULT NULL,
  `lsnid` int DEFAULT NULL,
  `t_lsnid` int NOT NULL,
  `norder` int DEFAULT NULL,
  `nrank` int DEFAULT NULL,
  `point` int DEFAULT NULL,
  `increment` int DEFAULT NULL,
  `status` int DEFAULT '0',
  PRIMARY KEY (`eventid`,`userid`,`ts`,`t_lsnid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `eventuser`
--

DROP TABLE IF EXISTS `eventuser`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `eventuser` (
  `eventid` varchar(128) NOT NULL,
  `userno` int NOT NULL,
  `istarget` char(1) DEFAULT NULL,
  `iscntrbpoints` char(1) DEFAULT 'N',
  `graph` char(1) DEFAULT NULL,
  `color` char(20) DEFAULT NULL,
  `point` int DEFAULT NULL,
  `vld` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`eventid`,`userno`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `parameters`
--

DROP TABLE IF EXISTS `parameters`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `parameters` (
  `scope` int DEFAULT NULL,
  `name` char(20) DEFAULT NULL,
  `value` char(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `points`
--

DROP TABLE IF EXISTS `points`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `points` (
  `ts` datetime NOT NULL,
  `user_id` int NOT NULL,
  `eventid` varchar(128) NOT NULL,
  `point` int DEFAULT NULL,
  `rank` int DEFAULT NULL,
  `gap` int DEFAULT NULL,
  `pstatus` varchar(25) DEFAULT '',
  `qstatus` varchar(25) DEFAULT '',
  `ptime` varchar(25) DEFAULT '',
  `qtime` varchar(25) DEFAULT '',
  PRIMARY KEY (`eventid`,`user_id`,`ts`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `timeacq`
--

DROP TABLE IF EXISTS `timeacq`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `timeacq` (
  `idx` int NOT NULL AUTO_INCREMENT,
  `t` datetime DEFAULT NULL,
  KEY `idx` (`idx`)
) ENGINE=InnoDB AUTO_INCREMENT=19698 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `timetable`
--

DROP TABLE IF EXISTS `timetable`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `timetable` (
  `eventid` char(100) NOT NULL,
  `userid` int NOT NULL,
  `sampletm1` datetime NOT NULL,
  `sampletm2` datetime DEFAULT NULL,
  `stime` datetime DEFAULT NULL,
  `etime` datetime DEFAULT NULL,
  `target` int DEFAULT '-1',
  `totalpoint` int DEFAULT NULL,
  `earnedpoint` int DEFAULT NULL,
  `status` int DEFAULT '0',
  PRIMARY KEY (`eventid`,`userid`,`sampletm1`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `userno` int NOT NULL,
  `userid` char(100) DEFAULT NULL,
  `user_name` char(200) DEFAULT NULL,
  `longname` char(200) DEFAULT NULL,
  `shortname` char(10) DEFAULT NULL,
  `genre` char(100) DEFAULT NULL,
  `rank` char(16) DEFAULT NULL,
  `nrank` varchar(60) DEFAULT '-',
  `prank` varchar(60) DEFAULT '-',
  `level` int DEFAULT NULL,
  `followers` int DEFAULT NULL,
  `fans` int DEFAULT '-1',
  `fans_lst` int DEFAULT '-1',
  `ts` datetime DEFAULT NULL,
  `getp` char(1) DEFAULT NULL,
  `graph` char(1) DEFAULT NULL,
  `color` char(20) DEFAULT NULL,
  `currentevent` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`userno`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `userhistory`
--

DROP TABLE IF EXISTS `userhistory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `userhistory` (
  `userno` int DEFAULT NULL,
  `user_name` char(200) DEFAULT NULL,
  `genre` char(100) DEFAULT NULL,
  `rank` char(16) DEFAULT NULL,
  `nrank` varchar(60) DEFAULT '-',
  `prank` varchar(60) DEFAULT '-',
  `level` int DEFAULT NULL,
  `followers` int DEFAULT NULL,
  `fans` int DEFAULT '-1',
  `fans_lst` int DEFAULT '-1',
  `ts` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-05-13 16:28:12
