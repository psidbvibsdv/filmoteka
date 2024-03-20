package moviestorage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"filmoteka/internal/domain/models"
	"log"
	"time"
)

var DbTimeout = 5 * time.Second

type MovieStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *MovieStorage {
	return &MovieStorage{
		db: db,
	}

}

func (s *MovieStorage) GetAllMovies(sortParam string) ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
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
	rows, err := s.db.QueryContext(ctx, query)
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
		if err != nil {
			log.Println("Error scanning movie rows", err)
			return nil, err
		}

		movies = append(movies, &movie)
	}
	if len(movies) < 1 {
		log.Println("No movies found in the table")
		return nil, models.ErrNoRecord
	}
	return movies, nil
}

func (s *MovieStorage) GetMovieByID(id int) (*models.Movie, error) {
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

func (s *MovieStorage) CreateMovie(m *models.Movie) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `INSERT INTO movies (title, description, rating, releasedate)
	values ($1, $2, $3, $4)`
	_, err := s.db.ExecContext(ctx, query, m.Title, m.Description, m.Rating, m.ReleaseDate)
	if err != nil {
		log.Println("Error inserting movie into a table", err)
		encoder, _ := json.Marshal(m)
		log.Println("Movie: ", string(encoder))
		return nil, err
	}
	return m, nil
}

func (s *MovieStorage) UpdateMovie(m *models.Movie) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()
	formattedRD, _ := json.Marshal(m.ReleaseDate.Time)
	log.Println("Formatted DOB: ", string(formattedRD))
	query := `UPDATE movies SET title = CASE WHEN $1 = '' THEN title ELSE COALESCE($1) END, description = CASE WHEN $2 = '' THEN description ELSE COALESCE($2) END, rating = CASE WHEN $3 = 0 THEN rating ELSE COALESCE($3) END, releasedate = CASE WHEN $4 < '1000-1-1' THEN releasedate ELSE CAST(COALESCE($4) AS DATE) END WHERE movieid = ($5)`
	_, err := s.db.ExecContext(ctx, query, m.Title, m.Description, m.Rating, formattedRD, m.MovieID)
	if err != nil {
		log.Println("Error updating movie in the table", err)
		encoder, _ := json.Marshal(m)
		log.Println("Movie: ", string(encoder))
		return nil, err
	}
	res, err := s.GetMovieByID(m.MovieID)
	log.Println("Updated movie: ", res)
	if err != nil {
		log.Println("Error returning updated movie from the table ", err)
		return nil, err
	}
	return res, nil
}

func (s *MovieStorage) DeleteMovie(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()
	query := `DELETE FROM movies WHERE movieid = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Println("Error deleting movie from the table", err)
		return err
	}
	return nil
}

func (s *MovieStorage) GetMovieByMovieName(moviename string) ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()
	query := `SELECT movieid, title, description, rating, releasedate FROM movies WHERE title ILIKE $1 OR description ILIKE $2;`

	rows, err := s.db.QueryContext(ctx, query, "%"+moviename+"%", "%"+moviename+"%")
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
		if err != nil {
			log.Println("Error scanning actor rows", err)
			return nil, err
		}

		movies = append(movies, &movie)
	}

	if len(movies) < 1 {
		log.Println("No movies found in the table")
		return nil, models.ErrNoRecord
	}
	return movies, nil
}
