package internal

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type MotorcycleStats struct {
	Id           int
	Name         string
	Level        int
	Engine       int
	Brakes       int
	Agility      int
	Aerodynamics int
}

type RaceResult struct {
	Username         string
	MotorcycleId     int
	MotorcycleName   string
	MotorcycleLevel  int
	Position         int
	TotalMotorcycles int
	TrackName        string
	Time             time.Time
}

type RacingDB interface {
	StartMatchmaking(username string, stats *MotorcycleStats) (track int, left int, e error)
	CompleteRace(track int) ([]RaceResult, error)
	CheckIsRacing(username string, MotorcycleId int) (track string, e error)
	GetHistory(username string) ([]RaceResult, error)
}

type SQL_DB struct {
	db *sql.DB
}

func NewSQL_DB(conn *sql.DB) *SQL_DB {
	return &SQL_DB{db: conn}
}

func (s *SQL_DB) StartMatchmaking(username string, stats *MotorcycleStats) (track int, left int, e error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return -1, -1, err
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT Id FROM Tracks ORDER BY RAND() LIMIT 1")
	row.Scan(&track)

	_, e = tx.Exec("INSERT INTO Matchmaking VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", username, stats.Id, track, stats.Name, stats.Level, stats.Engine, stats.Brakes, stats.Agility, stats.Aerodynamics)

	row = tx.QueryRow("SELECT FreeSlots FROM DetailedMatchmaking WHERE TrackId=? LIMIT 1", track)
	row.Scan(&left)

	return track, left, tx.Commit()
}

func (s *SQL_DB) CompleteRace(track int) ([]RaceResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT PlayerUsername, MotorcycleId, Position, MaxMotorcycles, MotorcycleLevel, MotorcycleName, Trackname, CURRENT_TIMESTAMP AS Time FROM DetailedMatchmaking WHERE TrackId=?", track)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var res []RaceResult
	for rows.Next() {
		var result RaceResult
		err = rows.Scan(&result.Username, &result.MotorcycleId, &result.Position, &result.TotalMotorcycles, &result.MotorcycleLevel, &result.MotorcycleName, &result.TrackName, &result.Time)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		res = append(res, result)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	rows.Close()

	_, err = tx.Exec("INSERT INTO History (Position, TotalMotorcycles, PlayerUsername, TrackName, MotorcycleName, MotorcycleLevel) SELECT Position, MaxMotorcycles, PlayerUsername, TrackName, MotorcycleName, MotorcycleLevel FROM DetailedMatchmaking WHERE TrackId=?", track)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = tx.Exec("DELETE FROM Matchmaking WHERE TrackId=?", track)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res, tx.Commit()
}

func (s *SQL_DB) CheckIsRacing(username string, MotorcycleId int) (track string, e error) {
	track = ""
	row := s.db.QueryRow("SELECT TrackName FROM DetailedMatchmaking WHERE PlayerUsername=? AND MotorcycleId=?", username, MotorcycleId)
	e = row.Scan(&track)

	return track, e
}

func (s *SQL_DB) GetHistory(username string) ([]RaceResult, error) {
	rows, err := s.db.Query("SELECT Position, TotalMotorcycles, PlayerUsername, TrackName, MotorcycleName, MotorcycleLevel, Time FROM History WHERE PlayerUsername=? ORDER BY RaceId DESC", username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var history []RaceResult
	for rows.Next() {
		var result RaceResult
		err = rows.Scan(&result.Position, &result.TotalMotorcycles, &result.Username, &result.TrackName, &result.MotorcycleName, &result.MotorcycleLevel, &result.Time)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		history = append(history, result)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return history, err
}
