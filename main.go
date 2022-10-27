package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	dynamo = ConnectDynamo()
}

func main() {
	//createTable()
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/employees/", getEmployees).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/employees/{id}", deleteEmployee).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/v1/employees/", createEmployee).Methods("POST", "GET")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", router))
}
