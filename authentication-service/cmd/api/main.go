package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64 = 0

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service...")

	conn := connectToDB()
	if conn != nil {
		logSnippet := "[authentication-service][main.go] =>"
		log.Panic(logSnippet + "Can't connect to PostgreSQL DB")
	}

	// Set up configuration...
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	logSnippet := "[authentication-service][main.go][connectToDB()] =>"

	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)

		if err != nil {
			log.Printf("%s openDB() returned an error (incrementing count):\n", logSnippet)
			log.Println(err)
			counts++
		} else {
			log.Printf("%s openDB() did not return an error:\n", logSnippet)
			return connection
		}

		if counts > 10 {
			log.Printf("%s Error count has exceeded 10, could not connec to to database\n", logSnippet)
			return nil
		}

		log.Printf("%s Sleeping for 2 seconds...\n", logSnippet)
		time.Sleep(2 * time.Second)
		continue
	}

}
