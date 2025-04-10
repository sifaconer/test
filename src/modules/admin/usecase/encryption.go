package usecase

import (
	"api-test/src/common"
	"api-test/src/config"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

type encryption struct {
	log    common.Logger
	config *config.Config
}

func (e *encryption) GenerateRandomKey() ([]byte, error) {
	key := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

func (e *encryption) GenerateRandomPassword() (string, error) {
	key, err := e.GenerateRandomKey()
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(key), nil
}

func (e *encryption) decodeMasterKey() ([]byte, error) {
	masterKey, err := base64.StdEncoding.DecodeString(e.config.MasterEncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error decodificando clave maestra: %w", err)
	}

	// Validar longitud de la clave (AES-256 requiere 32 bytes)
	if len(masterKey) != 32 {
		return nil, errors.New("la clave maestra debe ser de 32 bytes (256 bits)")
	}
	return masterKey, nil
}

func (e *encryption) Encrypt(plaintext string) (ciphertext, iv []byte, err error) {
	masterKey, err := e.decodeMasterKey()
	if err != nil {
		return nil, nil, err
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, nil, err
	}

	iv = make([]byte, 12) // GCM 12 bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	ciphertext = aesGCM.Seal(nil, iv, []byte(plaintext), nil)

	return ciphertext, iv, nil
}

func (e *encryption) Decrypt(ciphertext, iv []byte) (string, error) {
	masterKey, err := e.decodeMasterKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesGCM.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func NewEncryption(log common.Logger, config *config.Config) *encryption {
	return &encryption{
		log:    log,
		config: config,
	}
}
