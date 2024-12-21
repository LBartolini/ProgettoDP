package internal

import "database/sql"

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
	IsRacing              bool
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
	GetAllMotorcycles() []*Motorcycle
	GetUserMotorcycles(username string) ([]*Ownership, error)
	GetUserMoney(username string) (int, error)
	BuyMotorcycle(username string, MotorcycleId int) error
	UpgradeMotorcycle(username string, MotorcycleId int) error
	StartRace(username string, MotorcycleId int) error
	EndRace(username string, MotorcycleId int) error
}

type SQL_DB struct {
	db *sql.DB
}

func NewSQL_DB(conn *sql.DB) *SQL_DB {
	return &SQL_DB{db: conn}
}

func (s *SQL_DB) GetAllMotorcycles() []*Motorcycle {
	return nil
}

func (s *SQL_DB) GetUserMotorcycles(username string) ([]*Ownership, error) {
	return nil, nil
}

func (s *SQL_DB) GetUserMoney(username string) (int, error) {
	return 0, nil
}

func (s *SQL_DB) BuyMotorcycle(username string, MotorcycleId int) error {
	return nil
}

func (s *SQL_DB) UpgradeMotorcycle(username string, MotorcycleId int) error {
	return nil
}

func (s *SQL_DB) StartRace(username string, MotorcycleId int) error {
	return nil
}

func (s *SQL_DB) EndRace(username string, MotorcycleId int) error {
	return nil
}
