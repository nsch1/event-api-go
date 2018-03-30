package main

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
)

type Event struct {
	gorm.Model
	Name string `gorm:"not null" validate:"min=3"`
}

var (
	db *gorm.DB
	err error
	validate *validator.Validate
)

func main() {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=secret sslmode=disable")
	if err != nil {
		panic("oops")
	}
	defer db.Close()

	db.AutoMigrate(&Event{})

	validate = validator.New()

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

	err := validate.Struct(event)
	if err != nil {
		http.Error(w, "Name must be at least 3 characters.", http.StatusBadRequest)
		return
	}

	db.Create(&event)
	json.NewEncoder(w).Encode(&event)
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var event Event
	db.Find(&event, params["id"])
	if event.ID == 0 {
		http.Error(w, "No event found.", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(&event)
}