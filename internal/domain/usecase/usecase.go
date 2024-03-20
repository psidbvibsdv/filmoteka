package usecase

import (
	"filmoteka/internal/domain/usecase/actormovieusecase"
	"filmoteka/internal/domain/usecase/actorusecase"
	"filmoteka/internal/domain/usecase/movieusecase"
	"filmoteka/internal/domain/usecase/userusecase"
)

type UseCase struct {
	UserUseCase       *userusecase.UserUseCase
	MovieUseCase      *movieusecase.MovieUseCase
	ActorMovieUseCase *actormovieusecase.ActorMovieUseCase
	ActorUseCase      *actorusecase.ActorUseCase
}

//func New(storage storage) *UseCase {
//	return &UseCase{
//		storage: storage,
//	}
//}
