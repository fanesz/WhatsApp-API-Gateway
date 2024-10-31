package utils

import (
	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(data string) (*[]byte, error) {
	var png []byte
	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	return &png, nil
}
