package handlers

import (
	"html/template"
	"net/http"
	"pokerscore/models"
	"pokerscore/repository"
)

func IndexPage(writer http.ResponseWriter, _ *http.Request) {
	indexViewModel := repository.GetIndexPageViewModel(models.Db)

	templates := template.Must(template.ParseFiles("templates/index.html"))
	if err := templates.ExecuteTemplate(writer, "index.html", indexViewModel); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
