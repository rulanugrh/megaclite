package repository

import (
	"fmt"
	"log"

	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
)

type UserInterface interface {
	Register(req domain.Register) (*domain.User, error)
	Login(req domain.Login) (*domain.User, error)
	Get(id uint) (*domain.User, error)
	GetMail(email string) (*domain.User, error)
	UpdatePassword(email string, password string) error
	UpdateProfile(email string, req domain.User) error
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
	find := u.connection.DB.Raw("SELECT * FROM users WHERE email = ?", req.Email).Scan(&response)
	if find.RowsAffected != 0 {
		return nil, web.InternalServerError("Sorry email has been taken")
	}

	request := domain.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		KeygenID: req.KeygenID,
	}

	err := u.connection.DB.Create(&request).Error
	if err != nil {
		return nil, web.InternalServerError("Cannot save data user to Database")
	}

	return &request, nil
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
		log.Println(err.Error())
		return nil, web.InternalServerError("Cant find user with this email")
	}

	return &response, nil
}

func (u *user) UpdatePassword(email string, password string) error {
	err := u.connection.DB.Where("email = ?", email).Update("password", &password).Error
	if err != nil {
		return web.InternalServerError("Sorry cannot update password")
	}

	return nil
}

func (u *user) UpdateProfile(email string, req domain.User) error {
	if err := u.connection.DB.Where("email = ?", email).Update("avatar", &req.Avatar).Update("address", &req.Address).Update("username", &req.Username).Error; err != nil {
		return web.InternalServerError(fmt.Sprintf("Sorry cannot update profile somthing error: %s", err.Error()))
	}

	return nil
}
