package repository

import (
	"database/sql"
	"log"
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
	Games GameInfoViewModel
	Debts PlayersDebtsViewModel
}

func GetIndexPageViewModel(db *sql.DB) IndexPageViewModel {
	return IndexPageViewModel{
		Games: GetGamesInfo(db),
		Debts: *GetPlayersDebts(db),
	}
}

func GetGamesInfo(db *sql.DB) GameInfoViewModel {
	rows, err := db.Query("select * from poker.gamesinfo();")
	if err != nil {
		log.Panic(err)
	}

	result := GameInfoViewModel{
		Games: make([]Game, 0, 8),
	}
	var (
		gameName, playerName                       string
		gameDate                                   sql.NullTime
		amount, gameAmount, gameId, previousGameId int
		isOddRow                                   = true
	)
	for rows.Next() {
		err = rows.Scan(&gameId, &gameName, &playerName, &amount, &gameDate, &gameAmount)
		if err != nil {
			log.Panic(err)
		}

		if gameId != previousGameId {
			result.Games = append(result.Games, Game{
				Name:           gameName,
				Date:           gameDate.Time,
				IsDateValid:    gameDate.Valid,
				PlayersResults: make([]PlayerGameResult, 0, 8),
				IsOdd:          isOddRow,
				Amount:         gameAmount,
			})

			isOddRow = !isOddRow
		}

		indexOfCurrentGame := len(result.Games) - 1
		result.Games[indexOfCurrentGame].PlayersResults =
			append(result.Games[indexOfCurrentGame].PlayersResults, PlayerGameResult{
				Name:   playerName,
				Amount: amount,
			})

		previousGameId = gameId
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
