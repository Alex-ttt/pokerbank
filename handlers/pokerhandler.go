package handlers

import (
	"html/template"
	"net/http"
	"pokerscore/models"
	"pokerscore/repository"
)

func IndexPage(writer http.ResponseWriter, _ *http.Request) {
	playersDebtsViewModel := *repository.GetDebtsData(models.Db)

	templates := template.Must(template.ParseFiles("templates/index.html"))
	if err := templates.ExecuteTemplate(writer, "index.html", playersDebtsViewModel); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
