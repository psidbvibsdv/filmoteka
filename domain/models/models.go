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
	Actor      Actor
	Movie      Movie
	ActorMovie ActorMovie
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Actor struct {
	ActorID     int    `json:"actorid,omitempty"`
	Name        string `json:"name,omitempty"`
	Gender      string `json:"gender,omitempty"`
	DateOfBirth Date   `json:"dateofbirth,omitempty"`
}

type ActorMovie struct {
	ActorID int `json:"actor_id,omitempty"`
	MovieID int `json:"movie_id,omitempty"`
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
		Actor:      Actor{},
		Movie:      Movie{},
		ActorMovie: ActorMovie{},
	}
}

// Methods for User
func (u *User) auth() {
	bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(u.Password))
}

func (u *User) isAuthenticated() {
}

func (u *User) isAdmin() {
}

func (u *User) GetByEmail() (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT password, role FROM users WHERE email = $1`
	_, err := db.QueryContext(ctx, query, u.Email)
	if err != nil {
		log.Println("Error getting user by email from the table", err)
		return "", "", err
	}
	return "", "", err
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
		log.Println("Error getting actor by id from the table", err)
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
