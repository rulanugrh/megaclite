package middleware

import (
	"log"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
)

type PGPInterface interface {
	GeneratePGPKey(req domain.User) (*web.PGPResponse, error)
}

type pgp struct {
	utils *crypto.PGPHandle
}

func NewPGPMiddleware(utils *crypto.PGPHandle) PGPInterface {
	return &pgp{
		utils: utils,
	}
}
func (p *pgp) GeneratePGPKey(req domain.User) (*web.PGPResponse, error) {
	session, err := p.utils.GenerateSessionKey()
	if err != nil {
		return nil, web.InternalServerError("cannot generate session key: " + err.Error())
	}

	key, err := p.utils.KeyGeneration().
		AddUserId(req.Username, req.Email).
		New().
		GenerateKey()

	if err != nil {
		log.Printf("Cannot generated pgp key:" + err.Error())
	}

	pubkey, err := key.ToPublic()
	if err != nil {
		log.Printf("Cannot parsing into public key:" + err.Error())
	}

	keyRingPrivate, err := crypto.NewKeyRing(key)
	if err != nil {
		log.Printf("Cannot generate key ring private:" + err.Error())
	}

	keyRingPublic, err := crypto.NewKeyRing(pubkey)
	if err != nil {
		log.Printf("Cannot generate ring public:" + err.Error())
	}

	return &web.PGPResponse{
		PrivateKeyRing: keyRingPrivate,
		PrivateKey:     key,
		PublicKeyRing:  keyRingPublic,
		PublicKey:      pubkey,
		SessionKey:     session,
	}, nil

}

func (p *pgp) DecodePGPKey(private *crypto.KeyRing, public *crypto.KeyRing) (crypto.PGPDecryption, error) {
	decryption, err := p.utils.Decryption().
		DecryptionKeys(private).
		VerificationKeys(public).
		New()

	if err != nil {
		return nil, web.InternalServerError("cannot decode this pgp key")
	}

	return decryption, nil
}
