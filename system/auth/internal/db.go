package internal

import (
	"database/sql"
	"log"
)

type AuthDB interface {
	Login(username string, password string) (bool, error)
	Register(username string, password string, email string, phone string) (bool, error)
}

type SQL_DB struct {
	db *sql.DB
}

func NewSQL_DB(conn *sql.DB) *SQL_DB {
	return &SQL_DB{db: conn}
}

func (s *SQL_DB) Login(username string, password string) (bool, error) {
	stmt, err := s.db.Prepare("SELECT Username FROM Users WHERE Username=? AND Password=?")
	if err != nil {
		log.Println(err)
		return false, err
	}

	var temp string
	err = stmt.QueryRow(username, password).Scan(&temp)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func (s *SQL_DB) Register(username string, password string, email string, phone string) (bool, error) {
	stmt, err := s.db.Prepare("INSERT INTO Users VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(username, password, email, phone)
	if err != nil {
		log.Println(err)
		return false, err
	}

	rows_affected, err := res.RowsAffected()

	return rows_affected != 0, err
}
