package main

type CreateEmployeeRequest struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	CityID int    `json:"city_id"`
}

type Employee struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	CityID int    `json:"city_id"`
}

type City struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewEmployee(name, phone string, city_id int) *Employee {
	return &Employee{
		Name:   name,
		Phone:  phone,
		CityID: city_id,
	}
}
