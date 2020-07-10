package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
)

var Db *sql.DB

type DatabaseSettings struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func InitDb(connectionSettings DatabaseSettings) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		connectionSettings.Host,
		connectionSettings.Port,
		connectionSettings.User,
		connectionSettings.Password,
		connectionSettings.DbName,
	)
	var err error
	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic(err)
	}

	initDb()

	if err = Db.Ping(); err != nil {
		log.Panic(err)
	}
}
func initDb() {
	schemaRow := Db.QueryRow(
		`select exists(
				SELECT schema_name
				FROM information_schema.schemata
				WHERE schema_name = 'poker')`)

	var isSchemaExists bool
	err := schemaRow.Scan(&isSchemaExists)
	if err != nil {
		panic(err)
	}

	if isSchemaExists {
		return
	}

	pathWd, _ := os.Getwd()
	absolutePath := pathWd + "\\models\\sql\\DbInit.sql"
	query, _ := ioutil.ReadFile(absolutePath)

	_, err = Db.Exec(string(query))
	if err != nil {
		panic(err)
	}
}
