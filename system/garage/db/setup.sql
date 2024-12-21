DROP DATABASE IF EXISTS Garage;
CREATE DATABASE IF NOT EXISTS Garage;
USE Garage;

DROP TABLE IF EXISTS Users;
CREATE TABLE IF NOT EXISTS Users (
  Username varchar(32) NOT NULL,
  Money int NOT NULL,
  PRIMARY KEY (Username)
) ENGINE=InnoDB;

DROP TABLE IF EXISTS Motorcycles;
CREATE TABLE IF NOT EXISTS Motorcycles (
  Id int NOT NULL AUTO_INCREMENT,
  Name varchar(32) NOT NULL,
  PriceToBuy int NOT NULL,
  PriceToUpgrade int NOT NULL,
  MaxLevel int NOT NULL,
  Engine int NOT NULL,
  EngineIncrement int NOT NULL,
  Agility int NOT NULL,
  AgilityIncrement int NOT NULL,
  Brakes int NOT NULL,
  BrakesIncrement int NOT NULL,
  Aerodynamics int NOT NULL,
  AerodynamicsIncrement int NOT NULL,
  PRIMARY KEY (Id)
) ENGINE=InnoDB;

DROP TABLE IF EXISTS Owners;
CREATE TABLE IF NOT EXISTS Owners (
  Username varchar(32) NOT NULL,
  MotorcycleId int NOT NULL,
  Level int NOT NULL,
  IsRacing boolean NOT NULL,
  PRIMARY KEY (Username, MotorcycleId),
  FOREIGN KEY (Username) REFERENCES Users(Username),
  FOREIGN KEY (MotorcycleId) REFERENCES Motorcycles(Id)
) ENGINE=InnoDB;

INSERT INTO Users VALUES ("Lorenzo", 100), ("Matteo", 50);
INSERT INTO Motorcycles VALUES (1, "Ducati Panigale V4", 100, 20, 15, 10, 3, 8, 2, 12, 2, 15, 5), (2, "KTM SuperDuke 1290 RR", 120, 15, 10, 16, 5, 5, 1, 10, 3, 8, 3);
INSERT INTO Owners VALUES ("Lorenzo", 1, 5, false), ("Matteo", 2, 4, false);