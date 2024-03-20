package main

import (
	"database/sql"
	"filmoteka/internal/delivery/http/routes"
	"filmoteka/internal/domain/usecase"
	"filmoteka/internal/domain/usecase/actormovieusecase"
	"filmoteka/internal/domain/usecase/actorusecase"
	"filmoteka/internal/domain/usecase/movieusecase"
	"filmoteka/internal/domain/usecase/userusecase"
	"filmoteka/internal/storage/actormoviestorage"
	"filmoteka/internal/storage/actorstorage"
	"filmoteka/internal/storage/moviestorage"
	"filmoteka/internal/storage/userstorage"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
	"time"
)

var count uint8

func main() {

	conn := connectToDB()
	defer conn.Close()

	sessionManager := newSessionManager()

	userStorage := userstorage.New(conn)
	movieStorage := moviestorage.New(conn)
	actorStorage := actorstorage.New(conn)
	actormovieStorage := actormoviestorage.New(conn)

	userUseCase := userusecase.New(userStorage)
	movieUseCase := movieusecase.New(movieStorage)
	actorUseCase := actorusecase.New(actorStorage)
	actormovieUseCase := actormovieusecase.New(actormovieStorage)

	uc := usecase.UseCase{
		UserUseCase:       userUseCase,
		MovieUseCase:      movieUseCase,
		ActorUseCase:      actorUseCase,
		ActorMovieUseCase: actormovieUseCase,
	}

	r := routes.Routes(&uc, sessionManager)

	srv := http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: r,
	}

	log.Println("Starting server on port: ", srv.Addr)

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

func newSessionManager() *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = true
	sessionManager.Store = memstore.New()
	return sessionManager
}
