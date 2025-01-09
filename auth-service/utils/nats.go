package utils

import (
	"log"
	"github.com/nats-io/nats.go"
)

var NatsConn *nats.Conn
var KvStore nats.KeyValue

func InitNATSConnection(natsURL string) {
	log.Printf("Connecting to NATS at %s", natsURL)
	var err error
	NatsConn, err = nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	log.Println("Connected to NATS")
}

// InitKVStore initializes the NATS Key-Value store
func InitKVStore(bucketName string) {
	var err error
	var js nats.JetStreamContext

	js, err = NatsConn.JetStream()
	if err != nil {
		log.Fatalf("Error getting JetStream context: %v", err)
	}

	KvStore, err = js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: bucketName,
	})
	if err != nil {
		log.Fatalf("Error creating KV store: %v", err)
	}

	log.Println("NATS KV Store initialized")
}
