DROP DATABASE IF EXISTS Leaderboard;
CREATE DATABASE IF NOT EXISTS Leaderboard;
USE Leaderboard;

DROP TABLE IF EXISTS Users;
CREATE TABLE IF NOT EXISTS Users (
  Username varchar(32) NOT NULL,
  Points int NOT NULL,
  PRIMARY KEY (Username)
) ENGINE=InnoDB;

CREATE VIEW RankedUsers AS
SELECT Username, Points, RANK() OVER (ORDER BY Points DESC) AS Position
FROM Users;

INSERT INTO Users VALUES ("user", 10);