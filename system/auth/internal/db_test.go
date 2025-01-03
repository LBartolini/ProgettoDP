package internal

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"testing"
)

func NewSQLConnection(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:admin@tcp(test_auth_db:3306)/Auth")
	if err != nil {
		t.Errorf("failed to connect to db: %s", err)
	}
	if err := db.Ping(); err != nil {
		t.Errorf("error pinging database: %v", err)
	}

	return db
}

func TestDBLoginCorrect(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)
	res, err := db.Login("test", "12345")

	if res && err == nil {
		return
	}

	t.Errorf("Login not accepted but should be (correct credentials provided)")
}

func TestDBLoginWrongPassword(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)

	res, _ := db.Login("test", "67890")

	if !res {
		return
	}

	t.Errorf("Login accepted but should not be (wrong password provided)")
}

func TestDBLoginUserNotRegistered(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)

	res, _ := db.Login("not_registered", "00000")

	if !res {
		return
	}

	t.Errorf("Login accepted but should not be (user not registered)")
}

func TestDBRegisterCorrect(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)

	res, err := db.Register("user", "password", "user@test.com", "123456789")

	if !res || err != nil {
		t.Errorf("Register not accepted but should be (correct details provided)")
	}

	res, err = db.Login("user", "password")

	if res && err == nil {
		return
	}

	t.Errorf("After Registering Login not accepted but should be (correct credentials provided)")

}

func TestDBRegisterUserAlreadyRegistered(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)

	res, _ := db.Register("foo", "123", "foo@test.com", "123456789")

	if res {
		t.Errorf("Register accepted but should not be (user already registered)")
	}

	res, _ = db.Login("foo", "123")

	if !res {
		return
	}

	t.Errorf("After Registering Login accepted but should not be (wrong password provided)")
}
