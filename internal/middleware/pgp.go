package middleware

import (
	"fmt"
	"log"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
)

type PGPInterface interface {
	Encryption(req domain.MailRegister) ([]byte, error)
	Decryption(req domain.Mail) ([]byte, error)
	GenerateKeygen(req domain.Register) (*web.PGPResponse, error)
	VerificationKey(private string) (string, bool, error)
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
func (p *pgp) Encryption(req domain.MailRegister) ([]byte, error) {
	password := fmt.Sprintf("%s-%s", req.From, req.To)
	encryption, err := p.utils.Encryption().Password([]byte(password)).New()
	if err != nil {
		return nil, web.InternalServerError("Error while encrypt new PGP")
	}

	message, err := encryption.Encrypt([]byte(req.Message))
	if err != nil {
		return nil, web.InternalServerError("Cannot parsing secret server")
	}

	armored, err := message.ArmorBytes()
	if err != nil {
		return nil, web.InternalServerError("Cannot get armor bytes msg")
	}

	return armored, nil
}

func (p *pgp) Decryption(req domain.Mail) ([]byte, error) {
	password := fmt.Sprintf("%s-%s", req.From, req.To)
	decryption, err := p.utils.Decryption().Password([]byte(password)).New()
	if err != nil {
		return nil, web.InternalServerError("Cannot decryption this password")
	}

	decrypted, err := decryption.Decrypt([]byte(req.Message), crypto.Armor)
	if err != nil {
		return nil, web.InternalServerError("Cannot get real secret")
	}

	message := decrypted.Bytes()
	return message, nil
}

func (p *pgp) GenerateKeygen(req domain.Register) (*web.PGPResponse, error) {
	keygen, err := p.utils.KeyGeneration().AddUserId(req.Username, req.Email).New().GenerateKey()
	if err != nil {
		return nil, web.InternalServerError("Cannot generate key")
	}

	armored, err := keygen.Armor()
	if err != nil {
		return nil, web.InternalServerError("Cannot get armord public key")
	}

	private, err := crypto.NewPrivateKeyFromArmored(armored, []byte(req.Username))
	if err != nil {
		return nil, web.InternalServerError("Cannot get private key")
	}

	armoredPrivate, err := private.Armor()
	if err != nil {
		log.Fatal("Parsing armor from private key error: " + err.Error())
	}

	return &web.PGPResponse{
		Private:  armoredPrivate,
		HexKeyID: private.GetEntity().PrivateKey.KeyIdString(),
	}, nil
}

func (p *pgp) VerificationKey(private string) (string, bool, error) {
	privateKey := strings.NewReader(private)

	keyReader, err := openpgp.ReadArmoredKeyRing(privateKey)
	if err != nil {
		return "", false, web.BadRequest("Cannot Read Private Key")
	}

	if len(keyReader) == 0 {
		return "", false, web.BadRequest("Sorry Private Key Not Found")
	}

	id := keyReader[0].PrivateKey.KeyIdString()
	log.Println(id)
	return id, true, nil

}
