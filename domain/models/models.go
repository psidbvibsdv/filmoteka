package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var ErrNoRecord = errors.New("no entries found")

type Date struct {
	sql.NullTime
}

func (dt *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.Time.Format("2006-01-02"))
}

func (dt *Date) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "null" {
		dt.Valid = false
		return
	}
	dt.Valid = true
	dt.Time, err = time.Parse("2006-01-02", s)
	return
}

type Models struct {
	Actor          Actor
	Movie          Movie
	MovieWithActor MovieWithActor
	ActorMovies    ActorMovies
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Actor struct {
	ActorID     int    `json:"actorid,omitempty"`
	Name        string `json:"name,omitempty"`
	Gender      string `json:"gender,omitempty"`
	DateOfBirth Date   `json:"dateofbirth,omitempty"`
}

type MovieWithActor struct {
	MovieID     int     `json:"movieid,omitempty"`
	Title       string  `json:"Title"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	ReleaseDate Date    `json:"releasedate"`
	ActorName   string  `json:"actorname"`
}

type ActorMovies struct {
	ActorId int
	Name    string
	Movies  []*Movie
}

type Movie struct {
	MovieID     int     `json:"movieid,omitempty"`
	Title       string  `json:"Title"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	ReleaseDate Date    `json:"releasedate"`
}

var db *sql.DB
var dbTimeout = 5 * time.Second

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Actor:          Actor{},
		Movie:          Movie{},
		MovieWithActor: MovieWithActor{},
		ActorMovies:    ActorMovies{},
	}
}

// Methods for User
func (u *User) auth() {
	bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(u.Password))
}

func (u *User) GetByEmail() (password string, role string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT password, role FROM users WHERE email = $1`
	_, err = db.QueryContext(ctx, query, u.Email)
	if err != nil {
		log.Println("Error getting user by email from the table", err)
		return "", "", err
	}

	err = db.QueryRowContext(ctx, query, u.Email).Scan(&password, &role)
	if err != nil {
		log.Println("Error getting user by email from the table", err)
		return "", "", err
	}

	return password, role, err
}

//Methods for Actor

func (a *Actor) GetAll() ([]*Actor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT actorid, name, gender, dateofbirth
	FROM actors ORDER BY name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error getting all actors from the table", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows", err)
		}
	}(rows)

	var actors []*Actor

	for rows.Next() {
		var actor Actor
		err = rows.Scan(
			&actor.ActorID,
			&actor.Name,
			&actor.Gender,
			&actor.DateOfBirth,
		)
		if err != nil {
			log.Println("Error scanning actor rows", err)
			return nil, err
		}

		actors = append(actors, &actor)
	}
	if len(actors) < 1 {
		log.Println("No actors found in the table")
		return nil, ErrNoRecord
	}

	return actors, nil

}

func (a *Actor) Create() (*Actor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `INSERT INTO actors (name, gender, dateofbirth)
	values ($1, $2, $3)`
	_, err := db.ExecContext(ctx, query, a.Name, a.Gender, a.DateOfBirth)
	if err != nil {
		log.Println("Error inserting actor into a table", err)
		encoder, _ := json.Marshal(a)
		log.Println("Actor: ", string(encoder))
		return nil, err
	}
	return a, nil
}

func (a *Actor) GetByID() (*Actor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT actorid, name, gender, dateofbirth FROM actors WHERE actorid = $1`
	rows, err := db.QueryContext(ctx, query, a.ActorID)
	if err != nil {
		log.Println("Error getting actor by id from the table", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows", err)
		}
	}(rows)

	actor := &Actor{}
	for rows.Next() {
		err = rows.Scan(
			&actor.ActorID,
			&actor.Name,
			&actor.Gender,
			&actor.DateOfBirth,
		)
		if err != nil {
			log.Println("Error scanning actor rows", err)
			return nil, err
		}
	}

	//defaultDate := time.Date(1000, 01, 01, 0, 0, 0, 0, time.UTC)
	if !actor.DateOfBirth.Valid {
		return nil, errors.New("no actor found")
	}
	log.Println("Actor: ", actor.ActorID, actor.DateOfBirth.Time, actor.DateOfBirth.Valid)
	return actor, nil
}

func (a *Actor) Update() (*Actor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	log.Println(a.DateOfBirth)
	formattedDOB, _ := json.Marshal(a.DateOfBirth.Time)
	log.Println("Formatted DOB: ", string(formattedDOB))
	query := `UPDATE actors SET name = CASE WHEN $1 = '' THEN name ELSE COALESCE($1) END, gender = CASE WHEN $2 = '' THEN gender ELSE COALESCE($2) END, dateofbirth = CASE WHEN $3 < '1000-1-1' THEN dateofbirth ELSE CAST(COALESCE($3) AS DATE) END WHERE actorid = ($4)`
	_, err := db.ExecContext(ctx, query, a.Name, a.Gender, formattedDOB, a.ActorID)
	if err != nil {
		log.Println("Error updating actor in the table", err)
		encoder, _ := json.Marshal(a)
		log.Println("Actor: ", string(encoder))
		return nil, err
	}
	res, err := a.GetByID()
	log.Println("Updated actor: ", res)
	if err != nil {
		log.Println("Error error returning updated actor from the table", err)
		return nil, err
	}
	return res, nil
}

