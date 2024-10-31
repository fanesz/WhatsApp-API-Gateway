package utils

import (
	"io"
	"os"
)

func LoadImage(imageName string) (*[]byte, error) {
	file, err := os.Open(imageName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &byteData, nil
}
