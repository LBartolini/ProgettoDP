DROP DATABASE IF EXISTS Garage;
CREATE DATABASE IF NOT EXISTS Garage;
USE Garage;

DROP TABLE IF EXISTS Users;
CREATE TABLE IF NOT EXISTS Users (
  Username varchar(32) NOT NULL,
  Money int NOT NULL,
  PRIMARY KEY (Username),
  CHECK (Money >= 0)
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
  Level int NOT NULL DEFAULT 1,
  PRIMARY KEY (Username, MotorcycleId),
  FOREIGN KEY (Username) REFERENCES Users(Username),
  FOREIGN KEY (MotorcycleId) REFERENCES Motorcycles(Id)
) ENGINE=InnoDB;

DELIMITER $$

CREATE TRIGGER OwnersBeforeUpdate
BEFORE UPDATE ON Owners
FOR EACH ROW
BEGIN
  DECLARE max_level INT;

  SELECT MaxLevel INTO max_level
  FROM Motorcycles
  WHERE Id = NEW.MotorcycleId;

  IF NEW.Level > max_level THEN
    SIGNAL SQLSTATE '45000'
    SET MESSAGE_TEXT = 'Max Level Reached';
  END IF;
END$$

DELIMITER ;

INSERT INTO Users VALUES ("Lorenzo", 200), ("Matteo", 130);
INSERT INTO Motorcycles VALUES (1, "Ducati Panigale V4", 100, 20, 15, 10, 3, 8, 2, 12, 2, 15, 5), (2, "KTM SuperDuke 1290 RR", 120, 15, 10, 16, 5, 5, 1, 10, 3, 8, 3);
INSERT INTO Owners VALUES ("Lorenzo", 1, 5), ("Matteo", 2, 4);