func (a *Actor) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `DELETE FROM actors WHERE actorid = $1`
	_, err := db.ExecContext(ctx, query, a.ActorID)
	if err != nil {
		log.Println("Error deleting actor from the table", err)
		return err
	}
	return nil
}

//Methods for Movie

func (m *Movie) GetAll(sortParam string) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var query string
	switch sortParam {
	case "date":
		query = `SELECT movieid, title, description, rating, releasedate FROM movies ORDER BY releasedate`
	case "title":
		query = `SELECT movieid, title, description, rating, releasedate FROM movies ORDER BY title`
	case "rating":
		fallthrough
	case "":
		query = `SELECT movieid, title, description, rating, releasedate FROM movies ORDER BY rating DESC`
	default:
		log.Println("Invalid sort parameter")
		return nil, errors.New("invalid sort parameter")
	}
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error getting all movies from the table", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows", err)
		}
	}(rows)
	var movies []*Movie

	for rows.Next() {
		var movie Movie
		err = rows.Scan(
			&movie.MovieID,
			&movie.Title,
			&movie.Description,
			&movie.Rating,
			&movie.ReleaseDate,
		)
		if err != nil {
			log.Println("Error scanning movie rows", err)
			return nil, err
		}

		movies = append(movies, &movie)
	}
	if len(movies) < 1 {
		log.Println("No movies found in the table")
		return nil, ErrNoRecord
	}
	return movies, nil
}

func (m *Movie) GetByID() (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT movieid, title, description, rating, releasedate FROM movies WHERE movieid = $1`
	rows, err := db.QueryContext(ctx, query, m.MovieID)
	if err != nil {
		log.Println("Error getting movie by id from the table", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows", err)
		}
	}(rows)

	movie := &Movie{}
	for rows.Next() {
		err = rows.Scan(
			&movie.MovieID,
			&movie.Title,
			&movie.Description,
			&movie.Rating,
			&movie.ReleaseDate,
		)
		if err != nil {
			log.Println("Error scanning movie rows", err)
			return nil, err
		}
	}

	if !movie.ReleaseDate.Valid {
		return nil, errors.New("no movie found")
	}
	log.Println("Movie: ", movie)
	return movie, nil
}

func (m *Movie) Create() (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `INSERT INTO movies (title, description, rating, releasedate)
	values ($1, $2, $3, $4)`
	_, err := db.ExecContext(ctx, query, m.Title, m.Description, m.Rating, m.ReleaseDate)
	if err != nil {
		log.Println("Error inserting movie into a table", err)
		encoder, _ := json.Marshal(m)
		log.Println("Movie: ", string(encoder))
		return nil, err
	}
	return m, nil
}

func (m *Movie) Update() (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	formattedRD, _ := json.Marshal(m.ReleaseDate.Time)
	log.Println("Formatted DOB: ", string(formattedRD))
	query := `UPDATE movies SET title = CASE WHEN $1 = '' THEN title ELSE COALESCE($1) END, description = CASE WHEN $2 = '' THEN description ELSE COALESCE($2) END, rating = CASE WHEN $3 = 0 THEN rating ELSE COALESCE($3) END, releasedate = CASE WHEN $4 < '1000-1-1' THEN releasedate ELSE CAST(COALESCE($4) AS DATE) END WHERE movieid = ($5)`
	_, err := db.ExecContext(ctx, query, m.Title, m.Description, m.Rating, formattedRD, m.MovieID)
	if err != nil {
		log.Println("Error updating movie in the table", err)
		encoder, _ := json.Marshal(m)
		log.Println("Movie: ", string(encoder))
		return nil, err
	}
	res, err := m.GetByID()
	log.Println("Updated movie: ", res)
	if err != nil {
		log.Println("Error returning updated movie from the table ", err)
		return nil, err
	}
	return res, nil
}

func (m *Movie) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `DELETE FROM movies WHERE movieid = $1`
	_, err := db.ExecContext(ctx, query, m.MovieID)
	if err != nil {
		log.Println("Error deleting movie from the table", err)
		return err
	}
	return nil
}

func (m *Movie) GetByActorName(name string, surname string) ([]*MovieWithActor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT m.*, a.name AS actor_name
FROM Movies m
         JOIN actormovie am ON m.movieid = am.movieid
         JOIN Actors a ON am.actorid = a.actorid 
WHERE a.name ILIKE  $1 OR a.name ILIKE $2;`

	rows, err := db.QueryContext(ctx, query, "%"+name+"%", "%"+surname+"%")
	if err != nil {
		log.Println("Error getting movies by actor name from the table", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows", err)
		}
	}(rows)

	var movies []*MovieWithActor

	for rows.Next() {
		var movie MovieWithActor
		err = rows.Scan(
			&movie.MovieID,
			&movie.Title,
			&movie.Description,
			&movie.Rating,
			&movie.ReleaseDate,
			&movie.ActorName,
		)
		if err != nil {
			log.Println("Error scanning actor rows", err)
			return nil, err
		}

		movies = append(movies, &movie)
	}
	if len(movies) < 1 {
		log.Println("No movies found in the table")
		return nil, ErrNoRecord
	}

	return movies, nil

}

