package utils

import (
	"log"
	"github.com/nats-io/nats.go"
)

var NatsConn *nats.Conn
var KvStore nats.KeyValue

func InitNATSConnection(natsURL string) int {
	log.Printf("Connecting to NATS at %s", natsURL)
	var err error
	NatsConn, err = nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
		return -1
	}
	log.Println("Connected to NATS")
	return 0
}

// InitKVStore initializes the NATS Key-Value store
func InitKVStore(bucketName string) int {
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
		return -1
	}

	log.Println("NATS KV Store initialized")
	return 0
}
