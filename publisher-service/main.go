package main

import (
	"log"
	"net/http"
	"publisher-service/handlers"

	"github.com/nats-io/nats.go"

	"publisher-service/config"
)

func main() {
	log.Printf("Connection to NATS")
	cfg := config.LoadConfig()

	nc, err := nats.Connect(cfg.NATS_URL)
	log.Printf("Connectin to %s ...", cfg.NATS_URL)

	if err != nil {
		log.Printf("Error connecting %v", err)
	}

	log.Printf("Connection successful")
	/*
		err = nc.Publish("topic", []byte("Ciao"))
		if err != nil {
			log.Printf("error %v\n", err)
		}*/

	publisherHandler := handlers.NewPublisherHandler(nc, "functions.execution", 30)

	http.HandleFunc("/publisher/publish", publisherHandler.PublishHandlerMethod)
	log.Printf("server listening on port 8083")
	if err := http.ListenAndServe(cfg.SERVER_ADDRESS, nil); err != nil { //prima inizializzazione e poi condizione vera
		log.Printf("error running server: %v\n", err)
	}

}
