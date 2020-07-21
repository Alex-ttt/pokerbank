package handlers

import (
	"database/sql/driver"
	"fmt"
	"html/template"
	"net/http"
	"pokerscore/models"
	"pokerscore/repository"
)

type PlayerGameResult struct {
	WinnerId int
	LoserId  int
	Amount   int
}

func (i *PlayerGameResult) Value() (driver.Value, error) {
	return fmt.Sprintf("(%v, %v, %v)", i.WinnerId, i.LoserId, i.Amount), nil
}

func IndexPage(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		writer.WriteHeader(http.StatusOK)

		return
	}

	indexViewModel := repository.GetIndexPageViewModel(models.Db)
	//gameResults := make([]*PlayerGameResult, 2, 2)
	//gameResults[0] = new(PlayerGameResult)
	//gameResults[0].WinnerId = 4
	//gameResults[0].LoserId = 1
	//gameResults[0].Amount = 777
	//gameResults[1] = new(PlayerGameResult)
	//gameResults[1].WinnerId = 2
	//gameResults[1].LoserId = 1
	//gameResults[1].Amount = 444

	//_, _err := models.Db.Exec(
	//	"select poker.insertgameresult('Из кода динамически', now()::date, $1::poker.playergameresult[])",
	//	pq.Array(gameResults))

	//if _err != nil {
	//	log.Panic(_err)
	//}
	addFunc := template.FuncMap{"add": func(x, y int) int {
		return x + y
	}}

	templates := template.Must(template.New("index.html").Funcs(addFunc).ParseFiles("templates/index.html"))
	if err := templates.ExecuteTemplate(writer, "index.html", indexViewModel); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
