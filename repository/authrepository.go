package repository

import (
	"database/sql"
	"log"
)

func CheckLoginExists(db *sql.DB, login string) bool {
	row := db.QueryRow("select poker.doesloginexist($1);", login)

	var isLoginExists bool
	err := row.Scan(&isLoginExists)
	if err != nil {
		log.Panic(err)
	}

	return isLoginExists
}

func SetPassword(db *sql.DB, login string, password string) {
	_, err := db.Exec("select poker.setPassword($1, $2);", login, password)
	if err != nil {
		log.Panic(err)
	}
}

func GetPassword(db *sql.DB, login string) string {
	row := db.QueryRow("select poker.getpassword($1);", login)

	var password string
	err := row.Scan(&password)
	if err != nil {
		log.Panic(err)
	}

	return password
}
