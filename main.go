package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/alfaifiisa/MessagePigeon/dispatchers"
	"github.com/alfaifiisa/MessagePigeon/handlers"
	"github.com/alfaifiisa/MessagePigeon/repository"
)

// TODO: add api keys and authentication when providing the api

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "port to use for the server")
	flag.Parse()
	log.Println("!Welcome to MessagePigeon!")
	log.Println("starting the server on port", port)
	// load external Messaging APIs, initilize keys and check validity and conectivity.
	err := dispatchers.InitilizeDispachers()
	if err != nil {
		log.Fatal(err.Error())
	}
	// start messaging worker to receive messages from clients and regulate dispatching
	repository.StartWorker()

	log.Println("message dispatchers are ready :)")
	http.HandleFunc("/messages", handlers.PostMessageHandler)

	// TODO: add the support to tls, serve both on 8080 and 8443
	http.ListenAndServe(":"+port, nil)
}
