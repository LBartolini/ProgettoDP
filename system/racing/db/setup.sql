DROP DATABASE IF EXISTS Racing;
CREATE DATABASE IF NOT EXISTS Racing;
USE Racing;

DROP TABLE IF EXISTS Tracks;
CREATE TABLE IF NOT EXISTS Tracks (
  Id int NOT NULL,
  Name varchar(32) NOT NULL,
  MaxMotorcycles int NOT NULL,
  EngineValue int NOT NULL,
  BrakesValue int NOT NULL,
  AgilityValue int NOT NULL,
  AerodynamicsValue int NOT NULL,
  PRIMARY KEY (Id),
) ENGINE=InnoDB;

DROP TABLE IF EXISTS Matchmaking;
CREATE TABLE IF NOT EXISTS Matchmaking (
  PlayerUsername varchar(32) NOT NULL,
  MotorcycleId int NOT NULL,
  TrackId int NOT NULL,
  MotorycleName varchar(32) NOT NULL,
  MotorycleLevel int NOT NULL,
  MotorcycleEngine int NOT NULL,
  MotorcycleBrakes int NOT NULL,
  MotorcycleAgility int NOT NULL,
  MotorcycleAerodynamics int NOT NULL,
  PRIMARY KEY (PlayerUsername, MotorcycleId, TrackId),
  FOREIGN KEY TrackId REFERENCES Tracks(Id)
) ENGINE=InnoDB;

CREATE VIEW MatchmakingTotalPower AS
SELECT PlayerUsername, MotorocycleId, MotorcycleName, MotorcycleLevel, TrackId, 
    MaxMotorcycles, (MotorcycleEngine*EngineValue + MotorcycleAgility*AgilityValue + MotorcycleBrakes*BrakesValue + MotorcycleAerodynamics*AerodynamicsValue) as Power
FROM Matchmaking M
INNER JOIN Tracks T ON M.TrackId=T.Id
GROUP BY TrackId;

DELIMITER $$

CREATE TRIGGER MatchmakingBeforeUpdate
BEFORE UPDATE ON MatchmakingBeforeUpdate
FOR EACH ROW
BEGIN
  DECLARE max_motorcycles INT;

  SELECT MaxMotorcycles INTO max_motorcycles
  FROM Tracks
  WHERE Id = NEW.TrackId;

  DECLARE motorcycle_racing INT;

  SELECT COUNT(*) INTO motorcycle_racing
  FROM Matchmaking
  WHERE TrackId = NEW.TrackId;

  IF motorcycle_racing > max_motorcycles THEN
    SIGNAL SQLSTATE '45000'
    SET MESSAGE_TEXT = 'Max Motorcycle for this track reached';
  END IF;
END$$

DELIMITER ;

DROP TABLE IF EXISTS History;
CREATE TABLE IF NOT EXISTS History (
  RaceId int AUTO_INCREMENT NOT NULL,
  Position int NOT NULL,
  TotalMotorcycles int NOT NULL,
  PlayerUsername varchar(32) NOT NULL,
  TrackName varchar(32) NOT NULL,
  MotorycleName varchar(32) NOT NULL,
  MotorycleLevel int NOT NULL,
  PRIMARY KEY (RaceId)
) ENGINE=InnoDB;

