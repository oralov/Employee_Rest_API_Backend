package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// A struct that is used to create models(objects) of employees
type Employee struct {
	Id     int
	Name   string
	Age    int
	Salary int
}

// this function requests the entire list of employees from the database table
func getEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	json.NewEncoder(w).Encode(GetItems())

}

// this  function handles the request to delete an item from a database by id
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	fmt.Println(id)
	DeleteItem(id)
}

// this function handles the request to create a new employee in the employees table
func createEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	var employee Employee
	_ = json.NewDecoder(r.Body).Decode(&employee)
	PutItem(employee)
	json.NewEncoder(w).Encode(employee)
}
