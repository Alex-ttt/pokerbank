package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var Db *sql.DB

func InitDb(dataSourceName string){
	var err error
	Db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = Db.Ping(); err != nil {
		log.Panic(err)
	}
}
