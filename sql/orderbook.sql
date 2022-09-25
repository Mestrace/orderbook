CREATE DATABASE `db_orderbook_main`; 

USE `db_orderbook_main`; 

CREATE TABLE tab_orderbook_exchange_metadata 
  ( 
     `id`         BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT, 
     `exchange`   VARCHAR(32) NOT NULL, 
     `metadata`   TEXT NOT NULL, 
     `created_at` DATETIME NOT NULL, 
     `updated_at` DATETIME NOT NULL, 
     UNIQUE KEY `uniq_exchange` (`exchange`) 
  ) 
AUTO_INCREMENT = 1000, 
ENGINE=InnoDB, 
DEFAULT CHARSET = utf8mb4; 