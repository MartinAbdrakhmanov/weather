package main

import (
	"log"
	"net/http"
	"weather/internal/api"
	"weather/internal/db"

	// _ "weather/docs"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found")
	}
	router := mux.NewRouter()

	router.Use(logRequest)
	fs := http.FileServer(http.Dir("templates/css"))
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", fs))
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/report", api.GetReport)
	router.HandleFunc("/form", api.GetForm)
	router.HandleFunc("/submit", api.SubmitForm).Methods(http.MethodPost)
	router.HandleFunc("/suggestion", db.GetSuggestion)
	router.HandleFunc("/", api.GetCity)
	router.HandleFunc("/{city}", api.GetWeather)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
