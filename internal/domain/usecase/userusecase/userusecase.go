package userusecase

type UserUseCase struct {
	userStorage userStorage
}

func New(userStorage userStorage) *UserUseCase {
	return &UserUseCase{
		userStorage: userStorage,
	}
}

type userStorage interface {
	//auth(u *models.User)
	GetUserByEmail(email string) (string, string, error)
}

//func (uc *UserUseCase) auth(u *models.User) {
//	uc.userStorage.auth(u)
//}

func (uc *UserUseCase) GetUserByEmail(email string) (string, string, error) {
	return uc.userStorage.GetUserByEmail(email)
}
