package main

import (
	"database/sql"
	"filmoteka/domain/models"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
	"time"
)

var count uint8

type Config struct {
	Storage *sql.DB
	Models  models.Models
}

func main() {

	conn := connectToDB()
	defer conn.Close()

	app := Config{
		Storage: conn,
		Models:  models.New(conn),
	}

	srv := http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Println("Error starting server: ", err)
	}
}

func connectToDB() *sql.DB {
	// connect to postgres
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready...")
			count++
		} else {
			log.Println("Connected to database!")
			return connection
		}

		if count > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for three seconds...")
		time.Sleep(3 * time.Second)
		continue
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
