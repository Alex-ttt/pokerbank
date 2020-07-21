package repository

import (
	"database/sql"
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
	IsOdd          bool
}

type GameInfoViewModel struct {
	Games []Game
}

type IndexPageViewModel struct {
	Games      GameInfoViewModel
	Debts      PlayersDebtsViewModel
	Payments   PaymentsViewModel
	Offsetting OffsettingViewModel
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
		Games:      GetGamesInfo(db),
		Debts:      *GetPlayersDebts(db),
		Payments:   GetPlayersPayments(db),
		Offsetting: GetOffsetting(db),
	}
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
	isOdd := true
	for i := 0; i < len(result.Games); i++ {
		result.Games[i].IsOdd = isOdd
		isOdd = !isOdd
	}

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
