package service

import (
	"io/fs"
	"os"

	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/repository"
)

type UserInterface interface {
	Register(req domain.Register) (*web.GetUser, error)
	Login(req domain.Login) (*web.GetUser, error)
	GetEmail(email string) (*web.GetUser, error)
	GetByID(id uint) (*web.GetUser, error)
}

type user struct {
	repository repository.UserInterface
	middleware middleware.PGPInterface
}

func NewUserService(repository repository.UserInterface, middleware middleware.PGPInterface) UserInterface {
	return &user{
		repository: repository,
		middleware: middleware,
	}
}

func (u *user) Register(req domain.Register) (*web.GetUser, error) {
	data, err := u.repository.Register(req)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.GetUser{
		Username: data.Username,
		Email:    data.Email,
		Address:  data.Address,
	}

	err = os.Mkdir("./tmp/"+req.Username, fs.FileMode(os.O_RDWR))
	if err != nil {
		return nil, web.InternalServerError("Cannot create new folder")
	}

	pgp, err := u.middleware.GeneratePGPKey(req)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	println(pgp)

	return &response, nil
}

func (u *user) Login(req domain.Login) (*web.GetUser, error) {
	data, err := u.repository.Login(req)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.GetUser{
		Username: data.Username,
		Email:    data.Email,
		Address:  data.Address,
	}

	return &response, nil
}

func (u *user) GetEmail(email string) (*web.GetUser, error) {
	data, err := u.repository.GetMail(email)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.GetUser{
		Username: data.Username,
		Email:    data.Email,
		Address:  data.Address,
	}

	return &response, nil
}

func (u *user) GetByID(id uint) (*web.GetUser, error) {
	data, err := u.repository.Get(id)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.GetUser{
		Username: data.Username,
		Email:    data.Email,
		Address:  data.Address,
	}

	return &response, nil
}
