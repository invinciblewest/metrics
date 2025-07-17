package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

var (
	ErrPublicKeyDecoding  = errors.New("не удалось декодировать публичный ключ")
	ErrPrivateKeyDecoding = errors.New("не удалось декодировать приватный ключ")
	ErrNoPublicKey        = errors.New("публичный ключ не задан")
	ErrNoPrivateKey       = errors.New("приватный ключ не задан")
)

type Cryptor struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewCryptor(publicKeyPath, privateKeyPath string) (*Cryptor, error) {
	var publicKey *rsa.PublicKey
	var privateKey *rsa.PrivateKey

	if publicKeyPath != "" {
		publicBytes, err := os.ReadFile(publicKeyPath)
		if err != nil {
			return nil, err
		}

		publicKeyBlock, _ := pem.Decode(publicBytes)
		var pk interface{}
		pk, err = x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
		if err != nil {
			return nil, ErrPublicKeyDecoding
		}

		var ok bool
		publicKey, ok = pk.(*rsa.PublicKey)
		if !ok {
			return nil, ErrPublicKeyDecoding
		}
	}

	if privateKeyPath != "" {
		privateBytes, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, err
		}

		privateKeyBlock, _ := pem.Decode(privateBytes)
		privateKey, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
		if err != nil {
			return nil, ErrPrivateKeyDecoding
		}
	}

	return &Cryptor{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (c *Cryptor) Encrypt(plainText []byte) ([]byte, error) {
	if c.publicKey == nil {
		return nil, ErrNoPublicKey
	}
	return rsa.EncryptPKCS1v15(rand.Reader, c.publicKey, plainText)
}

func (c *Cryptor) Decrypt(cipherText []byte) ([]byte, error) {
	if c.privateKey == nil {
		return nil, ErrNoPrivateKey
	}
	return rsa.DecryptPKCS1v15(nil, c.privateKey, cipherText)
}
