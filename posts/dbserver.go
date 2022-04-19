package main

import (
	"database/sql"
	"log"
)

const datetimeLayer = "2006-01-02 15:04:05.999999"

func getDbConnectSource() string {
	return "host=db user=" + goDotEnvVariable("DB_USER") +
		" password=" + goDotEnvVariable("DB_PW") +
		" dbname=" + goDotEnvVariable("DB_NAME") +
		" sslmode=disable"
}

func getDbConnection() (db *sql.DB) {
	db, err := sql.Open("postgres", getDbConnectSource())
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	return
}
