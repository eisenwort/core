package ewc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

func encrypt(text string) ([]byte, error) {
	data := []byte(text)
	key := []byte("passphrasewhichneedstobe32bytes!")
	cip, err := aes.NewCipher(key)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(cip)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decrypt(cipherText []byte) (string, error) {
	key := []byte("passphrasewhichneedstobe32bytes!")
	cip, err := aes.NewCipher(key)

	if err != nil {
		log.Println(err)
		return "", err
	}

	gcm, err := cipher.NewGCM(cip)

	if err != nil {
		log.Println(err)
		return "", err
	}

	nonceSize := gcm.NonceSize()

	if len(cipherText) < nonceSize {
		log.Println(err)
		return "", err
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(plaintext), nil
}

func decryptDb() {
	cipherText, err := ioutil.ReadFile(connectionString)

	if err != nil {
		log.Println(err)
	}

	data, err := decrypt(cipherText)

	if err != nil {
		return
	}
	if err := ioutil.WriteFile(connectionString, []byte(data), 0777); err != nil {
		log.Println("decrypt db write file error:", err)
	}
}

func encryptDb() {
	data, err := ioutil.ReadFile(connectionString)

	if err != nil {
		log.Println(err)
	}

	cip, err := encrypt(string(data))

	if err != nil {
		return
	}
	if err := ioutil.WriteFile(connectionString, cip, 0777); err != nil {
		log.Println("encrypt db write file error:", err)
	}
}
