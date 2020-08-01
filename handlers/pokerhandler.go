package handlers

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"pokerscore/models"
	"pokerscore/repository"
	"strconv"
	"time"
)

type PlayerGameResult struct {
	WinnerId int `json:"winnerId"`
	LoserId  int `json:"looserId"`
	Amount   int `json:"amount"`
}

func (i *PlayerGameResult) Value() (driver.Value, error) {
	return fmt.Sprintf("(%v, %v, %v)", i.WinnerId, i.LoserId, i.Amount), nil
}

type GameResult struct {
	GameName    string
	GameDate    time.Time
	GameResults []PlayerGameResult
}

//
func (gameResult *GameResult) UnmarshalJSON(data []byte) error {
	obj := &map[string]string{}
	if err := json.Unmarshal(data, &obj); err != nil {
		fmt.Println(err)
		return err
	}

	gameResult.GameName = (*obj)["gameName"]
	gameDate, err := time.Parse("02.01.2006", (*obj)["gameDate"])
	if err != nil {
		fmt.Println(err)
		return err
	}
	gameResult.GameDate = gameDate
	results := (*obj)["gameResults"]
	var playersResults []map[string]string
	if err := json.Unmarshal([]byte(results), &playersResults); err != nil {
		fmt.Println(err)
		return err
	}

	gameResult.GameResults = make([]PlayerGameResult, 0, len(playersResults))
	for _, el := range playersResults {
		winnerId, _ := strconv.Atoi(el["winnerId"])
		looserId, _ := strconv.Atoi(el["looserId"])
		amount, _ := strconv.Atoi(el["amount"])

		gameResult.GameResults = append(gameResult.GameResults, PlayerGameResult{
			WinnerId: winnerId,
			LoserId:  looserId,
			Amount:   amount,
		})
	}

	return nil
}

func IndexPage(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {

		decoder := json.NewDecoder(request.Body)
		var t GameResult
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}

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
