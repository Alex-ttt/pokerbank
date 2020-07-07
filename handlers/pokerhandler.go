package handlers

import (
	"html/template"
	"net/http"
	"pokerscore/models"
	"pokerscore/repository"
)

type playersViewModel struct {
	Players []models.Player
}

func PlayersPage(writer http.ResponseWriter, _ *http.Request) {
	var playersVM playersViewModel = playersViewModel{
		Players: repository.GetAllPlayers(models.Db),
	}

	templates := template.Must(template.ParseFiles("templates/players-template.html"))
	if err := templates.ExecuteTemplate(writer, "players-template.html", playersVM); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
