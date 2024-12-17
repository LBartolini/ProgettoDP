DROP DATABASE IF EXISTS Auth;
CREATE DATABASE IF NOT EXISTS Auth;
USE Auth;

DROP TABLE IF EXISTS Users;
CREATE TABLE IF NOT EXISTS Users (
  Username varchar(32) NOT NULL,
  Password varchar(32) NOT NULL,
  Email varchar(64) NOT NULL,
  Phone varchar(10) NOT NULL,
  PRIMARY KEY (Username)
) ENGINE=InnoDB;

INSERT INTO Users VALUES ("pippo", "prova", "test@gmail.com", "123456789");

