package config

import (
	"fmt"
	"mywaclient/app/database"

	"go.mau.fi/whatsmeow"
)

var client *whatsmeow.Client

func Initialize() {
	fmt.Println("===== Initialize Client =====")

	dbInstance := database.GetDBInstance()
	device, err := dbInstance.GetFirstDevice()
	if err != nil {
		fmt.Println("failed to get first device from database")
		return
	}

	client = whatsmeow.NewClient(device, nil)

	if client.Store.ID != nil {
		if err := client.Connect(); err != nil {
			fmt.Println("failed to connect to whatsapp")
			return
		}
	}

	fmt.Println("✓ Client initialized")
}

func GetClient() *whatsmeow.Client {
	if client == nil {
		Initialize()
	}

	return client
}

func ResyncClient(callerClient **whatsmeow.Client) error {
	if client.IsConnected() {
		client.Disconnect()
	}

	if client != nil {
		client = nil
	}

	Initialize()
	if client == nil {
		return fmt.Errorf("failed to reinitialize client")
	}

	*callerClient = client
	return nil
}
