package repository

import (
	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
)

type UserInterface interface {
	Register(req domain.Register) (*domain.User, error)
	Login(req domain.Login) (*domain.User, error)
	Get(id uint) (*domain.User, error)
	GetMail(email string) (*domain.User, error)
}

type user struct {
	connection config.Database
}

func NewUserRepository(config config.Database) UserInterface {
	return &user{
		connection: config,
	}
}
func (u *user) Register(req domain.Register) (*domain.User, error) {
	var response domain.User
	find := u.connection.DB.Exec("SELECT * FROM users WHERE email = ?", req.Email)
	if find.RowsAffected > 0 {
		return nil, web.InternalServerError("Sorry email has been taken")
	}

	err := u.connection.DB.Exec("INSERT INTO users(username, email, password) VALUES (?,?,?)",
		req.Username,
		req.Email,
		req.Password,
	).Find(&response).Error

	if err != nil {
		return nil, web.InternalServerError("Cannot save data user to Database")
	}

	return &response, nil
}

func (u *user) Login(req domain.Login) (*domain.User, error) {
	var response domain.User

	err := u.connection.DB.Raw("SELECT * FROM users WHERE email = ?", req.Email).Scan(&response).Error
	if err != nil {
		return nil, web.InternalServerError("Cant find user with this email")
	}

	return &response, nil
}

func (u *user) Get(id uint) (*domain.User, error) {
	var response domain.User

	err := u.connection.DB.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&response).Error
	if err != nil {
		return nil, web.InternalServerError("Cant find user with this id")
	}

	return &response, nil
}

func (u *user) GetMail(email string) (*domain.User, error) {
	var response domain.User

	err := u.connection.DB.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&response).Error
	if err != nil {
		return nil, web.InternalServerError("Cant find user with this email")
	}

	return &response, nil
}
