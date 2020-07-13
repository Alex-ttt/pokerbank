package handlers

import (
	"html/template"
	"log"
	"net/http"
	"pokerscore/models"
	"pokerscore/repository"
	"strconv"
)

type PlayerGameResult struct {
	WinnerId int
	LoserId  int
	Amount   int
}

func createPlayerGameResult(pgResult *[]PlayerGameResult) string {
	var result = "ARRAY["
	for i, el := range *pgResult {
		if i > 0 {
			result = result + ", "
		}
		result = result + "(" + strconv.Itoa(el.WinnerId) + ", " + strconv.Itoa(el.LoserId) + ", " + strconv.Itoa(el.Amount) + ")"
	}
	result = result + "]"
	return result
}

func IndexPage(writer http.ResponseWriter, _ *http.Request) {
	indexViewModel := repository.GetIndexPageViewModel(models.Db)

	gameResults := make([]PlayerGameResult, 0, 2)
	gameResults = append(gameResults, PlayerGameResult{
		WinnerId: 4,
		LoserId:  1,
		Amount:   111,
	}, PlayerGameResult{
		WinnerId: 2,
		LoserId:  1,
		Amount:   333,
	})
	a := createPlayerGameResult(&gameResults)
	_, _err := models.Db.Exec(
		"select poker.insertgameresult('Из кода динамически', now()::date," + a + "::poker.playergameresult[])")

	if _err != nil {
		log.Panic(_err)
	}
	templates := template.Must(template.ParseFiles("templates/index.html"))
	if err := templates.ExecuteTemplate(writer, "index.html", indexViewModel); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
