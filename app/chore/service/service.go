package service

import (
	"context"
	"errors"
	"fmt"
	"mywaclient/app/chore/entity"
	"mywaclient/app/chore/interfaces"
	"mywaclient/app/config"
	"mywaclient/app/utils"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

var _ interfaces.WhatsappService = &whatsappService{}

type whatsappService struct {
	client *whatsmeow.Client
}

func NewWhatsappService(client *whatsmeow.Client) *whatsappService {
	return &whatsappService{
		client: client,
	}
}

func (s *whatsappService) CheckDevice() (bool, error) {
	if s.client == nil {
		return false, fmt.Errorf("client is not initialized")
	}

	isIDStored := s.client.Store.ID != nil  // Check if the ID is stored
	isAutenticated := s.client.IsLoggedIn() // Check if the client is authenticated

	return isIDStored && isAutenticated, nil
}

func (s *whatsappService) GetLoginQR() (<-chan *[]byte, error) {
	isLogin, err := s.CheckDevice()
	if err != nil {
		return nil, err
	}

	// If the device is already logged in, return nil
	if isLogin {
		return nil, nil
	}

	// If the client is nil, resync the client
	// this is to make sure that the client is not nil
	if s.client != nil {
		config.ResyncClient(&s.client)
	}

	// Get the QR channel
	qrChan, err := s.client.GetQRChannel(context.Background())
	if err != nil {
		if errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
			_ = s.client.Connect()
			if s.client.IsLoggedIn() {
				return nil, nil
			}
			return nil, err
		} else {
			return nil, err
		}
	}

	// Start a loop to listen to the QR channel
	// and keep generating the QR code if expired
	out := make(chan *[]byte)
	go func() {
		defer close(out)
		for evt := range qrChan {
			switch evt.Event {
			case "success":
				out <- nil
				return
			case "timeout":
				out <- nil
				return
			case "code":
				// You can display the QR code to the terminal if you want
				png, err := utils.GenerateQRCode(evt.Code)
				if err != nil {
					continue
				}
				out <- png
			}
		}
	}()

	// Connect the client after the QR code is scanned
	err = s.client.Connect()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *whatsappService) SendMessage(req *entity.MessageSend) error {
	logged, err := s.CheckDevice()
	if err != nil {
		return err
	}
	if !logged {
		return fmt.Errorf("device is not logged in")
	}

	// Costumize the utils.ParsePhoneNumber function
	// with the correct phone number format of your country
	targetJID := types.NewJID(utils.ParsePhoneNumber(req.To), types.DefaultUserServer)

	var message waE2E.Message
	if req.ImageID == "" { // If the image ID is empty, send a text message
		message.Conversation = proto.String(req.Message)
	} else { // If the image ID is not empty, send an image message
		byteData, err := utils.LoadImage(req.ImageID)
		if err != nil {
			return err
		}

		// Upload the image to the server
		uploadRes, _ := s.client.Upload(context.Background(), *byteData, whatsmeow.MediaImage)
		message.ImageMessage = &waE2E.ImageMessage{
			Caption:       proto.String("Hello, world!"),
			Mimetype:      proto.String("image/png"),
			URL:           &uploadRes.URL,
			DirectPath:    &uploadRes.DirectPath,
			MediaKey:      uploadRes.MediaKey,
			FileEncSHA256: uploadRes.FileEncSHA256,
			FileSHA256:    uploadRes.FileSHA256,
			FileLength:    &uploadRes.FileLength,
		}
	}

	// Send the message
	_, err = s.client.SendMessage(context.Background(), targetJID, &message)
	if err != nil {
		return err
	}

	return nil
}

func (s *whatsappService) ResetLoggedDevice() error {
	if s.client == nil {
		return fmt.Errorf("client is not initialized")
	}

	s.client.Store.Container.DeleteDevice(s.client.Store)
	s.client.Store.Delete()
	s.client.Disconnect()

	return nil
}
