package service

import (
	"log"

	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/middleware"
	"github.com/rulanugrh/megaclite/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserInterface interface {
	Register(req domain.Register) (*web.PGPResponse, error)
	Login(req domain.Login, private string) (*web.ResponseLogin, error)
	GetEmail(email string) (*web.GetUser, error)
	GetByID(id uint) (*web.GetUser, error)
	UpdatePassword(email string, password string) error
	UpdateProfile(email string, req domain.User) error
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

	key, err := u.middleware.GenerateKeygen(req)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	request := domain.Register{
		Username: req.Username,
		Password: string(hashed),
		Email:    req.Email,
		KeygenID: key.HexKeyID,
	}

	data, err := u.repository.Register(request)
	if err != nil {
		return nil, web.InternalServerError(err.Error())
	}

	response := web.PGPResponse{
		Private:  key.Private,
		Username: data.Username,
	}

	return &response, nil
}

func (u *user) Login(req domain.Login, private string) (*web.ResponseLogin, error) {
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

	id, check, err := u.middleware.VerificationKey(private)
	if !check && err != nil {
		return nil, web.BadRequest(err.Error())
	}

	if id != data.KeygenID {
		return nil, web.BadRequest("Sorry Your ID not matched")
	}

	response := web.ResponseLogin{
		KeyID:  data.KeygenID,
		UserID: data.ID,
		Email:  data.Email,
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
		Avatar:   data.Avatar,
		KeyID:    data.KeygenID,
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
		Avatar:   data.Avatar,
		KeyID:    data.KeygenID,
	}

	return &response, nil
}

func (u *user) UpdatePassword(email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return web.InternalServerError("Sorry cannot generate hash password")
	}

	if err = u.repository.UpdatePassword(email, string(hashedPassword)); err != nil {
		return web.InternalServerError(err.Error())
	}

	return nil
}

func (u *user) UpdateProfile(email string, req domain.User) error {
	if err := u.repository.UpdateProfile(email, req); err != nil {
		log.Println(err.Error())
		return web.InternalServerError(err.Error())
	}

	return nil
}
