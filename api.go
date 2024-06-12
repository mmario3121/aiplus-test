package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewApiServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/employee", makeHTTPHandleFunc(s.handleEmployee))
	router.HandleFunc("/employee/{id}", makeHTTPHandleFunc(s.handleGetEmployeeByID))

	router.HandleFunc("/city", makeHTTPHandleFunc(s.handleCities))

	err := http.ListenAndServe(s.listenAddr, router)

	if err != nil {
		log.Fatal(err)
	}
}

func (s *APIServer) handleEmployee(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handleGetEmployee(w, r)
	}
	if r.Method == http.MethodPost {
		return s.handleCreateEmployee(w, r)
	}
	return fmt.Errorf("method not allowed")
}

func (s *APIServer) handleCities(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handleGetCities(w, r)
	}
	return fmt.Errorf("method not allowed")
}

func (s *APIServer) handleGetEmployeeByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		id, err := getID(r)

		if err != nil {
			return err
		}

		employee, err := s.store.GetEmployeeByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, employee)
	}
	if r.Method == http.MethodDelete {
		return s.handleDeleteEmployee(w, r)
	}
	return fmt.Errorf("method not allowed")
}

func (s *APIServer) handleGetEmployee(w http.ResponseWriter, r *http.Request) error {
	employees, err := s.store.GetEmployees()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, employees)

}

//handleGetCities returns all cities

func (s *APIServer) handleGetCities(w http.ResponseWriter, r *http.Request) error {
	cities, err := s.store.GetCities()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, cities)
}

func (s *APIServer) handleCreateEmployee(w http.ResponseWriter, r *http.Request) error {
	createEmployeeReq := new(CreateEmployeeRequest)
	if err := json.NewDecoder(r.Body).Decode(createEmployeeReq); err != nil {
		return err
	}

	employee := NewEmployee(createEmployeeReq.Name, createEmployeeReq.Phone, createEmployeeReq.CityID)

	if err := ValidateEmployee(s.store, employee); err != nil {
		return err
	}

	if err := s.store.CreateEmployee(employee); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, employee)
}

func (s *APIServer) handleDeleteEmployee(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)

	if err != nil {
		return err
	}

	if err := s.store.DeleteEmployee(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idString := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idString)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idString)
	}
	return id, nil
}

func ValidateEmployee(s Storage, employee *Employee) error {
	if len(employee.Name) >= 30 || len(employee.Name) <= 3 {
		return fmt.Errorf("name is invalid")
	}

	var validPhone = regexp.MustCompile(`^\d{11}$`)
	if !validPhone.MatchString(employee.Phone) {
		return fmt.Errorf("phone is invalid")
	}

	exists, err := s.CityExists(employee.CityID)
	if err != nil {
		return fmt.Errorf("error checking city")
	}
	if !exists {
		return fmt.Errorf("city invalid")
	}

	return nil
}
