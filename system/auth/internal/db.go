package internal

import (
	"database/sql"
	"log"
)

// Interface with the Auth Database
type AuthDB interface {
	Login(username string, password string) (bool, error)
	Register(username string, password string, email string, phone string) (bool, error)
}

// Implementation for an SQL Database
type SQL_DB struct {
	db *sql.DB
}

func NewSQL_DB(conn *sql.DB) *SQL_DB {
	return &SQL_DB{db: conn}
}

func (s *SQL_DB) Login(username string, password string) (bool, error) {
	// Perform Login inside the DB
	// if credentials are not correct no row will be retrieved resulting in an erroneous login action

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
	// Perform Registration inside the DB
	// if username is already present the registration fails, the error is picked when executing the statement and also by checking the affected rows

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
