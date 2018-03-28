package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/coffee-shop/db"
	"github.com/coffee-shop/handlers"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatal("DB initialization failed")
		return
	}

	http.Handle("/", router())
	log.Println("Server listening to port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func router() *mux.Router {
	r := mux.NewRouter()

	r.Path("/drink").
		Methods("POST").
		HandlerFunc(handlers.DrinkCreateHandler)

	r.Path("/drink").
		Methods("DELETE").
		HandlerFunc(handlers.DrinkDeleteHandler)

	r.Path("/drinks").
		Methods("GET").
		HandlerFunc(handlers.DrinkSearchHandler)

	return r
}
