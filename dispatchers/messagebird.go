package dispatchers

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/alfaifiisa/MessagePigeon/models"
	messagebird "github.com/messagebird/go-rest-api"
)

var messageBirdClient *messagebird.Client
var m = &sync.Mutex{}

// TODO: change the recipient to recipients

// SendSMSMessage dispach an SMS message to the service.
func SendSMSMessage(smsMessage *models.SMSMessage, udh string) (string, error) {
	m.Lock()
	if smsMessage.GetMessageType() == "single" {
		message, err := messageBirdClient.NewMessage(smsMessage.GetOriginator(), smsMessage.GetRecipients(),
			smsMessage.GetSMSMessagePayload()[0].Content, nil)
		if err != nil {
			return "", err
		}
		return message.Body, nil
	} else if smsMessage.GetMessageType() == "multipart" {
		payloads := smsMessage.GetSMSMessagePayload()
		selectedPart := payloads[0]
		log.Println("payload contains", payloads)
		params := &messagebird.MessageParams{}
		params.Type = "binary"
		params.TypeDetails = make(map[string]interface{})
		params.TypeDetails["udh"] = selectedPart.UDH
		message, err := messageBirdClient.NewMessage(smsMessage.GetOriginator(), smsMessage.GetRecipients(),
			selectedPart.Content, params)
		if err != nil {
			// messagebird.ErrResponse means custom JSON errors.
			if err == messagebird.ErrResponse {
				for _, mbError := range message.Errors {
					fmt.Printf("Error: %#v\n", mbError)
				}
			}
		}
		if err != nil {
			return "", err
		}
		return message.Body, nil
	}
	m.Unlock()
	return "", errors.New("unsupported payload type: only support single and multipart payloads")
}
