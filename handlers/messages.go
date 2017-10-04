package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alfaifiisa/MessagePigeon/models"
	"github.com/alfaifiisa/MessagePigeon/repository"
)

// PostMessageHandler handler the request from a client to store an sms
// messages to be belivered.
func PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		ServeResponse(w, false, nil, apiError{1, "wrong Content-Type, it has to be application/json"})
		return
	}

	// only taking one receiping as the requirments stated
	var requestTemplate struct {
		Recipient  int    `json:"recipient"`
		Originator string `json:"originator"`
		Message    string `json:"message"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestTemplate)
	if err != nil {
		ServeResponse(w, false, nil, apiError{1, err.Error()})
		return
	}

	if !validateRecipientMSISDN(strconv.Itoa(requestTemplate.Recipient)) {
		ServeResponse(w, false, nil, apiError{2, "recipient is not valid"})
		return
	}

	if len(requestTemplate.Originator) == 0 {
		ServeResponse(w, false, nil, apiError{3, "originator not provided"})
		return
	}
	if len(requestTemplate.Message) == 0 {
		ServeResponse(w, false, nil, apiError{4, "message not provided"})
		return
	} else if len(requestTemplate.Message) > 1377 {
		ServeResponse(w, false, nil, apiError{5, "message exceeds the 1377 character long"})
		return
	}

	// only taking one receiping as the requirments stated
	oneRecipient := []string{strconv.Itoa(requestTemplate.Recipient)}
	smsMessage, err := models.NewSMSMessage(requestTemplate.Originator, oneRecipient, requestTemplate.Message)

	if err != nil {
		ServeResponse(w, false, nil, apiError{6, err.Error()})
		return
	}

	// to add the constructed message to the worker queue
	err = repository.SendSMSMessage(smsMessage)
	if err != nil {
		ServeResponse(w, false, nil, apiError{7, err.Error()})
		return
	}

	ServeResponse(w, true, "message is queued to be sent")

}
