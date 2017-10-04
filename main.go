package main

import (
	"log"
	"net/http"

	"github.com/alfaifiisa/MessagePigeon/dispatchers"
	"github.com/alfaifiisa/MessagePigeon/handlers"
)

var (
	port = "7070"
)

// TODO: add api keys and authentication when providing the api

func main() {
	//message, _ := models.NewMessage("3982435792", "lksdjg", test)
	//log.Println(message)
	log.Println("!Welcome to MessagePigeon!")
	log.Println("starting the server on port", port)
	err := dispatchers.InitilizeDispachers()

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("message dispatchers are ready :)")
	http.HandleFunc("/messages", handlers.PostMessageHandler)

	// TODO: add the support to tls, serve both on 8080 and 8443
	http.ListenAndServe(":"+port, nil)
}
