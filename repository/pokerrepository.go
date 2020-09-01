package repository

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"log"
	"sort"
	"time"
)

type Winner struct {
	Name      string
	Wins      []int
	CommonWin int
}

type PlayersDebtsViewModel struct {
	Losers  []string
	Winners []Winner
}

type GameResultInsertDto struct {
	GameName    string
	GameDate    time.Time
	GameResults []PlayerGameResultInsertDto
}

type PlayerGameResultInsertDto struct {
	WinnerId int `json:"winnerId"`
	LoserId  int `json:"looserId"`
	Amount   int `json:"amount"`
}

func (gameResult *GameResultInsertDto) UnmarshalJSON(data []byte) error {
	obj := &map[string]string{}
	if err := json.Unmarshal(data, &obj); err != nil {
		fmt.Println(err)
		return err
	}

	gameResult.GameName = (*obj)["gameName"]
	gameDate, err := time.Parse("2006-01-02", (*obj)["gameDate"])
	if err != nil {
		fmt.Println(err)
		return err
	}
	gameResult.GameDate = gameDate
	results := (*obj)["gameResults"]
	var playersResults []map[string]int
	if err := json.Unmarshal([]byte(results), &playersResults); err != nil {
		fmt.Println(err)
		return err
	}

	gameResult.GameResults = make([]PlayerGameResultInsertDto, 0, len(playersResults))
	for _, el := range playersResults {
		gameResult.GameResults = append(gameResult.GameResults, PlayerGameResultInsertDto{
			WinnerId: el["winnerId"],
			LoserId:  el["looserId"],
			Amount:   el["amount"],
		})
	}

	return nil
}

func (i PlayerGameResultInsertDto) Value() (driver.Value, error) {
	return fmt.Sprintf("(%v, %v, %v)", i.WinnerId, i.LoserId, i.Amount), nil
}

type PlayerGameResult struct {
	Name   string
	Amount int
}

type Game struct {
	Name           string
	Date           time.Time
	IsDateValid    bool
	Amount         int
	PlayersResults []PlayerGameResult
}

type GameInfoViewModel struct {
	Games []Game
}

type PlayersSourceViewModel struct {
	Players []Player
}

type Player struct {
	Id   int
	Name string
}

type IndexPageViewModel struct {
	Games         GameInfoViewModel
	Debts         PlayersDebtsViewModel
	Payments      PaymentsViewModel
	Offsetting    OffsettingViewModel
	PlayersSource PlayersSourceViewModel
}

type Payment struct {
	Payer     string
	Recipient string
	Amount    int
}

type PaymentsViewModel struct {
	Payments []Payment
}

type Offsetting struct {
	Recipient string
	OldDebtor string
	NewDebtor string
	Amount    int
}

type OffsettingViewModel struct {
	Offsets []Offsetting
}

func GetIndexPageViewModel(db *sql.DB) IndexPageViewModel {
	return IndexPageViewModel{
		Games:         GetGamesInfo(db),
		Debts:         *GetPlayersDebts(db),
		Payments:      GetPlayersPayments(db),
		Offsetting:    GetOffsetting(db),
		PlayersSource: GetPlayersSource(db),
	}
}

func GetPlayersSource(db *sql.DB) PlayersSourceViewModel {
	rows, err := db.Query("select * from poker.playerslist();")
	if err != nil {
		log.Panic(err)
	}

	var (
		id     int
		name   string
		result = PlayersSourceViewModel{
			Players: make([]Player, 0, 8),
		}
	)

	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Panic(err)
		}

		result.Players = append(result.Players, Player{
			Id:   id,
			Name: name,
		})
	}

	return result
}

func GetOffsetting(db *sql.DB) OffsettingViewModel {
	rows, err := db.Query("select * from poker.playersoffsetting();")
	if err != nil {
		log.Panic(err)
	}

	var (
		recipient, oldDebtor, newDebtor string
		amount                          int
	)
	result := OffsettingViewModel{Offsets: make([]Offsetting, 0, 8)}

	for rows.Next() {
		err = rows.Scan(&recipient, &oldDebtor, &newDebtor, &amount)
		if err != nil {
			log.Panic(err)
		}

		result.Offsets = append(result.Offsets, Offsetting{
			Recipient: recipient,
			OldDebtor: oldDebtor,
			NewDebtor: newDebtor,
			Amount:    amount,
		})
	}

	return result
}

