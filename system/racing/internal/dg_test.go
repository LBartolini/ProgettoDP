package internal

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"testing"
)

func NewSQLConnection(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:admin@tcp(test_racing_db:3306)/Racing?parseTime=true")
	if err != nil {
		t.Errorf("failed to connect to db: %s", err)
	}
	if err := db.Ping(); err != nil {
		t.Errorf("error pinging database: %v", err)
	}

	return db
}

func TestDBCheckIsRacing(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)
	track, err := db.CheckIsRacing("user", 3)

	if err != nil || track == "" {
		return
	}

	t.Errorf("Wrong information, user not racing but got a valid response")
}

func TestDBMatchmaking(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	stats1 := MotorcycleStats{
		Id:           1,
		Name:         "KTM",
		Level:        5,
		Engine:       100,
		Brakes:       100,
		Agility:      100,
		Aerodynamics: 100,
	}

	stats2 := MotorcycleStats{
		Id:           2,
		Name:         "Ducati",
		Level:        1,
		Engine:       5,
		Brakes:       5,
		Agility:      5,
		Aerodynamics: 5,
	}

	db := NewSQL_DB(conn)
	track, left, err := db.StartMatchmaking("user", &stats1)

	if err != nil || track == -1 || left != 1 {
		t.Errorf("Wrong response starting matchmaking motorcycle 1")
	}

	tr, err := db.CheckIsRacing("user", 1)

	if err != nil || tr != "Mugello" {
		t.Errorf("Wrong information user racing")
	}

	track, left, err = db.StartMatchmaking("user", &stats2)

	if err != nil || track == -1 || left != 0 {
		t.Errorf("Wrong response starting matchmaking motorcycle 2")
	}

	tr, err = db.CheckIsRacing("user", 2)

	if err != nil || tr != "Mugello" {
		t.Errorf("Wrong information user racing")
	}

	results, err := db.CompleteRace(1)

	if err != nil || results == nil {
		t.Log(err)
		t.Errorf("Got error while completing race but should not")
	}

	if results[0].Position == 1 && results[0].MotorcycleId == 1 && results[1].Position == 2 && results[1].MotorcycleId == 2 {
		return
	}

	t.Errorf("Wrong position of motorcycles after race")
}
