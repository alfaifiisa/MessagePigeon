package repository

import (
	"log"

	"github.com/alfaifiisa/MessagePigeon/dispatchers"
)

func StartWorker() {
	go func() {
		for {
			log.Println("start processing a message from the queue")
			dispatchers.SendSMSMessage(<-queue)
			log.Println("finish processing a message from the queue")
		}
	}()
}
