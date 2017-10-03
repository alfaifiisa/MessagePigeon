package dispatchers

import (
	"fmt"

	messagebird "github.com/messagebird/go-rest-api"
)

// TODO: read from a configuration file.
var (
	messageBirdAPIKey = "4XSQaL57j0qn5X37g43oxCoZG"
)

// InitilizeDispachers initilize all available dispatchers and make sure they work.
func InitilizeDispachers() error {
	messageBirdClient = messagebird.New(messageBirdAPIKey)
	balance, err := messageBirdClient.Balance()
	if err != nil {
		// messagebird.ErrResponse means custom JSON errors.
		if err == messagebird.ErrResponse {
			for _, mbError := range balance.Errors {
				fmt.Printf("Error: %#v\n", mbError)
			}
		}
		return err
	}

	// NOTE: here another dispatchers can be initilized
	return nil
}
