package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type Motorcycle struct {
	Id                    int
	Name                  string
	PriceToBuy            int
	PriceToUpgrade        int
	MaxLevel              int
	Engine                int
	EngineIncrement       int
	Agility               int
	AgilityIncrement      int
	Brakes                int
	BrakesIncrement       int
	Aerodynamics          int
	AerodynamicsIncrement int
}

type Ownership struct {
	Username              string
	MotorcycleId          int
	Name                  string
	Level                 int
	PriceToBuy            int
	PriceToUpgrade        int
	MaxLevel              int
	Engine                int
	EngineIncrement       int
	Agility               int
	AgilityIncrement      int
	Brakes                int
	BrakesIncrement       int
	Aerodynamics          int
	AerodynamicsIncrement int
}

type GarageDB interface {
	GetRemainingMotorcycles(username string) ([]*Motorcycle, error)
	GetUserMotorcycles(username string) ([]*Ownership, error)
	GetUserMotorcycleStats(username string, MotorcycleId int) (*Ownership, error)
	GetUserMoney(username string) (int, error)
	IncreaseUserMoney(username string, value int) error
	BuyMotorcycle(username string, MotorcycleId int) error
	UpgradeMotorcycle(username string, MotorcycleId int) error
}

// Implementation for an SQL Database
type SQL_DB struct {
	db *sql.DB
}

func NewSQL_DB(conn *sql.DB) *SQL_DB {
	return &SQL_DB{db: conn}
}

func (s *SQL_DB) GetUserMotorcycles(username string) ([]*Ownership, error) {
	// Retrieve motorcycles owned by user

	rows, err := s.db.Query("SELECT * FROM DetailedOwnership WHERE Username=?", username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var owned []*Ownership

	for rows.Next() {
		var row Ownership
		err := rows.Scan(&row.Username, &row.MotorcycleId, &row.Level,
			&row.Name, &row.PriceToBuy, &row.PriceToUpgrade, &row.MaxLevel, &row.Engine, &row.EngineIncrement,
			&row.Agility, &row.AgilityIncrement, &row.Brakes, &row.BrakesIncrement, &row.Aerodynamics, &row.AerodynamicsIncrement)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		owned = append(owned, &row)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return owned, nil
}

func (s *SQL_DB) GetUserMotorcycleStats(username string, MotorcycleId int) (*Ownership, error) {
	// Get stats of specific motorcycle

	row := s.db.QueryRow("SELECT * FROM DetailedOwnership WHERE Username=? AND MotorcycleId=?", username, MotorcycleId)
	var owned Ownership
	err := row.Scan(&owned.Username, &owned.MotorcycleId, &owned.Level,
		&owned.Name, &owned.PriceToBuy, &owned.PriceToUpgrade, &owned.MaxLevel, &owned.Engine, &owned.EngineIncrement,
		&owned.Agility, &owned.AgilityIncrement, &owned.Brakes, &owned.BrakesIncrement, &owned.Aerodynamics, &owned.AerodynamicsIncrement)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &owned, nil
}

func (s *SQL_DB) GetRemainingMotorcycles(username string) ([]*Motorcycle, error) {
	// Retrieve motorcycle not owned

	rows, err := s.db.Query("SELECT * FROM Motorcycles M WHERE M.Id NOT IN (SELECT O.MotorcycleId FROM Owners O WHERE O.Username=?)", username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var not_owned []*Motorcycle

	for rows.Next() {
		var row Motorcycle
		err := rows.Scan(&row.Id, &row.Name, &row.PriceToBuy, &row.PriceToUpgrade, &row.MaxLevel, &row.Engine, &row.EngineIncrement,
			&row.Agility, &row.AgilityIncrement, &row.Brakes, &row.BrakesIncrement, &row.Aerodynamics, &row.AerodynamicsIncrement)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		not_owned = append(not_owned, &row)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return not_owned, nil
}

func (s *SQL_DB) GetUserMoney(username string) (int, error) {
	row := s.db.QueryRow("SELECT Money FROM Users WHERE Username=?", username)

	money := 0
	row.Scan(&money)

	return money, row.Err()
}

func (s *SQL_DB) IncreaseUserMoney(username string, value int) error {
	// Increase money, used also for registration setting value=0

	if value < 0 {
		return errors.New("increase value can not be negative")
	}

	_, err := s.db.Exec("INSERT INTO Users VALUES (?, ?) ON DUPLICATE KEY UPDATE Money = Money + VALUES(Money)", username, value)

	return err
}

func (s *SQL_DB) BuyMotorcycle(username string, MotorcycleId int) error {
	// Buy motorcycle

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Begin transaction, if errors happen during execution the transaction is rolled back
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	var price int
	err = tx.QueryRow("SELECT PriceToBuy FROM Motorcycles WHERE Id=?", MotorcycleId).Scan(&price)
	if err != nil {
		log.Println(err)
		return err
	}

	var money int
	err = tx.QueryRow("SELECT Money FROM Users WHERE Username=?", username).Scan(&money)
	if err != nil {
		log.Println(err)
		return err
	}

	// Check if user can pay the price
	if money < price {
		return errors.New("not enough money to perform payment")
	}

	_, err = tx.Exec("INSERT INTO Owners (Username, MotorcycleId) VALUES (?, ?)", username, MotorcycleId)
	if err != nil {
		log.Println(err)
		return err
	}

	// Subtract price from money of user
	_, err = tx.Exec("UPDATE Users SET Money=Money-? WHERE Username=?", price, username)
	if err != nil {
		log.Println(err)
		return err
	}

	return tx.Commit()
}

func (s *SQL_DB) UpgradeMotorcycle(username string, MotorcycleId int) error {
	// Upgrade motorcycle (similar behaviour of buying)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	var price int
	err = tx.QueryRow("SELECT PriceToUpgrade FROM Motorcycles WHERE Id=?", MotorcycleId).Scan(&price)
	if err != nil {
		log.Println(err)
		return err
	}

	var money int
	err = tx.QueryRow("SELECT Money FROM Users WHERE Username=?", username).Scan(&money)
	if err != nil {
		log.Println(err)
		return err
	}

	if money < price {
		return errors.New("not enough money to perform payment")
	}

	_, err = tx.Exec("UPDATE Owners SET Level=Level+1 WHERE Username=? AND MotorcycleId=?", username, MotorcycleId)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.Exec("UPDATE Users SET Money=Money-? WHERE Username=?", price, username)
	if err != nil {
		log.Println(err)
		return err
	}

	return tx.Commit()
}
