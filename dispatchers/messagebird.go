package dispatchers

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alfaifiisa/MessagePigeon/models"
	messagebird "github.com/messagebird/go-rest-api"
)

type ThrottledClient struct {
	messageBirdClient *messagebird.Client
	sync.Mutex
	timer    *time.Timer
	duration time.Duration
}

// self destroyed session :)
func (t *ThrottledClient) throttle(duration time.Duration) {
	t.duration = duration
	t.timer = time.NewTimer(t.duration)
	t.timer.Stop()
	go func() {
		for {
			<-t.timer.C
			t.Unlock()
		}
	}()
}

func (t *ThrottledClient) NewMessage(originator string, recipients []string, body string, msgParams *messagebird.MessageParams) (*messagebird.Message, error) {
	t.Lock()
	t.timer.Reset(t.duration)
	//log.Println("sending request to api") //this is used to measure the time between each request
	return t.messageBirdClient.NewMessage(originator, recipients, body, msgParams)
}

var throttledCLient *ThrottledClient

// SendSMSMessage dispach an SMS message to the service.
func SendSMSMessage(smsMessage *models.SMSMessage) (string, error) {
	if smsMessage.GetMessageType() == "single" {
		message, err := throttledCLient.NewMessage(smsMessage.GetOriginator(), smsMessage.GetRecipients(),
			smsMessage.GetSMSMessagePayload()[0].Content, nil)

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
	} else if smsMessage.GetMessageType() == "multipart" {
		payloads := smsMessage.GetSMSMessagePayload()
		log.Println("number of partitions:", len(payloads))
		for index, selectedPart := range payloads {
			params := &messagebird.MessageParams{}
			params.Type = "binary"
			params.TypeDetails = make(map[string]interface{})
			params.TypeDetails["udh"] = selectedPart.UDH
			message, err := throttledCLient.NewMessage(smsMessage.GetOriginator(), smsMessage.GetRecipients(),
				selectedPart.Content, params)
			if err != nil {
				// messagebird.ErrResponse means custom JSON errors.
				if err == messagebird.ErrResponse {
					for _, mbError := range message.Errors {
						fmt.Printf("Error: %#v\n", mbError)
					}
				}
				continue
			}
			smsMessage.UpdatePayloadStatus("sent", index)
			log.Println(message.Body)
		}

		return "all portion has been sent", nil
	}
	return "", errors.New("unsupported payload type: only support single and multipart payloads")
}
