package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

type MockStorage struct {
	employee *Employee
	err      error
}

func (m *MockStorage) GetEmployeeByID(id int) (*Employee, error) {
	return m.employee, m.err
}

func (m *MockStorage) CreateEmployee(e *Employee) error        { return nil }
func (m *MockStorage) GetEmployees() ([]*Employee, error)      { return nil, nil }
func (m *MockStorage) DeleteEmployee(id int) error             { return nil }
func (m *MockStorage) UpdateEmployee(employee *Employee) error { return nil }
func (m *MockStorage) GetCities() ([]*City, error)             { return nil, nil }
func (m *MockStorage) CityExists(id int) (bool, error)         { return true, nil }

func TestGetEmployeeByID(t *testing.T) {
	storage := &MockStorage{
		employee: &Employee{
			ID:     1,
			Name:   "Test",
			Phone:  "1234567890",
			CityID: 1,
		},
		err: nil,
	}
	server := NewApiServer(":8081", storage)

	router := mux.NewRouter()
	router.HandleFunc("/employee/{id}", func(w http.ResponseWriter, r *http.Request) {
		err := server.handleGetEmployeeByID(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	request := httptest.NewRequest("GET", "/employee/1", nil)
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":1,"name":"Test","phone":"1234567890","city_id":1}`
	actual := strings.TrimSpace(responseRecorder.Body.String())
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestCreateEmployee(t *testing.T) {
	storage := &MockStorage{
		employee: nil,
		err:      nil,
	}
	server := NewApiServer(":8081", storage)

	router := mux.NewRouter()
	router.HandleFunc("/employee", func(w http.ResponseWriter, r *http.Request) {
		err := server.handleEmployee(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	employee := `{"name":"Test","phone":"87076219166","city_id":1}`
	request := httptest.NewRequest("POST", "/employee", strings.NewReader(employee))
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expected := `{"id":0,"name":"Test","phone":"87076219166","city_id":1}`
	actual := strings.TrimSpace(responseRecorder.Body.String())
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestDeleteEmployee(t *testing.T) {
	storage := &MockStorage{
		employee: nil,
		err:      nil,
	}
	server := NewApiServer(":8081", storage)

	router := mux.NewRouter()
	router.HandleFunc("/employee/{id}", func(w http.ResponseWriter, r *http.Request) {
		err := server.handleGetEmployeeByID(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	request := httptest.NewRequest("DELETE", "/employee/1", nil)
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"deleted"}`
	actual := strings.TrimSpace(responseRecorder.Body.String())
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}
