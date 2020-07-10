package main

import (
	"net/http"
	"pokerscore/handlers"
	"pokerscore/models"
)

func main() {
	connectionSettings := models.DatabaseSettings{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "admin",
		DbName:   "pokerdb",
	}
	models.InitDb(connectionSettings)

	http.Handle("/static/", //final url can be anything
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handlers.IndexPage)
	_ = http.ListenAndServe(":8080", nil)
}
