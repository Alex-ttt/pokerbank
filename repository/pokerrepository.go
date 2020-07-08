package repository

import (
	"database/sql"
	. "pokerscore/models"
)

func GetAllPlayers(db *sql.DB) []Player {
	rows, err := db.Query("SELECT Id, Name, Surname FROM players")
	if err != nil {
		panic(err)
	}
	players := make([]Player, 0, 4)
	for rows.Next() {
		var (
			id            int
			name, surname string
		)
		err = rows.Scan(&id, &name, &surname)
		if err != nil {
			panic(err)
		}

		player := Player{
			Id:      id,
			Name:    name,
			Surname: surname,
		}

		players = append(players, player)
	}

	return players
}

func GetPlayersDebts(db *sql.DB) []PlayerWins {
	query := "select " +
		"\n win.Name as WinnerName, lose.Name as LoserName, " +
		"\n sum(d.amount) - ifnull(sum(dp.amount), 0) - ifnull(sum(invert_d.amount), 0) as CommonDebt" +
		"\n from debts as d" +
		"\n left join debts as invert_d ON invert_d.winnerId = d.loserId and invert_d.loserId = d.winnerId" +
		"\n left join debtPayments as dp on dp.fromPlayerId = d.loserId and dp.toPlayerId = d.winnerId" +
		"\n inner join players as win on win.Id = d.winnerId" +
		"\n inner join players as lose on lose.Id = d.loserId" +
		"\n group by win.Name, lose.Name" +
		"\n having CommonDebt > 0 " +
		"\n order by WinnerName, LoserName"

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	playersWins := make([]PlayerWins, 0)

	var currentWinner, currentLoser, previousWinner string
	var currentAmount int

	for rows.Next() {
		err = rows.Scan(&currentWinner, &currentLoser, &currentAmount)
		if err != nil {
			panic(err)
		}

		if previousWinner != currentWinner {
			playersWins = append(playersWins, PlayerWins{
				PlayerName: currentWinner,
				Win:        make(map[string]int, 0),
				Sum:        0,
			})
		}
		actualWinnerIndex := len(playersWins) - 1
		playersWins[actualWinnerIndex].Sum += currentAmount
		playersWins[actualWinnerIndex].Win[currentLoser] = currentAmount

		previousWinner = currentWinner
	}

	return playersWins

}
