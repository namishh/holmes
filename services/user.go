package services

import (
	"log"

	"github.com/namishh/holmes/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	Level     int    `json:"level"`
	CreatedAt string `json:"created_at"`
}

type UserService struct {
	User      User
	UserStore database.DatabaseStore
}

func NewUserService(user User, userStore database.DatabaseStore) *UserService {
	return &UserService{
		User:      user,
		UserStore: userStore,
	}
}

func (us *UserService) CreateUser(u User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// create user himself
	stmt := `INSERT INTO teams (email, password, name) VALUES ($1, $2, $3)`

	s, err := us.UserStore.DB.Exec(stmt, u.Email, string(hashedPassword), u.Username)

	log.Print(s)

	return err
}

func (us *UserService) CheckUsername(usr string) (User, error) {
	query := `SELECT id, email, password, username FROM teams
		WHERE name = ?`

	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	us.User.Username = usr
	err = stmt.QueryRow(
		us.User.Username,
	).Scan(
		&us.User.ID,
		&us.User.Email,
		&us.User.Password,
		&us.User.Username,
	)
	if err != nil {
		return User{}, err
	}

	return us.User, nil
}

func (us *UserService) CheckEmail(email string) (User, error) {
	query := `SELECT id, email, password, name FROM teams
		WHERE email = ?`

	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	us.User.Email = email
	err = stmt.QueryRow(
		us.User.Email,
	).Scan(
		&us.User.ID,
		&us.User.Email,
		&us.User.Password,
		&us.User.Username,
	)
	if err != nil {
		return User{}, err
	}

	return us.User, nil
}
