-- MySQL dump 10.17  Distrib 10.3.23-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: db_fishrider
-- ------------------------------------------------------
-- Server version	10.3.23-MariaDB-0+deb10u1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `admin_master`
--

DROP TABLE IF EXISTS `admin_master`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `admin_master` (
  `admin_uname` varchar(50) NOT NULL DEFAULT '0',
  `admin_passwd` varchar(70) NOT NULL,
  `admin_email` varchar(50) DEFAULT NULL,
  `admin_sesid` varchar(100) DEFAULT '0',
  PRIMARY KEY (`admin_uname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admin_master`
--

LOCK TABLES `admin_master` WRITE;
/*!40000 ALTER TABLE `admin_master` DISABLE KEYS */;
INSERT INTO `admin_master` VALUES ('martin','$2a$14$p1UPB3sbG43d.AHX8fq/1u3X1JtF4wLjSSwZwnzaR/3xBRfczsgBe',NULL,'0607faaf-775b-449c-bf18-023111da8cc4');
/*!40000 ALTER TABLE `admin_master` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `country_codes`
--

DROP TABLE IF EXISTS `country_codes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_codes` (
  `country` varchar(50) DEFAULT NULL,
  `code` varchar(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `country_codes`
--

LOCK TABLES `country_codes` WRITE;
/*!40000 ALTER TABLE `country_codes` DISABLE KEYS */;
/*!40000 ALTER TABLE `country_codes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cust_master`
--

DROP TABLE IF EXISTS `cust_master`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cust_master` (
  `cust_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `cust_mob` bigint(20) NOT NULL DEFAULT 0,
  `cust_name` varchar(50) NOT NULL DEFAULT '',
  `cust_adr1` varchar(50) NOT NULL DEFAULT '',
  `cust_adr2` varchar(50) NOT NULL DEFAULT '',
  `cust_lmark` varchar(50) DEFAULT '',
  `cust_sesid` varchar(100) DEFAULT '0',
  PRIMARY KEY (`cust_id`),
  UNIQUE KEY `cust_mob` (`cust_mob`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cust_master`
--

LOCK TABLES `cust_master` WRITE;
/*!40000 ALTER TABLE `cust_master` DISABLE KEYS */;
INSERT INTO `cust_master` VALUES (1,9846500400,'Arun Kumar','Ambady(Kannimittom)','Velorvattom','Near Temple','b51c780b-01f2-47ff-876c-37c49053cb37'),(2,9846500401,'Ajith','Ajith Nivas','Railway Station Road','','a1cce203-1f7a-4d9a-ba54-9c9a41231844'),(3,9846500402,'Alex','House No 34','Brigade road','','0'),(4,9846500403,'Aravind ','Athira','kannan Kara','','0'),(5,9846500404,'Syam','Kalathil','Manoramakkavala','','0'),(6,9846500405,'Sathesh','House No 35','Thannermukkam','','0'),(7,9846500406,'Santhosh Joy','Pulloruhtikari','GreenLand Street','',''),(8,9846500407,'Anil','Tharayil House','Velorvattom','',''),(9,9846500410,'dfdf','dfd','dfd','','');
/*!40000 ALTER TABLE `cust_master` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `item_master`
--

DROP TABLE IF EXISTS `item_master`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `item_master` (
  `item_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `item_desc` varchar(50) NOT NULL DEFAULT '0',
  `item_stock` int(11) DEFAULT 0,
  `item_sel_price` float DEFAULT 0,
  `item_buy_price` float DEFAULT 0,
  `item_unit` varchar(10) NOT NULL,
  PRIMARY KEY (`item_id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `item_master`
--

LOCK TABLES `item_master` WRITE;
/*!40000 ALTER TABLE `item_master` DISABLE KEYS */;
INSERT INTO `item_master` VALUES (1,'Trevally-Fresh',97,40.5,20,'Number'),(2,'White Sardine',110,60.5,50,'Kg'),(3,'King Fish',0,0,0,'Kg'),(4,'Prawns',19,40,30,'Kg'),(5,'King Fish(Medium)',87,60,50,'Number'),(13,'Varala-Small',100,60,50,'Number');
/*!40000 ALTER TABLE `item_master` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `order_detail`
--

DROP TABLE IF EXISTS `order_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `order_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `order_id` int(10) unsigned NOT NULL,
  `item_id` int(10) unsigned NOT NULL,
  `item_qty` int(11) NOT NULL,
  `item_sold_price` float DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `FK__order_master` (`order_id`),
  KEY `FK_order_detail_item_master` (`item_id`),
  CONSTRAINT `FK__order_master` FOREIGN KEY (`order_id`) REFERENCES `order_master` (`order_id`),
  CONSTRAINT `FK_order_detail_item_master` FOREIGN KEY (`item_id`) REFERENCES `item_master` (`item_id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `order_detail`
--

LOCK TABLES `order_detail` WRITE;
/*!40000 ALTER TABLE `order_detail` DISABLE KEYS */;
INSERT INTO `order_detail` VALUES (1,82,4,4,NULL),(2,83,4,7,NULL),(3,83,2,9,NULL),(4,84,5,7,NULL),(5,85,5,6,NULL),(6,86,5,5,NULL),(7,87,5,8,NULL),(8,87,2,10,NULL),(9,88,1,1,NULL),(10,89,5,5,NULL),(11,89,1,2,NULL);
/*!40000 ALTER TABLE `order_detail` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `order_master`
--

DROP TABLE IF EXISTS `order_master`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `order_master` (
  `order_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `order_date` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `cust_id` int(10) unsigned NOT NULL,
  `order_amt` float unsigned NOT NULL DEFAULT 0,
  `order_status` varchar(50) NOT NULL DEFAULT 'pending',
  `p_mode` varchar(50) NOT NULL DEFAULT 'cod',
  PRIMARY KEY (`order_id`),
  KEY `FK_order_master_cust_master` (`cust_id`),
  CONSTRAINT `FK_order_master_cust_master` FOREIGN KEY (`cust_id`) REFERENCES `cust_master` (`cust_id`)
) ENGINE=InnoDB AUTO_INCREMENT=90 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `order_master`
--

LOCK TABLES `order_master` WRITE;
/*!40000 ALTER TABLE `order_master` DISABLE KEYS */;
INSERT INTO `order_master` VALUES (82,'2020-10-06 03:38:51',2,160,'approved','cod'),(83,'2020-10-05 08:13:36',2,820,'pending','online'),(84,'2020-10-05 08:23:07',2,560,'pending','cod'),(85,'2020-10-05 08:26:01',2,480,'pending','cod'),(86,'2020-10-05 08:27:02',1,400,'pending','cod'),(87,'2020-10-05 18:40:54',1,1085,'pending','cod'),(88,'2020-10-06 03:43:35',1,40.5,'approved','cod'),(89,'2020-10-06 06:09:47',1,381,'approved','online');
/*!40000 ALTER TABLE `order_master` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `temp_cart`
--

DROP TABLE IF EXISTS `temp_cart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `temp_cart` (
  `temp_sesid` varchar(100) DEFAULT NULL,
  `item_id` int(11) unsigned DEFAULT NULL,
  `item_qty` int(11) DEFAULT NULL,
  `time_stamp` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  KEY `FK_temp_cart_item_master` (`item_id`),
  CONSTRAINT `FK_temp_cart_item_master` FOREIGN KEY (`item_id`) REFERENCES `item_master` (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `temp_cart`
--

LOCK TABLES `temp_cart` WRITE;
/*!40000 ALTER TABLE `temp_cart` DISABLE KEYS */;
/*!40000 ALTER TABLE `temp_cart` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-10-11 11:06:09
