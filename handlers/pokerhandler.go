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
	_ = repository.GetPlayersDebts(models.Db)

	templates := template.Must(template.ParseFiles("templates/index.html"))
	if err := templates.ExecuteTemplate(writer, "index.html", playersVM); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
