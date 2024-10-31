package interfaces

import "mywaclient/app/chore/entity"

type WhatsappService interface {
	CheckDevice() (bool, error)
	GetLoginQR() (<-chan *[]byte, error)
	SendMessage(req *entity.MessageSend) error
	ResetLoggedDevice() error
}
