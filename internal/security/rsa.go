package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

// ReadRSAPublicKeyFromFile читает открытый ключ RSA из файла
func ReadRSAPublicKeyFromFile(filePath string) (*rsa.PublicKey, error) {
	// Читаем PEM-блок из файла
	block, err := readPemBlockFromFile(filePath)
	if err != nil {
		return nil, err
	}

	//
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

// ReadRSAPrivateKeyFromFile читает закрытый ключ RSA из файла
func ReadRSAPrivateKeyFromFile(filePath string) (*rsa.PrivateKey, error) {
	// Читаем PEM-блок из файла
	block, err := readPemBlockFromFile(filePath)
	if err != nil {
		return nil, err
	}

	// Парсим закрытый ключ из PEM-блока
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key.(*rsa.PrivateKey), nil
}

// EncryptRSA шифрует данные с использованием публичного ключа
func EncryptRSA(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
}

// DecryptRSA расшифровывает данные с использованием приватного ключа
func DecryptRSA(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
}
