package main

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"encoding/json"
	"net/url"
	"time"
)

type TempEvent struct {
	Name string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
}

type Event struct {
	gorm.Model
	Name string `gorm:"not null" json:"name"`
	StartDate  time.Time `json:"startDate"`
	EndDate time.Time `json:"endDate"`
}

func (e *TempEvent) validate() url.Values {
	errs := url.Values{}

	if e.Name == "" {
		errs.Add("name", "Event name is required.")
	}

	if len(e.Name) < 3 {
		errs.Add("name", "Event names must be at least 3 characters long.")
	}

	if e.StartDate == "" {
		errs.Add("startDate", "A start date is required.")
	}

	if e.EndDate == "" {
		errs.Add("endDate", "An end date is required.")
	}

	return errs
}

func (e *TempEvent) parseDates(errs url.Values) (time.Time, time.Time, url.Values) {
	const timeFormat = "2006-01-02"

	sD, err := time.Parse(timeFormat, e.StartDate)
	if err != nil {
		errs.Add("startDate", "Date format should be: YYYY-MM-DD")
	}

	eD, err := time.Parse(timeFormat, e.EndDate)
	if err != nil {
		errs.Add("endDate", "Date format should be: YYYY-MM-DD")
	}

	return sD, eD, errs
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
	var tempEvent TempEvent

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tempEvent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	validationErrs := tempEvent.validate()
	encoder := json.NewEncoder(w)

	sD, eD, validationErrs := tempEvent.parseDates(validationErrs)
	if len(validationErrs) > 0 {
		err := map[string]interface{}{"ValidationError": validationErrs}
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(err)
		return
	}

	event := Event{Name:tempEvent.Name, StartDate:sD, EndDate:eD}

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