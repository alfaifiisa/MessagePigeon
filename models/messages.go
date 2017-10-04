package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
)

// SMSMessage represents an sms message to be send, and has two types single or multipart
type SMSMessage struct {
	messageType string // single or mutipart
	originator  string
	recipients  []string
	payload     []SMSMessagePayload
}

// GetMessageType is an accessor for messageType
func (s *SMSMessage) GetMessageType() string {
	return s.messageType
}

// GetOriginator is an accessor for originator
func (s *SMSMessage) GetOriginator() string {
	return s.originator
}

// GetRecipients is an accessor for recipients
func (s *SMSMessage) GetRecipients() []string {
	return s.recipients
}

// GetSMSMessagePayload is an accessor for recipient
func (s *SMSMessage) GetSMSMessagePayload() []SMSMessagePayload {

	return s.payload
}

func (s *SMSMessage) UpdatePayloadStatus(status string, index int) error {
	if index < len(s.payload) && (status == "new" || status == "sent") {

		s.payload[index].Status = status
		return nil
	}

	// TODO: this should return an error for each case, NOT a general error
	return errors.New("invalid put values")

}

// SMSMessagePayload serves the purpose of sending multipart SMS along with simple SMS message
type SMSMessagePayload struct {
	UDH     string
	Content string
	Status  string // send or new
}

// NewSMSMessage is to create a new message object (constructor)
func NewSMSMessage(originator string, recipients []string, messageBody string) (*SMSMessage, error) {
	smsMessage := &SMSMessage{originator: originator, recipients: recipients}
	if len(messageBody) <= 160 {
		smsMessage.messageType = "single"
		smsMessage.payload = append(smsMessage.payload, SMSMessagePayload{"", messageBody, "new"})
	} else {
		/* sample udh 05 00 03 34 02 01
				   1. generate a starting random byte for message reference number byte 3
		           2. calculate byte 4 by deviding message length by 153. +1 if there is a reminder
		           3. keep incrementing byte 5 every time
		*/
		/* partitioning algorithm
		   1. start with the full sms message
		   2. slice by 153 and store as a partition until length is less than 153
		   3. save the remanining of the message as a partition
		*/
		referneceNumber := make([]byte, 1)
		if _, err := rand.Read(referneceNumber); err != nil {
			return nil, errors.New("cannot generate a random reference number for messages")
		}
		// TODO: watch for maximum number for the messages partitions it will overflow the sequence number when its in bytes
		numberOfParts := len(messageBody) / 153
		if len(messageBody)%10 > 0 {
			numberOfParts++
		}
		udh := []byte{0x05, 0x00, 0x03, referneceNumber[0], byte(numberOfParts), 0x00}

		smsMessage.messageType = "multipart"

		var part string
		for len(messageBody) > 0 {
			udh[5]++
			if len(messageBody) >= 153 {
				part = messageBody[:153]
				messageBody = messageBody[153:]
			} else {
				part = messageBody[:]
				messageBody = messageBody[len(messageBody):]
			}
			smsMessagePayload := SMSMessagePayload{hex.EncodeToString(udh), hex.EncodeToString([]byte(part)), "new"}
			smsMessage.payload = append(smsMessage.payload, smsMessagePayload)
		}
		log.Println("partitioning finished")
	}
	return smsMessage, nil
}
