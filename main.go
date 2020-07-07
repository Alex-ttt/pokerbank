package main

import (
	"net/http"
	"pokerscore/handlers"
	"pokerscore/models"
)

func main() {
	models.InitDb("./poker.sqlite")

	http.Handle("/static/", //final url can be anything
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	http.HandleFunc("/players", handlers.PlayersPage)
	_ = http.ListenAndServe(":8080", nil)

}
