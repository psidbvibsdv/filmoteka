package actormoviestorage

import (
	"context"
	"database/sql"
	"filmoteka/internal/domain/models"
	"github.com/pkg/errors"
	"log"
	"time"
)

var DbTimeout = 5 * time.Second

type ActorMovieStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *ActorMovieStorage {
	return &ActorMovieStorage{
		db: db,
	}
}

func (s *ActorMovieStorage) GetMovieByID(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()
	query := `SELECT movieid, title, description, rating, releasedate FROM movies WHERE movieid = $1`
	rows, err := s.db.QueryContext(ctx, query, id)
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

	movie := &models.Movie{}
	for rows.Next() {
		err = rows.Scan(
			&movie.MovieID,
			&movie.Title,
			&movie.Description,
			&movie.Rating,
			&movie.ReleaseDate,
		)
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No results", err)
			return nil, models.ErrNoRecord
		} else if err != nil {
			log.Println("Error getting rows", err)
			return nil, err
		}
	}

	if !movie.ReleaseDate.Valid {
		return nil, models.ErrNoRecord
	}
	log.Println("Movie: ", movie)
	return movie, nil
}

func (s *ActorMovieStorage) GetActorByID(id int) (*models.Actor, error) {
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
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No results", err)
			return nil, models.ErrNoRecord
		} else if err != nil {
			log.Println("Error scanning rows", err)
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

func (s *ActorMovieStorage) AddActorToMovie(actorid int, movieid int) (*models.Actor, *models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()
	a := &models.Actor{ActorID: actorid}
	m := &models.Movie{MovieID: movieid}
	a, err := s.GetActorByID(a.ActorID)
	if errors.Is(err, models.ErrNoRecord) {
		log.Println("No results", err)
		return nil, nil, models.ErrNoRecord
	} else if err != nil {
		log.Println("Error getting actor from the table", err)
		return nil, nil, err
	}
	m, err = s.GetMovieByID(m.MovieID)
	if errors.Is(err, models.ErrNoRecord) {
		log.Println("No results", err)
		return nil, nil, err
	} else if err != nil {
		log.Println("Error getting movies for actor", err)
		return nil, nil, err
	}
	query := `INSERT INTO actormovie (actorid, movieid) VALUES ($1, $2)`
	_, err = s.db.ExecContext(ctx, query, actorid, movieid)
	if err != nil {
		log.Println("Error adding actor to movie in the table", err)
		return nil, nil, err
	}
	return a, m, nil
}

func (s *ActorMovieStorage) DeleteActorFromMovie(actorid int, movieid int) (*models.Actor, *models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()
	a := &models.Actor{ActorID: actorid}
	m := &models.Movie{MovieID: movieid}
	a, err := s.GetActorByID(a.ActorID)
	if errors.Is(err, models.ErrNoRecord) {
		log.Println("No results", err)
		return nil, nil, err
	} else if err != nil {
		log.Println("Error getting actor by id from the table", err)
		return nil, nil, err
	}
	m, err = s.GetMovieByID(m.MovieID)
	if errors.Is(err, models.ErrNoRecord) {
		log.Println("No results", err)
		return nil, nil, err
	} else if err != nil {
		log.Println("Error getting movies by id from the table", err)
		return nil, nil, err
	}
	query := `DELETE FROM actormovie WHERE actorid = $1 AND movieid = $2`
	_, err = s.db.ExecContext(ctx, query, actorid, movieid)
	if err != nil {
		log.Println("Error deleting actor from movie in the table", err)
		return nil, nil, err
	}
	return a, m, nil
}

func (s *ActorMovieStorage) GetActorsForMovie(id int) ([]*models.Actor, *models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `SELECT a.ActorID, a.Name, a.Gender, a.DateOfBirth FROM Actors a JOIN actormovie am ON a.actorid = am.actorid WHERE am.movieid = $1`

	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		log.Println("Error getting actors for movie from the table", err)
		return nil, nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows", err)
		}
	}(rows)

	var actors []*models.Actor
	for rows.Next() {
		var actor models.Actor
		err = rows.Scan(
			&actor.ActorID,
			&actor.Name,
			&actor.Gender,
			&actor.DateOfBirth,
		)
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No results", err)
			return nil, nil, models.ErrNoRecord
		} else if err != nil {
			log.Println("Error scanning rows", err)
			return nil, nil, err
		}

		actors = append(actors, &actor)

	}
	movie, err := s.GetMovieByID(id)
	if errors.Is(err, models.ErrNoRecord) {
		log.Println("No results", err)
		return nil, nil, err
	} else if err != nil {
		log.Println("Error scanning rows", err)
		return nil, nil, err
	}
	return actors, movie, nil
}

