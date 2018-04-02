package main

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"encoding/json"
	"net/url"
)

type Event struct {
	gorm.Model
	Name string `gorm:"not null" json:"name"`
}

func (e *Event) validate() url.Values {
	errs := url.Values{}

	if e.Name == "" {
		errs.Add("name", "Event name is required.")
	}

	if len(e.Name) < 3 {
		errs.Add("name", "Event names must be at least 3 characters long.")
	}

	return errs
}

var (
	db *gorm.DB
	err error
)

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

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	encoder := json.NewEncoder(w)
	if validationErrs := event.validate(); len(validationErrs) > 0 {
		err := map[string]interface{}{"ValidationError": validationErrs}
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(err)
	}

	db.Create(&event)
	encoder.Encode(&event)
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