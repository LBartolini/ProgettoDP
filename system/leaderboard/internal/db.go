package internal

import (
	"database/sql"
)

type LeaderboardInfo struct {
	username string
	points   int32
	position int32
}

type LeaderboardDB interface {
	GetLeaderboard() []LeaderboardInfo
	GetUserInfo(username string) (*LeaderboardInfo, error)
	IncrementPoints(username string, points int) error
}

type SQL_DB struct {
	db *sql.DB
}

func NewSQL_DB(conn *sql.DB) *SQL_DB {
	return &SQL_DB{db: conn}
}

func (s *SQL_DB) GetLeaderboard() []LeaderboardInfo {
	rows, err := s.db.Query("SELECT * FROM RankedUsers ORDER BY Position ASC")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var info []LeaderboardInfo

	for rows.Next() {
		var row LeaderboardInfo
		if err := rows.Scan(&row.username, &row.points, &row.position); err != nil {
			return nil
		}
		info = append(info, row)
	}

	if err := rows.Err(); err != nil {
		return nil
	}

	return info
}

func (s *SQL_DB) GetUserInfo(username string) (*LeaderboardInfo, error) {
	stmt, err := s.db.Prepare("SELECT Username, Points, Position FROM RankedUsers WHERE Username=?")

	if err != nil {
		return nil, err
	}

	var info LeaderboardInfo
	err = stmt.QueryRow(username).Scan(&info.username, &info.points, &info.position)

	return &info, err
}

func (s *SQL_DB) IncrementPoints(username string, points int) error {
	stmt, err := s.db.Prepare("INSERT INTO Users VALUES (?, ?) ON DUPLICATE KEY UPDATE Points = Points + VALUES(Points)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, points)

	return err
}
