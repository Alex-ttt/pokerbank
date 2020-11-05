package handlers

import (
	"encoding/json"
	"github.com/Alex-ttt/pokerbank/models"
	"github.com/Alex-ttt/pokerbank/repository"
	"html/template"
	"net/http"
)

func AddDebtPayment(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		decoder := json.NewDecoder(request.Body)
		var insertDto repository.DebtPaymentInsertDto

		err := decoder.Decode(&insertDto)
		if err != nil {
			panic(err)
		}

		err = repository.AddDebtPayment(models.Db, &insertDto)
		if err != nil {
			panic(err)
		}

		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte("{}"))
		return
	}
}

func IndexPage(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		decoder := json.NewDecoder(request.Body)
		var gameResult repository.GameResultInsertDto
		var err error

		if err = decoder.Decode(&gameResult); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = repository.AddGameWithResults(models.Db, &gameResult); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte("{}"))
		return
	}

	indexViewModel := repository.GetIndexPageViewModel(models.Db)
	//indexViewModel := repository.GetMockPageViewModel()
	addFunc := template.FuncMap{"add": func(x, y int) int {
		return x + y
	}}
	seqFunc := template.FuncMap{"seq": func(n int) []int {
		return make([]int, n, n)
	}}

	templates := template.Must(
		template.
			New("index.html").
			Funcs(addFunc).
			Funcs(seqFunc).
			ParseFiles("templates/index.html"))

	if err := templates.ExecuteTemplate(writer, "index.html", indexViewModel); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

}
