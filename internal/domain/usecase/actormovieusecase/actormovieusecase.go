package actormovieusecase

import (
	"filmoteka/internal/domain/models"
)

type ActorMovieUseCase struct {
	storage ActorMovieStorage
}

type ActorMovieStorage interface {
	GetActorsForMovie(id int) ([]*models.Actor, *models.Movie, error)
	GetMoviesForActor(actorid int) ([]*models.Movie, *models.Actor, error)
	GetActorsAndMoviesForMovie(id int) ([]*models.ActorMovies, *models.Movie, error)
	AddActorToMovie(actorid int, movieid int) (*models.Actor, *models.Movie, error)
	DeleteActorFromMovie(actorid int, movieid int) (*models.Actor, *models.Movie, error)
	GetMovieByActorName(name string, surname string) ([]*models.MovieWithActor, error)
}

func New(storage ActorMovieStorage) *ActorMovieUseCase {
	return &ActorMovieUseCase{
		storage: storage,
	}
}

func (uc *ActorMovieUseCase) GetActorsForMovie(id int) ([]*models.Actor, *models.Movie, error) {
	return uc.storage.GetActorsForMovie(id)
}

func (uc *ActorMovieUseCase) GetMoviesForActor(actorid int) ([]*models.Movie, *models.Actor, error) {
	return uc.storage.GetMoviesForActor(actorid)
}

func (uc *ActorMovieUseCase) GetActorsAndMoviesForMovie(id int) ([]*models.ActorMovies, *models.Movie, error) {
	return uc.storage.GetActorsAndMoviesForMovie(id)
}

func (uc *ActorMovieUseCase) AddActorToMovie(actorid int, movieid int) (*models.Actor, *models.Movie, error) {
	return uc.storage.AddActorToMovie(actorid, movieid)
}

func (uc *ActorMovieUseCase) DeleteActorFromMovie(actorid int, movieid int) (*models.Actor, *models.Movie, error) {
	return uc.storage.DeleteActorFromMovie(actorid, movieid)
}

func (uc *ActorMovieUseCase) GetMovieByActorName(name string, surname string) ([]*models.MovieWithActor, error) {
	return uc.storage.GetMovieByActorName(name, surname)
}
