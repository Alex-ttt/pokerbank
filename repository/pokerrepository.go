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
