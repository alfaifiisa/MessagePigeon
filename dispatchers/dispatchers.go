package dispatchers

import (
	"fmt"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
)

// TODO: read from a configuration file.
var (
	messageBirdAPIKey = "4XSQaL57j0qn5X37g43oxCoZG"
)

// InitilizeDispachers initilize all available dispatchers and make sure they work.
func InitilizeDispachers() error {

	throttledCLient = &ThrottledClient{}
	throttledCLient.throttle(time.Second * 1)
	throttledCLient.messageBirdClient = messagebird.New(messageBirdAPIKey)
	balance, err := throttledCLient.messageBirdClient.Balance()
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
