package RestAPI

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	host 		= "0.0.0.0"
	port 		= 5432
	username	= "iqbal"
	password	= "root"
	dbname		= "daily_diary_dev"
	sslmode		= "disable"
)

var db *sql.DB

type MyDatabase interface {
	InitDB() *sql.DB
}

func InitDB() *sql.DB {
	fmt.Println(username +" - "+ password +" - "+ dbname)
	var connectionString = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, username, password, dbname, sslmode)

	var err error
	db, err = sql.Open("postgres", connectionString)

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	return db
}
