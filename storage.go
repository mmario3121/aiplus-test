package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateEmployee(employee *Employee) error
	GetEmployees() ([]*Employee, error)
	GetEmployeeByID(id int) (*Employee, error)
	DeleteEmployee(id int) error
	UpdateEmployee(employee *Employee) error
	GetCities() ([]*City, error)
	CityExists(id int) (bool, error)
}

type PostgresStore struct {
	db *sql.DB
}

type Scanner interface {
	Scan(rows *sql.Rows) error
}

func WaitPostgres(psqlInfo string) (*sql.DB, error) {
	for i := 0; i < 5; i++ {
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Printf("Error while opening database: %v. Retrying...", err)
			time.Sleep(time.Duration(1<<i) * time.Second)
			continue
		}
		err = db.Ping()
		if err == nil {
			log.Println("Connected to database")
			return db, nil
		}
		log.Printf("Error connecting to database: %v. Retrying...", err)
		time.Sleep(time.Duration(1<<i) * time.Second)
	}
	return nil, fmt.Errorf("Error connecting to database")
}

func NewPostgresStore() (*PostgresStore, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := WaitPostgres(psqlInfo)

	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	err := s.CreateEmployeeTable()
	if err != nil {
		return err
	}
	err = s.CreateCityTable()
	if err != nil {
		return err
	}
	err = s.SeedCities()
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) CreateEmployeeTable() error {
	query := `CREATE TABLE IF NOT EXISTS employees (
		id SERIAL PRIMARY KEY,
		name TEXT,
		phone TEXT,
		city_id INT
		)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateCityTable() error {
	query := `CREATE TABLE IF NOT EXISTS cities (
		id SERIAL PRIMARY KEY,
		name TEXT
		)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) SeedCities() error {
	cities := []string{"Almaty", "Astana", "Taraz", "Jerusalem", "Kyiv"}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, city := range cities {
		tx.Exec("insert into cities (name) values ($1)", city)
	}
	log.Printf("Seeded Cities")
	return tx.Commit()
}

func (s *PostgresStore) CreateEmployee(employee *Employee) error {
	query := `insert into employees (name, phone, city_id) values ($1, $2, $3) returning id`
	_, err := s.db.Query(
		query,
		employee.Name,
		employee.Phone,
		employee.CityID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetEmployees() ([]*Employee, error) {
	rows, err := s.db.Query("select * from employees")
	if err != nil {
		return nil, err
	}
	employees := []*Employee{}
	for rows.Next() {
		employee := &Employee{}
		err := scanInto(employee, rows)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	defer rows.Close()
	return employees, nil
}

func (s *PostgresStore) UpdateEmployee(employee *Employee) error {
	return nil
}

func (s *PostgresStore) DeleteEmployee(id int) error {
	_, err := s.db.Exec("delete from employees where id = $1", id)
	return err
}

func (s *PostgresStore) GetEmployeeByID(id int) (*Employee, error) {
	rows, err := s.db.Query("select * from employees where id = $1", id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		employee := &Employee{}
		return employee, scanInto(employee, rows)
	}
	return nil, fmt.Errorf("employee %d not found", id)
}

func (s *PostgresStore) GetCities() ([]*City, error) {
	rows, err := s.db.Query("select * from cities")
	if err != nil {
		return nil, err
	}
	cities := []*City{}
	for rows.Next() {
		city := &City{}
		err := scanInto(city, rows)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}
	defer rows.Close()
	return cities, nil
}

func (s *PostgresStore) CityExists(id int) (bool, error) {
	query := `select exists(select 1 from cities where id = $1)`
	var exists bool
	err := s.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (e *Employee) Scan(rows *sql.Rows) error {
	return rows.Scan(&e.ID, &e.Name, &e.Phone, &e.CityID)
}

func (c *City) Scan(rows *sql.Rows) error {
	return rows.Scan(&c.ID, &c.Name)
}

func scanInto(s Scanner, rows *sql.Rows) error {
	return s.Scan(rows)
}
