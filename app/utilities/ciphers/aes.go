package ciphers

import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/aes"
	"os"
	"io"
)

func EncryptWithAES(message []byte) ([]byte, error) {
	key := getEncryptionKey()
	aesc, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesc)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	cipherbytes := gcm.Seal(nonce, nonce, []byte(message), nil)

	return cipherbytes, nil
}

func DecryptFromAES(cipherbytes []byte) ([]byte, error) {
	key := getEncryptionKey()

	aesc, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesc)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphermessage := cipherbytes[:nonceSize], cipherbytes[nonceSize:]
	message, err := gcm.Open(nil, nonce, ciphermessage, nil)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func getEncryptionKey() []byte {
	return []byte(os.Getenv("CHI_YT_DB_ENC_KEY"))
}
