package userstorage

import (
	"context"
	"database/sql"
	"log"
	"time"
)

var DbTimeout = 5 * time.Second

type UserStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}

}

//func (s *UserStorage) auth(u *models.User) {
//	bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(u.Password))
//}

func (s *UserStorage) GetUserByEmail(email string) (password string, role string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	query := `SELECT password, role FROM users WHERE email = $1`
	_, err = s.db.QueryContext(ctx, query, email)
	if err != nil {
		log.Println("Error getting user by email from the table", err)
		return "", "", err
	}

	err = s.db.QueryRowContext(ctx, query, email).Scan(&password, &role)
	if err != nil {
		log.Println("Error getting user by email from the table", err)
		return "", "", err
	}

	return password, role, err
}