func GetPlayersPayments(db *sql.DB) PaymentsViewModel {
	rows, err := db.Query("select * from poker.playerspayments();")
	if err != nil {
		log.Panic(err)
	}

	var (
		payer, recipient string
		amount           int
	)
	result := PaymentsViewModel{Payments: make([]Payment, 0, 8)}
	for rows.Next() {
		err = rows.Scan(&payer, &recipient, &amount)
		if err != nil {
			log.Panic(err)
		}

		result.Payments = append(result.Payments, Payment{
			Payer:     payer,
			Recipient: recipient,
			Amount:    amount,
		})

	}

	return result
}

func AddGameWithResults(db *sql.DB, gameResults *GameResultInsertDto) error {
	_, _err := db.Exec(
		"select poker.insertgameresult($1, $2::date, $3::poker.playergameresult[])",
		gameResults.GameName,
		gameResults.GameDate.Format("2006-01-02"),
		pq.Array(gameResults.GameResults))

	if _err != nil {
		log.Panic(_err)
		return _err
	}

	return nil
}

func GetGamesInfo(db *sql.DB) GameInfoViewModel {
	rows, err := db.Query("select * from poker.gamesinfo();")
	if err != nil {
		log.Panic(err)
	}

	var (
		gameName, playerName       string
		gameDate                   sql.NullTime
		amount, gameAmount, gameId int
	)

	resultMap := make(map[int]*Game, 0)
	for rows.Next() {
		err = rows.Scan(&gameId, &gameName, &playerName, &amount, &gameDate, &gameAmount)
		if err != nil {
			log.Panic(err)
		}

		game, exists := resultMap[gameId]
		if !exists {
			game = new(Game)
			game.Name = gameName
			game.Date = gameDate.Time
			game.IsDateValid = gameDate.Valid
			game.PlayersResults = make([]PlayerGameResult, 0, 8)
			game.Amount = gameAmount

			resultMap[gameId] = game
		}

		game.PlayersResults =
			append(game.PlayersResults, PlayerGameResult{
				Name:   playerName,
				Amount: amount,
			})
	}

	result := GameInfoViewModel{
		Games: make([]Game, 0, len(resultMap)),
	}

	for _, value := range resultMap {
		result.Games = append(result.Games, *value)
	}

	sort.SliceStable(result.Games, func(i, j int) bool {
		return result.Games[i].Date.After(result.Games[j].Date)
	})

	return result
}

func GetPlayersDebts(db *sql.DB) *PlayersDebtsViewModel {
	rows, err := db.Query("select * from poker.playersdebts();")
	if err != nil {
		log.Panic(err)
	}

	const DefaultSliceCapacity = 8
	var (
		result = PlayersDebtsViewModel{
			Losers:  make([]string, 0, DefaultSliceCapacity),
			Winners: make([]Winner, 0, DefaultSliceCapacity),
		}
		winner, loser                               string
		pWin, commonWin, winnerId, previousWinnerId int
	)
	for rows.Next() {
		_ = rows.Scan(&winnerId, &winner, &loser, &pWin, &commonWin)

		if winnerId != previousWinnerId {
			result.Winners = append(result.Winners, Winner{
				Name:      winner,
				Wins:      make([]int, 0, DefaultSliceCapacity),
				CommonWin: commonWin,
			})
		}

		indexOfCurrentWinner := len(result.Winners) - 1
		result.Winners[indexOfCurrentWinner].Wins = append(result.Winners[indexOfCurrentWinner].Wins, pWin)

		var losersSlice = &(result.Losers)
		result.Losers = *(addElementIfItNotContained(losersSlice, loser))

		previousWinnerId = winnerId
	}

	sort.SliceStable(result.Winners, func(i, j int) bool {
		return result.Winners[i].Name < result.Winners[j].Name
	})

	return &result
}

func addElementIfItNotContained(_array *[]string, element string) *[]string {
	for _, currentEl := range *_array {
		if currentEl == element {
			return _array
		}
	}

	result := append(*_array, element)
	return &result
}
