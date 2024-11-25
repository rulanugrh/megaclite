package middleware

import (
	"bytes"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
)

type PGPInterface interface {
	Encryption(req domain.Register) ([]byte, error)
	Decryption(req domain.User, armored []byte) (bool, error)
}

type pgp struct {
	utils *crypto.PGPHandle
	conf  config.App
}

func NewPGPMiddleware(utils *crypto.PGPHandle, conf config.App) PGPInterface {
	return &pgp{
		utils: utils,
		conf:  conf,
	}
}
func (p *pgp) Encryption(req domain.Register) ([]byte, error) {
	password := fmt.Sprintf("%s-%s", req.Username, req.Email)
	encryption, err := p.utils.Encryption().Password([]byte(password)).New()
	if err != nil {
		return nil, web.InternalServerError("Error while encrypt new PGP")
	}

	message, err := encryption.Encrypt([]byte(p.conf.Server.Secret))
	if err != nil {
		return nil, web.InternalServerError("Cannot parsing secret server")
	}

	armored, err := message.ArmorBytes()
	if err != nil {
		return nil, web.InternalServerError("Cannot get armor bytes msg")
	}

	return armored, nil
}

func (p *pgp) Decryption(req domain.User, armored []byte) (bool, error) {
	password := fmt.Sprintf("%s-%s", req.Username, req.Email)
	decryption, err := p.utils.Decryption().Password([]byte(password)).New()
	if err != nil {
		return false, web.InternalServerError("Cannot decryption this password")
	}

	decrypted, err := decryption.Decrypt(armored, crypto.Armor)
	if err != nil {
		return false, web.InternalServerError("Cannot get real secret")
	}

	message := decrypted.Bytes()
	if !bytes.Equal(message, []byte(p.conf.Server.Secret)) {
		return false, web.InternalServerError("Sorry secret not matched")
	}

	return true, nil
}
