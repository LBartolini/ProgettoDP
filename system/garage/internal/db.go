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
	GetAllMotorcycles() ([]*Motorcycle, error)
	GetRemainingMotorcycles(username string) ([]*Motorcycle, error)
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

func (s *SQL_DB) GetAllMotorcycles() ([]*Motorcycle, error) { // TODO Maybe Remove Because Unused
	return nil, nil
}

func (s *SQL_DB) GetUserMotorcycles(username string) ([]*Ownership, error) {
	rows, err := s.db.Query("SELECT * FROM Owners O INNER JOIN Motorcycles M ON O.MotorcycleId=M.Id WHERE O.Username=?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var owned []*Ownership

	for rows.Next() {
		var row Ownership
		err := rows.Scan(&row.Username, &row.MotorcycleId, &row.Level, &row.IsRacing, &row.MotorcycleId,
			&row.Name, &row.PriceToBuy, &row.PriceToUpgrade, &row.MaxLevel, &row.Engine, &row.EngineIncrement,
			&row.Agility, &row.AgilityIncrement, &row.Brakes, &row.BrakesIncrement, &row.Aerodynamics, &row.AerodynamicsIncrement)
		if err != nil {
			return nil, err
		}
		owned = append(owned, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return owned, nil
}

func (s *SQL_DB) GetRemainingMotorcycles(username string) ([]*Motorcycle, error) {
	rows, err := s.db.Query("SELECT * FROM Motorcycles M WHERE M.Id NOT IN (SELECT O.MotorcycleId FROM Owners O WHERE O.Username=?)", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var not_owned []*Motorcycle

	for rows.Next() {
		var row Motorcycle
		err := rows.Scan(&row.Id, &row.Name, &row.PriceToBuy, &row.PriceToUpgrade, &row.MaxLevel, &row.Engine, &row.EngineIncrement,
			&row.Agility, &row.AgilityIncrement, &row.Brakes, &row.BrakesIncrement, &row.Aerodynamics, &row.AerodynamicsIncrement)
		if err != nil {
			return nil, err
		}
		not_owned = append(not_owned, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return not_owned, nil
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
