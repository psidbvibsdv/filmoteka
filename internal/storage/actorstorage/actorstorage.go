package actorstorage

import (
	"context"
	"database/sql"
	"encoding/json"
	"filmoteka/internal/domain/models"
	"log"
	"time"
)

var DbTimeout = 5 * time.Second

type ActorStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *ActorStorage {
	return &ActorStorage{
		db: db,
	}
}

func (s *ActorStorage) GetAllActors() ([]*models.Actor, error) {
	var actors []*models.Actor
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `SELECT actorid, name, gender, dateofbirth
	FROM actors ORDER BY name`

	rows, err := s.db.QueryContext(ctx, query)
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

	for rows.Next() {
		var actor models.Actor
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
		return nil, models.ErrNoRecord
	}

	return actors, nil

}

func (s *ActorStorage) CreateActor(a *models.Actor) (*models.Actor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `INSERT INTO actors (name, gender, dateofbirth)
	values ($1, $2, $3)`
	_, err := s.db.ExecContext(ctx, query, a.Name, a.Gender, a.DateOfBirth)
	if err != nil {
		log.Println("Error inserting actor into a table", err)
		encoder, _ := json.Marshal(a)
		log.Println("Actor: ", string(encoder))
		return nil, err
	}
	return a, nil
}

func (s *ActorStorage) GetActorByID(id int) (*models.Actor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()
	query := `SELECT actorid, name, gender, dateofbirth FROM actors WHERE actorid = $1`
	rows, err := s.db.QueryContext(ctx, query, id)
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

	actor := &models.Actor{}
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
		return nil, models.ErrNoRecord
	}
	log.Println("Actor: ", actor.ActorID, actor.DateOfBirth.Time, actor.DateOfBirth.Valid)
	return actor, nil
}

func (s *ActorStorage) UpdateActor(a *models.Actor) (*models.Actor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()
	log.Println(a.DateOfBirth)
	formattedDOB, _ := json.Marshal(a.DateOfBirth.Time)
	log.Println("Formatted DOB: ", string(formattedDOB))
	query := `UPDATE actors SET name = CASE WHEN $1 = '' THEN name ELSE COALESCE($1) END, gender = CASE WHEN $2 = '' THEN gender ELSE COALESCE($2) END, dateofbirth = CASE WHEN $3 < '1000-1-1' THEN dateofbirth ELSE CAST(COALESCE($3) AS DATE) END WHERE actorid = ($4)`
	_, err := s.db.ExecContext(ctx, query, a.Name, a.Gender, formattedDOB, a.ActorID)
	if err != nil {
		log.Println("Error updating actor in the table", err)
		encoder, _ := json.Marshal(a)
		log.Println("Actor: ", string(encoder))
		return nil, err
	}
	res, err := s.GetActorByID(a.ActorID)
	log.Println("Updated actor: ", res)
	if err != nil {
		log.Println("Error error returning updated actor from the table", err)
		return nil, err
	}
	return res, nil
}

func (s *ActorStorage) DeleteActor(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()
	query := `DELETE FROM actors WHERE actorid = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Println("Error deleting actor from the table", err)
		return err
	}
	return nil
}
