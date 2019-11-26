package ewc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"log"
)

type CryptService struct {
	EncryptChan chan string
	DecryptChan chan string
	ResultChan  chan string
}

func NewCryptService() *CryptService {
	srv := new(CryptService)
	srv.EncryptChan = make(chan string, chanSize)
	srv.DecryptChan = make(chan string, chanSize)
	srv.ResultChan = make(chan string, chanSize)

	return srv
}

func (srv *CryptService) Encrypt(encryptData EncryptData) string {
	cipherText, err := rsa.EncryptOAEP(
		encryptData.Hash,
		rand.Reader,
		encryptData.PublicKey,
		encryptData.Message,
		encryptData.Label,
	)
	if err != nil {
		log.Println("encrypt error:", err)
	}
	return string(cipherText)
}

func (srv *CryptService) Decrypt(cipherText string, priv *rsa.PrivateKey) string {
	hash := sha512.New()
	msgBytes, err := rsa.DecryptOAEP(hash, rand.Reader, priv, []byte(cipherText), nil)

	if err != nil {
		log.Println("decrypt error:", err)
	}

	return string(msgBytes)
}
