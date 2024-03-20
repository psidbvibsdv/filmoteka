package movieusecase

import (
	"filmoteka/internal/domain/models"
)

type MovieUseCase struct {
	storage movieStorage
}

func New(movieStorage movieStorage) *MovieUseCase {
	return &MovieUseCase{
		storage: movieStorage,
	}
}

type movieStorage interface {
	GetAllMovies(param string) ([]*models.Movie, error)
	CreateMovie(m *models.Movie) (*models.Movie, error)
	GetMovieByID(id int) (*models.Movie, error)
	UpdateMovie(m *models.Movie) (*models.Movie, error)
	DeleteMovie(id int) error

	GetMovieByMovieName(moviename string) ([]*models.Movie, error)
}

func (uc *MovieUseCase) GetAllMovies(param string) ([]*models.Movie, error) {
	return uc.storage.GetAllMovies(param)
}

func (uc *MovieUseCase) CreateMovie(m *models.Movie) (*models.Movie, error) {
	return uc.storage.CreateMovie(m)
}

func (uc *MovieUseCase) GetMovieByID(id int) (*models.Movie, error) {
	return uc.storage.GetMovieByID(id)
}

func (uc *MovieUseCase) UpdateMovie(m *models.Movie) (*models.Movie, error) {
	return uc.storage.UpdateMovie(m)
}

func (uc *MovieUseCase) DeleteMovie(id int) error {
	return uc.storage.DeleteMovie(id)
}

func (uc *MovieUseCase) GetMovieByMovieName(moviename string) ([]*models.Movie, error) {
	return uc.storage.GetMovieByMovieName(moviename)
}
