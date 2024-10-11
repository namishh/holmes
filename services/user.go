package services

import (
	"github.com/minio/minio-go/v7"
	"github.com/namishh/holmes/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	Points    int    `json:"points"`
	CreatedAt string `json:"created_at"`
}

type UserService struct {
	User        User
	UserStore   database.DatabaseStore
	MinioClient *minio.Client
}

func NewUserService(user User, userStore database.DatabaseStore, mini *minio.Client) *UserService {
	return &UserService{
		User:        user,
		UserStore:   userStore,
		MinioClient: mini,
	}
}

func (us *UserService) CreateUser(u User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO teams (email, password, name, points) VALUES ($1, $2, $3, 1000)`

	_, err = us.UserStore.DB.Exec(stmt, u.Email, string(hashedPassword), u.Username)
	return err
}

func (us *UserService) CheckUsername(usr string) (User, error) {
	query := `SELECT id, email, password, name, points FROM teams
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
		&us.User.Points,
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

func (us *UserService) GetAllUsers() ([]User, error) {
	query := `SELECT id, email, name, points FROM teams`
	users := make([]User, 0)
	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return users, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Email, &u.Username, &u.Points)
		if err != nil {
			return users, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (us *UserService) DeleteTeam(id int) error {
	query := `DELETE FROM teams WHERE id = ?`
	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	stmt.Exec(id)

	return nil
}
