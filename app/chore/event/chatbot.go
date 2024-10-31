package event

import (
	"context"
	"fmt"
	"mywaclient/app/config"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func InitializeChatbot() {
	client := config.GetClient()
	client.AddEventHandler(Chatbot)
}

func Chatbot(evts interface{}) {
	client := config.GetClient()
	switch evt := evts.(type) {
	case *events.Message:
		msg := evt.Message.GetConversation()

		fmt.Println(msg)

		// Send a reply to the message
		client.SendMessage(
			context.Background(),
			evt.Info.Chat,
			&waE2E.Message{
				Conversation: proto.String(handleTextReply(msg)),
			})

	}
}

// You can add more text replies here
func handleTextReply(msg string) string {
	switch strings.ToLower(msg) {
	case "hi", "hello", "hey":
		return "Hello!"
	case "how are you":
		return "I'm fine, thank you"
	case "bye":
		return "Goodbye!"
	}

	return "I'm sorry, I don't understand"
}
