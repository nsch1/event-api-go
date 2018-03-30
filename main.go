package main

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"encoding/json"
)

type Event struct {
	gorm.Model
	Name string
}

var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=secret sslmode=disable")
	if err != nil {
		panic("oops")
	}
	defer db.Close()

	db.AutoMigrate(&Event{})

	router := mux.NewRouter()
	router.HandleFunc("/events", GetEvents).Methods("GET")
	router.HandleFunc("/events", CreateEvent).Methods("POST")
	router.HandleFunc("/events/{id}", GetEvent).Methods("GET")

	log.Fatal(http.ListenAndServe(":4008", router))
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	var events []Event
	db.Find(&events)
	json.NewEncoder(w).Encode(&events)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	json.NewDecoder(r.Body).Decode(&event)
	db.Create(&event)
	json.NewEncoder(w).Encode(&event)
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var event Event
	db.Find(&event, params["id"])
	json.NewEncoder(w).Encode(&event)
}