func (m *Movie) GetByMovieName(moviename string) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT movieid, title, description, rating, releasedate FROM movies WHERE title ILIKE $1 OR description ILIKE $2;`

	rows, err := db.QueryContext(ctx, query, "%"+moviename+"%", "%"+moviename+"%")
	if err != nil {
		log.Println("Error getting movies by movie name from the table", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows", err)
		}
	}(rows)

	var movies []*Movie

	for rows.Next() {
		var movie Movie
		err = rows.Scan(
			&movie.MovieID,
			&movie.Title,
			&movie.Description,
			&movie.Rating,
			&movie.ReleaseDate,
		)
		if err != nil {
			log.Println("Error scanning actor rows", err)
			return nil, err
		}

		movies = append(movies, &movie)
	}

	if len(movies) < 1 {
		log.Println("No movies found in the table")
		return nil, ErrNoRecord
	}
	return movies, nil
}

func (m *Movie) GetActorsForMovie() ([]*Actor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT a.ActorID, a.Name FROM Actors a JOIN actormovie am ON a.actorid = am.actorid WHERE am.movieid = $1`

	rows, err := db.QueryContext(ctx, query, m.MovieID)
	if err != nil {
		log.Println("Error getting actors for movie from the table", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows", err)
		}
	}(rows)

	var actors []*Actor
	for rows.Next() {
		var actor Actor
		err = rows.Scan(
			&actor.ActorID,
			&actor.Name,
		)
		if err != nil {
			log.Println("Error scanning actor rows", err)
			return nil, err
		}

		actors = append(actors, &actor)

	}
	return actors, nil
}

func (m *Movie) GetMoviesForActor(actorid int) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT m.movieid, m.title, m.description, m.rating, m.releasedate FROM Movies m JOIN actormovie am ON m.movieid = am.movieid WHERE am.actorid = $1`

	rows, err := db.QueryContext(ctx, query, actorid)
	if err != nil {
		log.Println("Error getting movies for actor from the table", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows", err)
		}
	}(rows)

	var movies []*Movie
	for rows.Next() {
		var movie Movie
		err = rows.Scan(
			&movie.MovieID,
			&movie.Title,
			&movie.Description,
			&movie.Rating,
			&movie.ReleaseDate,
		)
		if err != nil {
			log.Println("Error scanning movie rows", err)
			return nil, err
		}

		movies = append(movies, &movie)

	}
	return movies, nil

}

func (m *Movie) GetActorsAndMoviesForMovie() ([]*ActorMovies, error) {
	actors, err := m.GetActorsForMovie()
	if err != nil {
		log.Println("Error getting actors for movie", err)
		return nil, err
	}
	var result []*ActorMovies
	for _, actor := range actors {
		movies, err := m.GetMoviesForActor(actor.ActorID)
		if err != nil {
			log.Println("Error getting movies for actor", err)
			return nil, err
		}
		result = append(result, &ActorMovies{
			ActorId: actor.ActorID,
			Name:    actor.Name,
			Movies:  movies,
		})
	}
	return result, nil
}

func (m *Movie) AddActorToMovie(actorid int, movieid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	a := Actor{ActorID: actorid}
	m.MovieID = movieid
	_, err := a.GetByID()
	if err != nil {
		log.Println("Error getting actor by id from the table")
		return err
	}
	_, err = m.GetByID()
	if err != nil {
		log.Println("Error getting movie by id from the table")
		return err
	}
	query := `INSERT INTO actormovie (actorid, movieid) VALUES ($1, $2)`
	_, err = db.ExecContext(ctx, query, actorid, movieid)
	if err != nil {
		log.Println("Error adding actor to movie in the table", err)
		return err
	}
	return nil
}

func (m *Movie) DeleteActorFromMovie(actorid int, movieid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	a := Actor{ActorID: actorid}
	m.MovieID = movieid
	_, err := a.GetByID()
	if err != nil {
		log.Println("Error getting actor by id from the table")
		return err
	}
	_, err = m.GetByID()
	if err != nil {
		log.Println("Error getting movie by id from the table")
		return err
	}
	query := `DELETE FROM actormovie WHERE actorid = $1 AND movieid = $2`
	_, err = db.ExecContext(ctx, query, actorid, movieid)
	if err != nil {
		log.Println("Error deleting actor from movie in the table", err)
		return err
	}
	return nil
}
