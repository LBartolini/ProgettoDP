package internal

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"testing"
)

func NewSQLConnection(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:admin@tcp(test_garage_db:3306)/Garage")
	if err != nil {
		t.Errorf("failed to connect to db: %s", err)
	}
	if err := db.Ping(); err != nil {
		t.Errorf("error pinging database: %v", err)
	}

	return db
}

func TestDBBuyMotorcycleCorrect(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)
	err := db.BuyMotorcycle("user", 2)

	if err == nil {
		return
	}

	t.Errorf("Not able to buy but instead it should be")
}

func TestDBBuyMotorcycleAlreadyOwned(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)
	err := db.BuyMotorcycle("user", 1)

	if err != nil {
		return
	}

	t.Errorf("Able to buy motorcycle but should not be (motorcycle already owned)")
}

func TestDBBuyMotorcycleNotEnoughMoney(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)
	err := db.BuyMotorcycle("foo", 1)

	if err != nil {
		return
	}

	t.Errorf("Able to buy motorcycle but should not be (not enough money)")
}

func TestDBUpgradeMotorcycleCorrect(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)
	err := db.UpgradeMotorcycle("user", 1)

	if err == nil {
		return
	}

	t.Errorf("Upgrade not performed, but should be")
}

func TestDBUpgradeMotorcycleNotEnoughMoney(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)
	err := db.UpgradeMotorcycle("test", 1)

	if err != nil {
		return
	}

	t.Errorf("Upgrade performed, but should not be (not enough money)")
}

func TestDBUpgradeMotorcycleMaxLevelReached(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)
	err := db.UpgradeMotorcycle("foo", 1)

	if err != nil {
		t.Errorf("Upgrade not performed, but should be")
		return
	}

	err = db.UpgradeMotorcycle("foo", 1)

	if err != nil {
		return
	}

	t.Errorf("Upgrade performed, but should not be (max level reached)")
}

func TestDBIncreaseUserMoney(t *testing.T) {
	conn := NewSQLConnection(t)
	defer conn.Close()

	db := NewSQL_DB(conn)
	money, err := db.GetUserMoney("test")

	if money != 0 || err != nil {
		t.Errorf("Wrong value of money or error raised")
		return
	}

	err = db.IncreaseUserMoney("test", 5)

	if err != nil {
		t.Errorf("Increase not performed but should be")
	}

	money, err = db.GetUserMoney("test")

	if money == 5 && err == nil {
		return
	}

	t.Errorf("Wrong value of money or error raised")
}
