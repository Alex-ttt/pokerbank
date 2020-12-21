package handlers

import (
	"encoding/json"
	"github.com/Alex-ttt/pokerbank/repository"
	"github.com/Alex-ttt/pokerbank/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AddDebtPayment(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var insertDto repository.DebtPaymentInsertDto

	err := decoder.Decode(&insertDto)
	if err != nil {
		panic(err)
	}

	err = repository.AddDebtPayment(services.Db, &insertDto)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"Content-Type": "application/json",
	})
}

func AddGameResult(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var gameResult repository.GameResultInsertDto
	var err error

	if err = decoder.Decode(&gameResult); err != nil {
		log.Panic(http.StatusInternalServerError)
		return
	}

	if err = repository.AddGameWithResults(services.Db, &gameResult); err != nil {
		log.Panic(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Content-Type": "application/json",
	})
}

func IndexPage(c *gin.Context) {

	indexViewModel := repository.GetIndexPageViewModel(services.Db)
	//indexViewModel := repository.GetMockPageViewModel()

	c.HTML(http.StatusOK, "index.html", indexViewModel)

	//templates := template.Must(
	//	template.
	//		New("index.html").
	//		Funcs(addFunc).
	//		Funcs(seqFunc).
	//		ParseFiles("templates/index.html"))
	//
	//if err := templates.ExecuteTemplate(writer, "index.html", indexViewModel); err != nil {
	//	http.Error(writer, err.Error(), http.StatusInternalServerError)
	//	return
	//}
}
