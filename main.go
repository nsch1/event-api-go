package main

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Event struct {
	gorm.Model
	Name string
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=secret sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/events", CreateEvent).Methods("POST")

	log.Fatal(http.ListenAndServe(":4008", router))
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {

}