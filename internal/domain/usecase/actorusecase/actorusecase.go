package actorusecase

import (
	"filmoteka/internal/domain/models"
)

type ActorUseCase struct {
	storage ActorStorage
}

func New(storage ActorStorage) *ActorUseCase {
	return &ActorUseCase{
		storage: storage,
	}
}

type ActorStorage interface {
	GetAllActors() ([]*models.Actor, error)
	CreateActor(a *models.Actor) (*models.Actor, error)
	GetActorByID(id int) (*models.Actor, error)
	UpdateActor(a *models.Actor) (*models.Actor, error)
	DeleteActor(id int) error
}

func (uc *ActorUseCase) GetAllActors() ([]*models.Actor, error) {
	return uc.storage.GetAllActors()
}

func (uc *ActorUseCase) CreateActor(a *models.Actor) (*models.Actor, error) {
	return uc.storage.CreateActor(a)
}

func (uc *ActorUseCase) GetActorByID(id int) (*models.Actor, error) {
	return uc.storage.GetActorByID(id)
}

func (uc *ActorUseCase) UpdateActor(a *models.Actor) (*models.Actor, error) {
	return uc.storage.UpdateActor(a)
}

func (uc *ActorUseCase) DeleteActor(id int) error {
	return uc.storage.DeleteActor(id)
}
