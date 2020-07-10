package repository

import (
	"database/sql"
	"log"
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

func GetDebtsData(db *sql.DB) *PlayersDebtsViewModel {
	rows, err := db.Query("select * from poker.playersdebts();")
	if err != nil {
		log.Panic(err)
	}

	var result PlayersDebtsViewModel = PlayersDebtsViewModel{
		Losers:  make([]string, 0, 8),
		Winners: make([]Winner, 0, 8),
	}
	var previousWinner, winner, loser string
	var pWin, commonWin int
	for rows.Next() {
		_ = rows.Scan(&winner, &loser, &pWin, &commonWin)

		if winner != previousWinner {
			result.Winners = append(result.Winners, Winner{
				Name:      winner,
				Wins:      make([]int, 0, 8),
				CommonWin: commonWin,
			})
		}

		indexOfCurrentWinner := len(result.Winners) - 1
		result.Winners[indexOfCurrentWinner].Wins = append(result.Winners[indexOfCurrentWinner].Wins, pWin)

		var losersSlice = &(result.Losers)
		result.Losers = *(addElementIfItNotContained(losersSlice, loser))

		previousWinner = winner
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
