package repository

import (
	"errors"

	"github.com/alfaifiisa/MessagePigeon/models"
)

var (
	bufferSize = 10000
	queue      = make(chan *models.SMSMessage, bufferSize)
)

// SendSMSMessage using buffered channel to store message in the memory to be processed.
func SendSMSMessage(smsMessage *models.SMSMessage) error {
	// TODO: follwing check might fail when there a multiple go routines accesing this method
	// which might cause this to pass the error and try adding to a full channel which will cause blocking of this method
	if len(queue) == bufferSize {
		return errors.New("server is busy please send your request again")
	}

	queue <- smsMessage
	return nil
}
