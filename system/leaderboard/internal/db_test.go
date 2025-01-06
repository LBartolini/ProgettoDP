package internal

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"testing"
)

func NewSQLConnection(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:admin@tcp(test_leaderboard_db:3306)/Leaderboard")
	if err != nil {
		t.Errorf("failed to connect to db: %s", err)
	}
	if err := db.Ping(); err != nil {
		t.Errorf("error pinging database: %v", err)
	}

	return db
}

func TestDBIncrementPoints(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)
	info, err := db.GetUserInfo("user")

	if err != nil || info.points != 10 {
		t.Errorf("Unable to get correct leaderboard info")
	}

	err = db.IncrementPoints("user", 10)

	if err != nil {
		t.Errorf("Unable to increment user points")
	}

	info, err = db.GetUserInfo("user")

	if err != nil || info.points != 20 {
		t.Errorf("Unable to get correct leaderboard info")
	}

}