func (s *ActorMovieStorage) GetMoviesForActor(actorid int) ([]*models.Movie, *models.Actor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `SELECT m.movieid, m.title, m.description, m.rating, m.releasedate FROM Movies m JOIN actormovie am ON m.movieid = am.movieid WHERE am.actorid = $1`

	rows, err := s.db.QueryContext(ctx, query, actorid)
	if err != nil {
		log.Println("Error getting movies for actor from the table", err)
		return nil, nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows", err)
		}
	}(rows)

	var movies []*models.Movie
	for rows.Next() {
		var movie models.Movie
		err = rows.Scan(
			&movie.MovieID,
			&movie.Title,
			&movie.Description,
			&movie.Rating,
			&movie.ReleaseDate,
		)
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No results", err)
			return nil, nil, models.ErrNoRecord
		} else if err != nil {
			log.Println("Error scanning rows", err)
			return nil, nil, err
		}

		movies = append(movies, &movie)

	}
	actor, err := s.GetActorByID(actorid)
	if errors.Is(err, models.ErrNoRecord) {
		log.Println("No results", err)
		return nil, nil, err
	} else if err != nil {
		log.Println("Error getting actor by id from the table", err)
		return nil, nil, err

	}

	return movies, actor, nil

}

func (s *ActorMovieStorage) GetMovieByActorName(name string, surname string) ([]*models.MovieWithActor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `SELECT m.*, a.name AS actor_name
FROM Movies m
         JOIN actormovie am ON m.movieid = am.movieid
         JOIN Actors a ON am.actorid = a.actorid 
WHERE a.name ILIKE  $1 AND a.name ILIKE $2;`

	rows, err := s.db.QueryContext(ctx, query, "%"+name+"%", "%"+surname+"%")
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

	var movies []*models.MovieWithActor

	for rows.Next() {
		var movie models.MovieWithActor
		err = rows.Scan(
			&movie.MovieID,
			&movie.Title,
			&movie.Description,
			&movie.Rating,
			&movie.ReleaseDate,
			&movie.ActorName,
		)
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No results", err)
			return nil, models.ErrNoRecord
		} else if err != nil {
			log.Println("Error scanning rows", err)
			return nil, err
		}

		movies = append(movies, &movie)
	}

	return movies, nil

}

func (s *ActorMovieStorage) GetActorsAndMoviesForMovie(id int) ([]*models.ActorMovies, *models.Movie, error) {
	actors, _, err := s.GetActorsForMovie(id)
	movie := &models.Movie{}
	if errors.Is(err, models.ErrNoRecord) {
		log.Println("No results", err)
		return nil, nil, err
	} else if err != nil {
		log.Println("Error getting actors for movie", err)
		return nil, nil, err
	}
	var result []*models.ActorMovies
	for _, actor := range actors {
		movies, _, err := s.GetMoviesForActor(actor.ActorID)
		if errors.Is(err, models.ErrNoRecord) {
			log.Println("No results", err)
			return nil, nil, err
		} else if err != nil {
			log.Println("Error getting movies for actor", err)
			return nil, nil, err
		}
		result = append(result, &models.ActorMovies{
			ActorId: actor.ActorID,
			Name:    actor.Name,
			Movies:  movies,
		})
	}
	movie, err = s.GetMovieByID(id)
	if errors.Is(err, models.ErrNoRecord) {
		log.Println("No results", err)
		return nil, nil, err
	} else if err != nil {
		log.Println("Error scanning rows", err)
		return nil, nil, err
	}
	return result, movie, nil
}
