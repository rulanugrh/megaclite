package service

import (
	// "io/fs"

	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserInterface interface {
	Register(req domain.Register) (*web.PGPResponse, error)
	Login(req domain.Login, private string) (*web.GetUser, error)
	GetEmail(email string) (*web.GetUser, error)
	GetByID(id uint) (*web.GetUser, error)
}

type user struct {
	repository repository.UserInterface
	middleware middleware.PGPInterface
	validation middleware.IValidation
}

func NewUserService(repository repository.UserInterface, middlewares middleware.PGPInterface) UserInterface {
	return &user{
		repository: repository,
		middleware: middlewares,
		validation: middleware.NewValidation(),
	}
}

func (u *user) Register(req domain.Register) (*web.PGPResponse, error) {
	err := u.validation.Validate(req)
	if err != nil {
		return nil, u.validation.ValidationMessage(err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, web.BadRequest("Error while hashed password")
	}

	request := domain.Register{
		Username: req.Username,
		Password: string(hashed),
		Email:    req.Email,
	}

	data, err := u.repository.Register(request)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	key, err := u.middleware.GenerateKeygen(req)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.PGPResponse{
		Private:  key.Private,
		Username: data.Username,
	}

	return &response, nil
}

func (u *user) Login(req domain.Login, private string) (*web.GetUser, error) {
	err := u.validation.Validate(req)
	if err != nil {
		return nil, u.validation.ValidationMessage(err)
	}

	data, err := u.repository.Login(req)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(req.Password))
	if err != nil {
		return nil, web.BadRequest("Sorry password not matched")
	}

	check, err := u.middleware.VerificationKey(private)
	if !check && err != nil {
		return nil, web.BadRequest("Sorry secret message and keygen not matched")
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
