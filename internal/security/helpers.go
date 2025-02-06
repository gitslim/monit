package security

import (
	"encoding/pem"
	"errors"
	"io"
	"os"
)

// readPemBlockFromFile читает блок данных из файла с ключом в формате PEM.
func readPemBlockFromFile(keyFilePath string) (*pem.Block, error) {
	// Открываем файл с ключом
	file, err := os.Open(keyFilePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	// Читаем содержимое файла
	keyBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	//
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	return block, err
}
