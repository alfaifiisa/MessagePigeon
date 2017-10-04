package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/alfaifiisa/MessagePigeon/dispatchers"
	"github.com/alfaifiisa/MessagePigeon/models"
)

// PostMessageHandler handler the request from a client to store an sms
// messages to be belivered.
func PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var requestTemplate struct {
		Recipient  int    `json:"recipient"`
		Originator string `json:"originator"`
		Message    string `json:"message"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestTemplate)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	oneRecipient := []string{strconv.Itoa(requestTemplate.Recipient)}
	smsMessage, err := models.NewSMSMessage(requestTemplate.Originator, oneRecipient, requestTemplate.Message)
	//log.Println(smsMessage.GetSMSMessagePayload()[0])

	if err != nil {
		log.Println("error", err.Error())
		return
	}
	//	log.Println(smsMessage)
	result, err := dispatchers.SendSMSMessage(smsMessage)
	if err != nil {
		log.Println("error", err.Error())
		return
	}
	log.Println(result)

}
