package main

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/OnYyon/gRPCCalculator/web/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	apiHandler := handlers.NewAPIHandler("http://0.0.0.0:" + "8080")

	r.HandleFunc("/", apiHandler.HomeHandler).Methods("GET")
	r.HandleFunc("/login", apiHandler.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/register", apiHandler.RegisterHandler).Methods("GET", "POST")
	r.HandleFunc("/expressions", apiHandler.ExpressionsHandler).Methods("GET")
	r.HandleFunc("/calculate", apiHandler.CalculateHandler).Methods("GET", "POST")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join("web", "static")))))

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8081",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Starting web server on localhost:8081")
	log.Fatal(srv.ListenAndServe())
}
