package repository

import "github.com/alfaifiisa/MessagePigeon/dispatchers"

// StartWorker starts the messaging dispaching routine
func StartWorker() {
	go func() {
		for {
			//log.Println("start processing a message from the queue")
			dispatchers.SendSMSMessage(<-queue)
			//log.Println("finish processing a message from the queue")
		}
	}()
}
